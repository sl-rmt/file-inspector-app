package emlparse

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/mail"
	"os"

	"file-inspector/emails/msgparse"
)

type Eml struct {
	Message     *mail.Message
	Body        string
	Attachments []msgparse.Attachment
}

func ReadFromFile(filePath string) (*Eml, error) {
	var emlFile Eml

	// read the contents and parse them
	file, err := os.Open(filePath)

	if err != nil {
		return nil, fmt.Errorf("error opening file: %s", err.Error())
	}

	email, err := mail.ReadMessage(file)

	if err != nil {
		return nil, fmt.Errorf("error reading eml message: %s", err.Error())
	}

	// parsed fine so put it into the struct
	emlFile.Message = email

	// Accessing the body reader won't work after we close the file so 
	// read the body bytes out to a buffer then store them
	buf := new(bytes.Buffer)
	numRead, err := buf.ReadFrom(email.Body)

	if err != nil && err != io.EOF {
		log.Printf("error reading body bytes: %s", err.Error())
	} else if numRead == 0 {
		log.Printf("failed to read any body bytes")
	}

	// parsed fine so put it into the struct
	emlFile.Body = buf.String()

	// get any attachments out if there's body content
	if len(emlFile.Body) > 0 {
		attachments, err := extractAllAttachments(email, emlFile.Body)

		if err != nil && err.Error() != NoAttachments {
			return nil, err
		}
	
		// parsed fine so put it into the struct
		emlFile.Attachments = attachments
	}

	// done processing so close the file
	file.Close()

	return &emlFile, nil
}
