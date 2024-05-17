package files

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"fyne.io/fyne/v2/data/binding"

	"mvdan.cc/xurls/v2"

	"github.com/RedMapleTech/email-parse/emlparse"
	"github.com/RedMapleTech/email-parse/msgparse"

	"github.com/RedMapleTech/forensics-email-inspector/safelinks"
	"github.com/RedMapleTech/url-inspect/urls"
)

func processMsgFile(filePath string, displayText binding.String) (bool, error) {
	msg, err := msgparse.ReadMsgFile(filePath, false)

	if err != nil {
		return true, err
	}

	log.Println("Email parsing done")

	// print key fields
	keyFieldNames := []string{msgSender, msgDisplayName, msgSenderSMTP, msgSenderEmail, msgSenderEmail2, msgReceivedName, msgReceivedSMTP, msg7bitEmail, msgReceivedEmail, subject, messageTopic, msgMessageID}
	var analysis bytes.Buffer

	// Print values
	for _, fieldName := range keyFieldNames {
		field := msg.GetPropertyByName(fieldName)

		if len(field) > 0 {
			analysis.WriteString(fmt.Sprintf("%s: %q\n", fieldName, field))
		}
	}

	// set the analysis text in the bound UI element
	displayText.Set(analysis.String())

	// add details on authentication
	authHeader, err := msgparse.GetHeaderByName(msg.GetPropertyByName("Message Headers"), authResults)
	dangerous := false

	if err != nil {
		return true, err
	} else if authHeader != "" {
		dangerous = parseAuthResults(authHeader, &analysis)
		analysis.WriteString("\n")

		// set the analysis text in the bound UI element
		displayText.Set(analysis.String())
	}

	// body details
	inspectBody(msg.GetPropertyByName("Message body"), &analysis)

	// set the analysis text in the bound UI element
	displayText.Set(analysis.String())

	// add attachment details, if there are any
	if len(msg.Attachments) > 0 {
		addAttachmentDetails(msg.Attachments, &analysis)
	}

	// set the analysis text in the bound UI element
	displayText.Set(analysis.String())

	log.Println("Msg processing done")
	return dangerous, nil
}

// return true if dangerous
func parseAuthResults(authHeader string, buffer *bytes.Buffer) bool {
	fields := strings.Split(authHeader, ";")

	if len(fields) == 0 {
		return true
	}

	buffer.WriteString("\nAuthentication results:\n")
	dangerous := false

	for _, field := range fields {
		field = strings.TrimSpace(field)

		if strings.HasPrefix(field, "dkim=") || strings.HasPrefix(field, "spf=") || strings.HasPrefix(field, "dmarc=") {
			if strings.Contains(field, "=pass") {
				buffer.WriteString(fmt.Sprintf("\tGOOD: %s\n", field))
			} else {
				buffer.WriteString(fmt.Sprintf("\tBAD: %s\n", field))
				dangerous = true
			}
		}
	}

	log.Println("Auth processing done")
	return dangerous
}

// return true if dangerous
func processEmlFile(filePath string, displayText binding.String) (bool, error) {
	emlFile, err := emlparse.ReadFromFile(filePath)

	if err != nil {
		return true, err
	}

	keyHeaders := []string{emlFrom, emlReturnPath, emlTo, emlDate, subject, emlMessageID, emlContentType}
	var analysis bytes.Buffer

	// Print values
	for _, fieldName := range keyHeaders {
		field := emlFile.Message.Header.Get(fieldName)

		if len(field) > 0 {
			analysis.WriteString(fmt.Sprintf("%s: %q\n", fieldName, field))
		}
	}

	// set the analysis text in the bound UI element
	displayText.Set(analysis.String())

	// get the auth results and parse them
	authHeader := (emlFile.Message.Header.Get(authResults))
	dangerous := false

	if authHeader != "" {
		dangerous = parseAuthResults(authHeader, &analysis)
	}

	// set the analysis text in the bound UI element
	displayText.Set(analysis.String())

	// add attachment details, if there are any
	if len(emlFile.Attachments) > 0 {
		addAttachmentDetails(emlFile.Attachments, &analysis)

		// set the analysis text in the bound UI element
		displayText.Set(analysis.String())
	}

	// body details
	inspectBody(emlFile.Body, &analysis)
	displayText.Set(analysis.String())
	log.Println("Eml processing done")

	return dangerous, nil
}

func addAttachmentDetails(attachments []msgparse.Attachment, analysis *bytes.Buffer) {
	log.Println("Parsing attachments")
	analysis.WriteString(fmt.Sprintf("\nEmail has %d attachments:\n", len(attachments)))

	for i, a := range attachments {
		analysis.WriteString(fmt.Sprintf("\tAttachment %d:\n", i+1))

		if len(a.Filename) > 0 {
			analysis.WriteString(fmt.Sprintf("\tFilename: %q\n", a.Filename))
		}

		if len(a.LongFilename) > 0 {
			analysis.WriteString(fmt.Sprintf("\tLong Filename: %q\n", a.LongFilename))
		}

		if len(a.MimeTag) > 0 {
			analysis.WriteString(fmt.Sprintf("\tMIME tag: %q\n", a.MimeTag))
		}

		analysis.WriteString(fmt.Sprintf("\tSize: %d bytes\n", len(a.Bytes)))

		hash := sha256.New()
		hash.Write(a.Bytes)
		analysis.WriteString(fmt.Sprintf("\tSHA-256 hash: %q\n\n", hex.EncodeToString(hash.Sum(nil))))
	}
}

func inspectBody(body string, analysis *bytes.Buffer) {
	log.Println("Inspecting email body")
	analysis.WriteString("\nBody Details:\n")

	if len(body) == 0 {
		analysis.WriteString("\tEmpty body\n")
		return
	} else {
		// count lines
		// TODO this is lame - detect html vs plaintext and do something smarter
		// lines := strings.Split(body, "\r")
		// fmt.Printf("\tBody has %d lines of content\n", len(lines))

		analysis.WriteString("\tEmail body has content.\n")
	}

	err := inspectLinks(body, analysis)

	if err != nil {
		analysis.WriteString(fmt.Sprintf("\tError inspecting body for links: %s.", err.Error()))
	}

	analysis.WriteString("\n")
}

func inspectLinks(body string, analysis *bytes.Buffer) error {
	log.Println("Looking for links")

	// find all URLs in the body
	rxStrict := xurls.Strict()
	res := rxStrict.FindAllString(body, -1)

	log.Printf("Inspecting %d links\n", len(res))

	// if we found any, process them one by one
	if len(res) > 0 {
		analysis.WriteString(fmt.Sprintf("\n\tFound %d URLs in the email body:\n", len(res)))

		// check Alexa common 100k domains
		log.Println("Loading common URLs")
		commonChecker, err := urls.GetCommonURLChecker()

		if err != nil {
			return err
		}

		log.Printf("Loaded %d common Alexa domains\n", commonChecker.CountKnownDomains())

		for _, entry := range res {
			entry = strings.ToLower(entry)

			// skip email links, empty strings
			if strings.HasPrefix(entry, "mailto:") {
				continue
			} else if strings.HasPrefix(entry, "tel:") {
				continue
			}

			if len(strings.TrimSpace(entry)) == 0 {
				continue
			}

			// check if it's an Outlook safelink
			if safelinks.IsSafelink(entry) {
				original, err := safelinks.ExtractOriginalURL(entry)

				if err != nil {
					analysis.WriteString(fmt.Sprintf("\t\tError extracting URL: %s\n", err.Error()))
				} else {
					// check if it's from a common (most popular 100k) domain
					isCommon, err := commonChecker.Check(entry)

					if err != nil {
						return err
					}

					if isCommon {
						analysis.WriteString(fmt.Sprintf("\t\tSafelink redirects to common domain %q\n", original))
					} else {
						analysis.WriteString(fmt.Sprintf("\t\tSafelink redirects to *uncommon* domain %q\n", original))
					}
				}
			} else {
				// not a safelink
				// check if it's from a common (most popular 100k) domain
				isCommon, err := commonChecker.Check(entry)

				if err != nil {
					return err
				}

				if isCommon {
					analysis.WriteString(fmt.Sprintf("\t\tURL from common domain: %q\n", entry))
				} else {
					analysis.WriteString(fmt.Sprintf("\t\tURL from *uncommon* domain: %q\n", entry))
				}
			}
		}
	}

	return nil
}
