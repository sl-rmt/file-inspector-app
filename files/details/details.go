package details

import (
	"fmt"
	"os"
	"path"

	"github.com/dustin/go-humanize"
	"file-inspector/files/hashing"
	"file-inspector/files/checks"
)

// FileDetails contains details about a file. Duh.
type FileDetails struct {
	Path       string
	SHA256     string
	Size       int64
	SizeString string
	Mimetype   string
}

// GetFile tries to open the file, returns file pointer and error
func GetFile(rootFolder, filePath string) (*os.File, error) {

	fullPath := path.Join(rootFolder, filePath)

	if !checks.FileExists(fullPath) {
		return nil, fmt.Errorf("file %q Doesn't Exist", fullPath)
	}

	return os.Open(fullPath)
}

// GetFileSize returns the file size as an int64
func GetFileSize(fileString string) (int64, error) {

	info, err := os.Stat(fileString)

	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}

// GetFileSizeString returns the file size as a human readable size string
func GetFileSizeString(fileString string) (string, error) {

	info, err := os.Stat(fileString)

	if err != nil {
		return "", err
	}

	return humanize.Bytes(uint64(info.Size())), nil
}

// GetFileDetails returns the file details
func GetFileDetails(fileString string) (*FileDetails, error) {

	var details FileDetails
	details.Path = fileString

	f, err := os.Open(fileString)

	if err != nil {
		return nil, err
	}
	defer f.Close()

	info, err := os.Stat(fileString)

	if err != nil {
		return nil, err
	}

	details.Size = info.Size()
	details.SizeString = humanize.Bytes(uint64(info.Size()))
	details.SHA256 = hashing.GetFileSHA256HashString(fileString)
	details.Mimetype, err = GetFileType(fileString)

	if err != nil {
		return nil, err
	}

	return &details, nil
}
