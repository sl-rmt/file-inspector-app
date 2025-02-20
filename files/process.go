package files

import (
	"fmt"
	"log"
	"path"

	"file-inspector/files/details"
)

const (
	emlMimeType  = "text/plain; charset=utf-8"
	msgMimeType  = "application/vnd.ms-outlook"
	pdfMimeType  = "application/pdf"
	docxMimeType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
)

type FileProperties struct {
	FileName string
	Hash     string
	FileType string
	Size     string
}

type ProcessResult struct {
	FilePath  string
	Error     error
	Parsed    bool
	Completed bool
	Dangerous bool
	Metadata  string
	Analysis  string
}

func GetFileProperties(filePath string) (*FileProperties, error) {
	props := FileProperties{
		FileName: filePath,
	}

	// get details
	details, err := details.GetFileDetails(filePath)

	// set them if no error
	if err != nil {
		return nil, err
	} else {
		props.FileType = details.Mimetype
		props.Hash = details.SHA256
		props.Size = details.SizeString
	}

	return &props, nil
}

func CheckMime(extension, mime string) (bool, string) {
	switch extension {
	case ".msg":
		if mime != msgMimeType {
			return false, fmt.Sprintf("☠️ We expect %q for files with .msg extensions, but found %q.", msgMimeType, mime)
		}
	case ".eml":
		if mime != emlMimeType {
			return false, fmt.Sprintf("☠️ We expect %q for files with .eml extensions, but found %q.", emlMimeType, mime)
		}
	case ".pdf":
		if mime != pdfMimeType {
			return false, fmt.Sprintf("☠️ We expect %q for files with .pdf extensions, but found %q.", pdfMimeType, mime)
		}
	case ".docx":
		if mime != docxMimeType {
			return false, fmt.Sprintf("☠️ We expect %q for files with .docx extensions, but found %q.", docxMimeType, mime)
		}
	}

	return true, ""
}

func ProcessFile(filePath string) *ProcessResult {
	log.Printf("Processing file %q\n", filePath)
	fileExt := path.Ext(filePath)

	res := ProcessResult{
		FilePath: filePath,
		Error:    nil,
	}

	switch fileExt {
	case ".msg":
		log.Println("Parsing email file")
		processMsgFile(&res)
	case ".eml":
		log.Println("Parsing email file")
		processEmlFile(&res)
	case ".pdf":
		log.Println("Parsing PDF file")
		processPDFFile(&res)
	case ".docx":
	 	log.Println("Parsing document file")
		processDocxFile(&res)
	default:
		res.Completed = false
		res.Error = fmt.Errorf("unknown file extension %q", fileExt)
		return &res
	}

	return &res
}
