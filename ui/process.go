package main

import (
	"fmt"
	"log"
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
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

func getProps(filePath string, fileName, hash, fileType, size binding.String, window *fyne.Window) {
	// set file name
	fileName.Set(fmt.Sprintf("Filename: %q", filePath))

	// get details
	details, err := details.GetFileDetails(filePath)

	// set them if no error
	if err != nil {
		launchErrorDialog(err, window)
	} else {
		fileType.Set(fmt.Sprintf("File Type: %q", details.Mimetype))
		hash.Set(fmt.Sprintf("SHA256: %q", details.SHA256))
		size.Set(fmt.Sprintf("Size: %q", details.SizeString))
	}
}

func processFile(filePath string, window *fyne.Window, text binding.String) {
	log.Printf("Processing file %q\n", filePath)

	fileExt := path.Ext(filePath)

	switch fileExt {
	case ".msg":
		processMsgFile(filePath, window, text)
	case ".eml":
		log.Println("Parsing email file")
		processEmlFile(filePath, window)
	case ".docx":
		log.Println("Parsing document file")
	default:
		launchInfoDialog("Unknown file", fmt.Sprintf("Can only inspect known file types, not %s", fileExt), window)
	}
}

func processMsgFile(filePath string, window *fyne.Window, text binding.String) {
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

	text.Set(analysis)
}

func processEmlFile(filePath string, window *fyne.Window) {

}
