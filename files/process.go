package files

import (
	"fmt"
	"log"
	"path"

	"github.com/RedMapleTech/filehandling-go/details"
)

const (
	authResults = "Authentication-Results"
	subject     = "Subject"

	emlFrom        = "From"
	emlReturnPath  = "Return-Path"
	emlTo          = "To"
	emlDate        = "Date"
	emlMessageID   = "Message-ID"
	emlContentType = "Content-Type"

	msgSender        = "Sender name"
	msgDisplayName   = "Sender Simple Display Name"
	msgSenderSMTP    = "Sender SMTP Address"
	msgSenderEmail   = "Sender Email"
	msgSenderEmail2  = "Sender Email 2"
	msg7bitEmail     = "Seven Bit Email"
	msgReceivedName  = "Received by name"
	msgReceivedSMTP  = "Received By SMTP Address"
	msgReceivedEmail = "Received by email"
	messageTopic     = "Topic"
	msgMessageID     = "MessageID"
	msgMSIPLabel     = "Microsoft Information Protection (MSIP) Label"

	emlMimeType = "text/plain; charset=utf-8"
	msgMimeType = "application/vnd.ms-outlook"
	pdfMimeType = "application/pdf"
)

type FileProperties struct {
	FileName string
	Hash     string
	FileType string
	Size     string
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

	matches, explanation := checkMime(path.Ext(filePath), details.Mimetype)

	if !matches {
		return nil, fmt.Errorf("mismatched extension and MIME type. %s", explanation)
	}

	return &props, nil
}

func checkMime(extension, mime string) (bool, string) {
	switch extension {
	case ".msg":
		if mime != msgMimeType {
			return false, fmt.Sprintf("We expect %s for files with %s extensions, but found %s", msgMimeType, extension, mime)
		}
	case ".eml":
		if mime != emlMimeType {
			return false, fmt.Sprintf("We expect %s for files with %s extensions, but found %s", emlMimeType, extension, mime)
		}
	case ".pdf":
		if mime != pdfMimeType {
			return false, fmt.Sprintf("We expect %s for files with %s extensions, but found %s", emlMimeType, extension, mime)
		}
	}

	return true, ""
}

type ProcessResult struct {
	FilePath  string
	Error     error
	Parsed    bool
	Completed bool
	Dangerous bool
	Analysis  string
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
	// case ".docx":
	// 	log.Println("Parsing document file")
	default:
		res.Completed = false
		res.Error = fmt.Errorf("unknown file extension %q", fileExt)
		return &res
	}

	return &res
}
