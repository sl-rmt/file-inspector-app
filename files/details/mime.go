package details

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"

	"file-inspector/files/hashing"

	"github.com/gabriel-vasile/mimetype"
)

// FindFoldersWithFileTypes walks the provided rootpath and returns
// folders containing files matching the filetype
func FindFoldersWithFileTypes(rootPath string, mimetypes []string) ([]*FileDetails, map[string]int) {

	folders := make(map[string]int)
	var fileCount int

	// setup ID workers
	// at some point disk access will slow us, 8 routines probably enough.
	// Too many and it will get behind the file listing
	const numWorkers = 8

	// buffered channels
	// problem if more than 1000000 files
	results := make(chan *FileDetails, 1000000)

	// plenty to keep all the workers busy
	jobs := make(chan *FileDetails, 1000)

	var wg sync.WaitGroup

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go getFileTypeWorker(mimetypes, jobs, results, &wg)
	}

	// define inline function to walk filesystem
	err := filepath.Walk(rootPath, func(thisFilePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// ignore dirs
		if !info.IsDir() {
			fileCount++

			var details FileDetails
			details.Path = thisFilePath
			details.Size = info.Size()

			// send it to be identified
			jobs <- &details
		}

		return nil
	})

	// all jobs sent - close the channel
	close(jobs)

	// panic if the walker returned an error
	if err != nil {
		log.Panicf("%s\n", err)
	}

	// wait for them all to finish
	wg.Wait()

	// results received, close the channel
	close(results)

	// process results
	var files []*FileDetails

	// put them all into a single slice
	for item := range results {
		files = append(files, item)
	}

	folders = countFiles(files, folders)

	return files, folders
}

// FindFoldersWithFileType walks the provided rootpath and returns
// folders containing files matching the filetype
func FindFoldersWithFileType(rootPath, mimetype string) []string {

	// TODO replace with map[string]int
	var folders []string

	// define inline function to walk filesystem
	err := filepath.Walk(rootPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// ignore dirs
		if !info.IsDir() {
			// get mime type
			fileType, err := GetFileType(filePath)

			if err != nil {
				log.Printf("Error trying to type file %q: %s", filePath, err.Error())
			} else {
				if fileType == mimetype {
					//log.Printf("Found Binary file: %q (type %s)", filePath, mimetype)
					folderPath, _ := path.Split(filePath)
					folders = append(folders, folderPath)
				}
			}
		}

		return nil
	})

	// panic if the walker returned an error
	if err != nil {
		log.Panicf("%s\n", err)
	}

	return folders
}

// FindFilesInFolderByType walks the provided rootpath and
// returns any file found that match the mimetype
func FindFilesInFolderByType(rootPath, mimetype string) []string {
	var files []string

	// define inline function to walk filesystem
	err := filepath.Walk(rootPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// ignore dirs
		if !info.IsDir() {
			// get mime type
			fileType, err := GetFileType(filePath)

			if err != nil {
				log.Printf("Error trying to type file %q: %s", filePath, err.Error())
			} else {
				if fileType == mimetype {
					//log.Printf("Found Binary file: %q (type %s)", filePath, mimetype)
					files = append(files, filePath)
				}
			}
		}

		return nil
	})

	// panic if the walker returned an error
	if err != nil {
		log.Panicf("%s\n", err)
	}

	return files
}

// GetFileType returns the mime type of the file
func GetFileType(filePath string) (string, error) {

	mime, err := mimetype.DetectFile(filePath)

	if err != nil {
		return "", fmt.Errorf("error getting the type of file %q", filePath)
	}

	return mime.String(), nil
}

func getFileTypeWorker(mimetypes []string, jobs <-chan *FileDetails, results chan<- *FileDetails, wg *sync.WaitGroup) {

	defer wg.Done()

	for details := range jobs {
		// list of types here: https://github.com/gabriel-vasile/mimetype/blob/master/supported_mimes.md
		mime, err := mimetype.DetectFile(details.Path)

		if err == nil {
			for _, mimetype := range mimetypes {
				if mime.String() == mimetype {
					details.Mimetype = mime.String()
					details.SHA256 = hashing.GetFileSHA256HashString(details.Path)
					results <- details
					continue
				}
			}
		}
	}
}

func countFiles(files []*FileDetails, folders map[string]int) map[string]int {

	for _, file := range files {
		dir := filepath.Dir(file.Path)

		if _, ok := folders[dir]; ok {
			folders[dir] = folders[dir] + 1
		} else {
			folders[dir] = 1
		}
	}

	return folders
}
