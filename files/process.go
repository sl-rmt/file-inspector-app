package files

import (
	"fmt"
	"log"
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
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
)

func SetFileProperties(filePath string, fileName, hash, fileType, size binding.String, window *fyne.Window) error {
	// set file name
	fileName.Set(filePath)

	// get details
	details, err := details.GetFileDetails(filePath)

	// set them if no error
	if err != nil {
		return err
	} else {
		fileType.Set(details.Mimetype)
		hash.Set(details.SHA256)
		size.Set(details.SizeString)
	}

	matches, explanation := checkMime(path.Ext(filePath), details.Mimetype)

	if !matches {
		return fmt.Errorf("mismatched extension and MIME type. %s", explanation)
	}

	return nil
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
	}

	return true, ""
}

func ProcessFile(filePath string, window *fyne.Window, displayText binding.String) (  error) {
	log.Printf("Processing file %q\n", filePath)
	fileExt := path.Ext(filePath)
	var err error

	switch fileExt {
	case ".msg":
		log.Println("Parsing email file")
		err = processMsgFile(filePath, displayText)
	case ".eml":
		log.Println("Parsing email file")
		err = processEmlFile(filePath, displayText)
	// case ".docx":
	// 	log.Println("Parsing document file")
	default:
		return fmt.Errorf("unknown file extension %q", fileExt)
	}

	return err
}
