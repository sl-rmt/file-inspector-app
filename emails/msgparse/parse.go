package msgparse

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/richardlehane/mscfb"
)

func ReadMsgFile(filePath string, verbose bool) (*Message, error) {

	// open the file
	f, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	// parse it as an OLE doc
	doc, err := mscfb.New(f)

	if err != nil {
		return nil, err
	}

	// create the output and its maps
	msg := &Message{}
	msg.Properties = make(map[string]string)
	msg.UnknownProperties = make(map[int64]UnknownProperty)

	// extract the message content
	processDocEntries(doc, msg, verbose)

	return msg, nil
}

// Process each entry successively
func processDocEntries(doc *mscfb.Reader, msg *Message, verbose bool) {
	for {
		// get the next entry from the doc
		entry, err := doc.Next()

		// if there are no more entries
		if err != nil {
			// done processing
			break
		}

		// for any of the attachment prefixes
		if strings.Contains(entry.Name, attachmentPrefix) {
			var currentAttachment Attachment

			// add this first entry
			err := addEntryToAttachment(entry, &currentAttachment)

			if err != nil {
				log.Printf("\tError processing attachment entry: %s\n", err.Error())
			}

			// attachment entries are all sequential, so
			// collate subsequent entries related to attachments
			for {
				// get the next entry from the doc
				entry, err = doc.Next()

				// if there are no more entries or it's not attachment related
				if err != nil || !strings.HasPrefix(entry.Name, attachmentPrefix) {
					// done processing this attachment, save it and break
					msg.Attachments = append(msg.Attachments, currentAttachment)

					// TODO add this in as an option
					//DumpBinaryAttachment(currentAttachment)

					// reset it
					currentAttachment = Attachment{}

					break
				} else {
					err = addEntryToAttachment(entry, &currentAttachment)

					if err != nil {
						log.Printf("\tError processing attachment entry: %s\n", err.Error())
					}
				}
			}
		} else if strings.Contains(entry.Name, propertyStreamPrefix) {
			// for other properties
			prop, err := extractEntryProperty(entry)

			if err != nil {
				// print them if verbose
				if verbose {
					log.Printf("\tError parsing property %q from stream: %s\n", entry.Name, err.Error())
				}
			} else {
				err = msg.addPropertyToMessage(*prop, verbose)

				if err != nil {
					log.Printf("\tError adding property %q to message: %s\n", entry.Name, err.Error())
				}
			}
		} else {
			log.Printf("\tUnknown entry type: %q", entry.Name)
		}
	}
}

func extractEntryProperty(entry *mscfb.File) (*EntryProperty, error) {
	properties, err := determineEntryProperties(entry)

	if err != nil {
		return nil, err
	}

	data, err := decodeDataFromProperty(entry, *properties)

	if err != nil {
		return nil, err
	}

	// build the struct and return it
	messageProperty := EntryProperty{
		PropertyType: properties.PropertyType,
		Encoding:     properties.Encoding,
		Data:         data,
	}

	return &messageProperty, nil
}

func decodeDataFromProperty(entry *mscfb.File, info EntryProperty) (interface{}, error) {

	if info.PropertyType == "" {
		return nil, fmt.Errorf("empty property type")
	}

	switch info.Encoding {
	// ASCII
	case AsciiEncoding:
		rawBytes := make([]byte, entry.Size)
		entry.Read(rawBytes)

		decoded, err := decodeACSII(rawBytes)

		if err != nil {
			return nil, fmt.Errorf("error decoding ASCII: %s", err.Error())
		}

		return decoded, nil
	// UNICODE
	case UnicodeEncoding:
		rawBytes := make([]byte, entry.Size)
		entry.Read(rawBytes)

		decoded, err := decodeUTF16LE(rawBytes)

		if err != nil {
			return nil, fmt.Errorf("error decoding Unicode: %s", err.Error())
		}

		return decoded, nil
	// Binary
	case BinaryEncoding:
		rawBytes := make([]byte, entry.Size)
		entry.Read(rawBytes)
		return rawBytes, nil
	// Other
	default:
		typeName, err := GetEncodingName(info.Encoding)

		if err == nil {
			log.Printf("\tFound unknown field of type %s, ID: 0x%s\n", typeName, info.PropertyType)
		} else {
			log.Printf("\tFound unknown field of unknown type %s, ID: 0x%s\n", info.Encoding, info.PropertyType)
		}

		rawBytes := make([]byte, entry.Size)
		entry.Read(rawBytes)
		return rawBytes, nil
	}
}

// See https://learn.microsoft.com/en-us/office/client-developer/outlook/mapi/mapi-property-tags
func determineEntryProperties(entry *mscfb.File) (*EntryProperty, error) {
	name := entry.Name
	prop := EntryProperty{}

	// if it's a property stream i.e. entry "__substg1.0_00020102" it has property type: 0002, encoding: 0x0102
	if strings.HasPrefix(name, propertyStreamPrefix) {
		var property string
		var encoding string

		val := name[len(propertyStreamPrefix):]
		property = val[0:4]
		encoding = val[4:8]

		prop = EntryProperty{
			PropertyType: property,
			Encoding:     encoding,
		}

		return &prop, nil
	}

	return nil, fmt.Errorf("stream has the wrong prefix")
}

// Tries a few time encodings
// TODO replace with ParseDate in https://go.dev/src/net/mail/message.go?
func GetTimeFromString(s string, propertyID int64) (time.Time, error) {

	var t time.Time

	if s == "" {
		return t, fmt.Errorf("error converting time from property %x: Empty string provided", propertyID)
	}

	// try RFC1123
	t, err := time.Parse(time.RFC1123, s)

	if err == nil {
		return t, nil
	}

	// try variant of RFC1123
	t, err = time.Parse("02 Jan 2006 15:04:05 (MST)", s)

	if err == nil {
		return t, nil
	}

	// e.g. 31 Jul 2023 07:49:56.7052 (UTC)
	t, err = time.Parse("02 Jan 2006 15:04:05.0000 (MST)", s)

	if err == nil {
		return t, nil
	}

	// try another variant of RFC1123
	// 09 Aug 2023 04:31:38
	t, err = time.Parse("02 Jan 2006 15:04:05", s)

	if err == nil {
		return t, nil
	}

	return t, fmt.Errorf("error parsing time from property %x %q: %s", propertyID, s, err.Error())
}

func GetEncodingName(s string) (string, error) {
	switch s {
	case AsciiEncoding:
		return "ASCII", nil
	case UnicodeEncoding:
		return "Unicode", nil
	case BinaryEncoding:
		return "Binary", nil
	case PT_NULL:
		return "PT_NULL", nil // null value
	case PT_SHORT:
		return "PT_SHORT", nil // signed 16 bit value
	case PT_LONG:
		return "PT_LONG", nil // signed or unsigned 32 bit value
	case PT_FLOAT:
		return "PT_FLOAT", nil // 32 bit floating point
	case PT_DOUBLE:
		return "PT_DOUBLE", nil // 64 bit floating point
	case PT_CURRENCY:
		return "PT_CURRENCY", nil // currency (64 bit integer)
	case PT_APPTIME:
		return "PT_APPTIME", nil // date type
	case PT_ERROR:
		return "PT_ERROR", nil // 32 bit error value
	case PT_BOOLEAN:
		return "PT_BOOLEAN", nil // boolean
	case PT_OBJECT:
		return "PT_OBJECT", nil // embedded object
	case PT_LONGLONG:
		return "PT_LONGLONG", nil // 64 bit signed integer
	case PT_SYSTIME:
		return "PT_SYSTIME", nil // date type
	case OLEGUID:
		return "OLEGUID", nil // OLE GUID
	}

	return "", fmt.Errorf("unknown type ID: 0x%s", s)
}
