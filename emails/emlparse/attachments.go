package emlparse

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"net/mail"
	"strings"

	"mime"

	"file-inspector/emails/msgparse"
)

const (
	NoAttachments = "content type is not multipart"
)

func extractAllAttachments(message *mail.Message, bodyString string) ([]msgparse.Attachment, error) {
	
	// get the details from teh message header
	mediaType, params, err := mime.ParseMediaType(message.Header.Get("Content-Type"))

	if err != nil {
		return nil, fmt.Errorf("error parsing media type: %s", err.Error())
	}

	// no attachments
	if !strings.HasPrefix(mediaType, "multipart/mixed") {
		return nil, fmt.Errorf(NoAttachments)
	}

	// pull out the attachments
	attachments, err := extractAttachmentsFromBoundary(bodyString, params["boundary"])

	if err != nil {
		return nil, err
	}

	return attachments, nil
}

// Bit hacky, could instead use multipart.NewReader()
func extractAttachmentsFromBoundary(bodyString string, boundary string) ([]msgparse.Attachment, error) {
	var attachments []msgparse.Attachment

	lines := strings.Split(bodyString, "\n")

	var i int

	// run through each line looking for the boundary
	for i = 0; i < len(lines); i++ {
		line := lines[i]

		// looks like this:
		// --=-XNI3F2P8aCdwwxXQDLdRmw==   												<--- Boundary
		// Content-Type: application/octet-stream; name=G026730897.pdf
		// Content-Disposition: attachment; filename=G026730897.pdf
		// <optional extra fields>
		// Content-Transfer-Encoding: base64
		//
		// JVBERi0xLjcKJeLjz9MKMTEgMCBvYmoKPDwvQSAxMiAwIFIvQm9yZGVyWzAgMCAwXS9GIDQvUCA0 <--- Base64 blob
		// IDAgUi9SZWN0WzM4OC43NSA1NzEuMTYgNDY1Ljg2IDU3OS4xNF0vU3VidHlwZS9MaW5rPj4KZW5k
		// ...
		// ZWYKMTI1MzgwCiUlRU9GCg==
		//
		// --=-XNI3F2P8aCdwwxXQDLdRmw==--   											<--- Boundary again
		if strings.Contains(line, boundary) {
			const ct = "Content-Type: "
			const name = "name="
			const b64 = "Content-Transfer-Encoding: base64"

			if strings.HasPrefix(lines[i+1], ct) && strings.Contains(lines[i+1], name) {
				filename := lines[i+1][strings.Index(lines[i+1], name)+5:]
				
				// trim whitespace
				filename = strings.TrimSpace(filename)

				// trim quotes
				filename = strings.Trim(filename, "\"")

				// check the next few lines for the base64 tag
				for j := 1; j < 10; j++ {

					// if we find it
					if strings.HasPrefix(lines[i+j], b64) {

						// send off the rest of the lines to parse it out - start two off to account for the empty line
						// TODO could use the location of the next boundary instead
						rawBytes, err := getAttachmentBytes(lines[i+j+2:], boundary)

						if err != nil {
							log.Printf("Error extracting attachment %q: %s", filename, err.Error())
							return nil, err
						} else {
							log.Printf("Extracted %d bytes for attachment: %q\n", len(rawBytes), filename)
							attachments = append(attachments, msgparse.Attachment{Bytes: rawBytes, Filename: filename})
						}
					}
				}
			}
		}
	}

	return attachments, nil
}

func getAttachmentBytes(lines []string, boundary string) ([]byte, error) {
	var buf bytes.Buffer

	// stitch all the lines together
	for _, s := range lines {
		
		// we're done when we get to an empty line or the boundary
		if s == "" || strings.Contains(s, boundary) {
			break
		}

		trimmed := strings.TrimSpace(s)
		
		if trimmed != "" {
			buf.WriteString(trimmed)
		}
	}

	trimmed := strings.TrimSpace(buf.String())

	// decode Base64
	decoded, err := base64.StdEncoding.DecodeString(trimmed)

	if err != nil {
		return nil, err
	}

	return decoded, nil
}
