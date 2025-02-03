package files

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"mvdan.cc/xurls/v2"

	"file-inspector/emails/emlparse"
	"file-inspector/emails/msgparse"

	"file-inspector/utils/safelinks"
	"file-inspector/utils/urls"
)

func processMsgFile(result *ProcessResult) {
	msg, err := msgparse.ReadMsgFile(result.FilePath, false)

	if err != nil {
		result.Parsed = false
		result.Completed = false
		result.Error = err
		return
	}

	log.Println("Email parsing done")
	result.Parsed = true

	// print key fields
	keyFieldNames := []string{msgSender, msgDisplayName, msgSenderSMTP, msgSenderEmail, msgSenderEmail2, msgReceivedName, msgReceivedSMTP, msg7bitEmail, msgReceivedEmail, subject, messageTopic, msgMessageID}
	var analysis, metadata bytes.Buffer

	// Print values
	for _, fieldName := range keyFieldNames {
		field := msg.GetPropertyByName(fieldName)

		if len(field) > 0 {
			metadata.WriteString(fmt.Sprintf("%s: %q\n", fieldName, field))
		}
	}

	// add details on authentication
	authHeader, err := msgparse.GetHeaderByName(msg.GetPropertyByName("Message Headers"), authResults)
	result.Dangerous = false

	if err != nil {
		result.Completed = false
		result.Error = err
		return
	} else if authHeader != "" {
		result.Dangerous = parseAuthResults(authHeader, &metadata)
		metadata.WriteString("\n")
	}

	// body details
	inspectBody(msg.GetPropertyByName("Message body"), &analysis)

	// add attachment details, if there are any
	if len(msg.Attachments) > 0 {
		addAttachmentDetails(msg.Attachments, &analysis)
	}

	log.Println("Msg processing done")
	result.Metadata = metadata.String()
	result.Analysis = analysis.String()
	result.Completed = true
}

// return true if dangerous
func processEmlFile(result *ProcessResult) {
	emlFile, err := emlparse.ReadFromFile(result.FilePath)

	if err != nil {
		result.Parsed = false
		result.Completed = false
		result.Error = err
		return
	}

	keyHeaders := []string{emlFrom, emlReturnPath, emlTo, emlDate, subject, emlMessageID, emlContentType}
	var analysis, metadata bytes.Buffer

	// Print values
	for _, fieldName := range keyHeaders {
		field := emlFile.Message.Header.Get(fieldName)

		if len(field) > 0 {
			metadata.WriteString(fmt.Sprintf("%s: %q\n", fieldName, field))
		}
	}

	// get the auth results and parse them
	authHeader := (emlFile.Message.Header.Get(authResults))
	result.Dangerous = false

	if authHeader != "" {
		result.Dangerous = parseAuthResults(authHeader, &metadata)
	}

	// add attachment details, if there are any
	if len(emlFile.Attachments) > 0 {
		addAttachmentDetails(emlFile.Attachments, &analysis)
	}

	// body details
	inspectBody(emlFile.Body, &analysis)

	log.Println("Eml processing done")
	result.Metadata = metadata.String()
	result.Analysis = analysis.String()
	result.Completed = true
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

// TODO return result
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

// TODO return result
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
