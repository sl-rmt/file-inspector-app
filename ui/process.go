package main

import (
	"fmt"
	"log"
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/RedMapleTech/email-parse/emlparse"
	"github.com/RedMapleTech/email-parse/msgparse"
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
)

func getProps(filePath string, fileName, hash, fileType, size binding.String, window *fyne.Window) error {
	// set file name
	fileName.Set(filePath)

	// get details
	details, err := details.GetFileDetails(filePath)

	// set them if no error
	if err != nil {
		launchErrorDialog(err, window)
		return err
	} else {
		fileType.Set(details.Mimetype)
		hash.Set(details.SHA256)
		size.Set(details.SizeString)
	}

	matches := checkMime(path.Ext(filePath), details.Mimetype, window)

	if !matches {
		return fmt.Errorf("mismatched extension and MIME type")
	}

	return nil
}

func checkMime(extension, mime string, window *fyne.Window) bool {
	switch extension {
	case ".msg":
		if mime != "application/vnd.ms-outlook" {
			launchInfoDialog("Unexpected File Type", fmt.Sprintf("File MIME type %q is not the expected type for %q files. Parsing aborted.", mime, extension), window)
			return false
		}
	case ".eml":
		if mime != "text/plain; charset=utf-8" {
			launchInfoDialog("Unexpected File Type", fmt.Sprintf("File MIME type %q is not the expected type for %q files. Parsing aborted.", mime, extension), window)
			return false
		}
	}

	return true
}

func processFile(filePath string, window *fyne.Window, displayText binding.String) {
	log.Printf("Processing file %q\n", filePath)
	fileExt := path.Ext(filePath)

	switch fileExt {
	case ".msg":
		processMsgFile(filePath, window, displayText)
	case ".eml":
		log.Println("Parsing email file")
		processEmlFile(filePath, window, displayText)
	// case ".docx":
	// 	log.Println("Parsing document file")
	default:
		launchInfoDialog("Unknown file", fmt.Sprintf("Can only inspect known file types, not %s", fileExt), window)
	}
}

func processMsgFile(filePath string, window *fyne.Window, displayText binding.String) {
	msg, err := msgparse.ReadMsgFile(filePath, false)

	if err != nil {
		launchErrorDialog(err, window)
		return
	}

	log.Println("Email parsing done")

	keyFieldNames := []string{msgSender, msgDisplayName, msgSenderSMTP, msgSenderEmail, msgSenderEmail2, msgReceivedName, msgReceivedSMTP, msg7bitEmail, msgReceivedEmail, subject, messageTopic, msgMessageID}

	var analysis string

	// Print values
	for _, fieldName := range keyFieldNames {
		field := msg.GetPropertyByName(fieldName)
		analysis = fmt.Sprintf("%s\n%s: %q", analysis, fieldName, field)
	}

	displayText.Set(analysis)
}

func processEmlFile(filePath string, window *fyne.Window, displayText binding.String) {
	emlFile, err := emlparse.ReadFromFile(filePath)

	if err != nil {
		launchErrorDialog(err, window)
		return
	}

	keyHeaders := []string{emlFrom, emlReturnPath, emlTo, emlDate, subject, emlMessageID, emlContentType}

	var analysis string

	// Print values
	for _, fieldName := range keyHeaders {
		field := emlFile.Message.Header.Get(fieldName)
		analysis = fmt.Sprintf("%s\n%s: %q", analysis, fieldName, field)
	}

	displayText.Set(analysis)
}
