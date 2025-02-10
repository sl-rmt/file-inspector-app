package pdf

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"seehuhn.de/go/pdf"
)

func IsEncrypted(filePath string) (bool, error) {
	_, err := getReader(filePath)

	if err != nil {
		if strings.Contains(err.Error(), "authentication failed for document") {
			return true, nil
		}

		return false, err
	}

	return false, nil
}

func CheckForActiveContent(filePath string) (string, error) {

	fd, err := os.Open(filePath)

	if err != nil {
		return "❓ Failed to complete checks", err
	}

	reader, err := getReader(filePath)

	if err != nil {
		return "❓Failed to complete checks", err
	}

	info, err := pdf.SequentialScan(fd)

	if err != nil {
		return "❓Failed to complete checks", err
	}

	keywords := [][]string{
		{"/JavaScript","Javascript content is an embedded script that can run when the document is opened"}, // "<<\n/EmbeddedFiles 243 0 R\n/JavaScript 251 0 R\n>>"
		{"/AcroForm","Active content use to build an editable form"},   // "<<\n/AcroForm 249 0 R\n/Metadata 245 0 R\n/Names 250 0 R\n/Outlines 176 0 R\n/
		{"/JS","Javascript aka 'JS' content is an embedded script that can run when the document is opened"},
		{"/OpenAction","An active action that is designed to run when the PDF is opened"},
		{"/Launch","An active action that is designed to run when the PDF is opened"},
		{"/AA","An active action 'AA' that is designed to run when the PDF is opened"},
	}

	var result bytes.Buffer
	counts := make([]int, len(keywords))

	for _, section := range info.Sections {
		for _, fileObject := range section.Objects {
			n := fileObject.Reference.Number()

			if fileObject.Broken {
				continue
			}

			//fmt.Printf("%d: %q %s, %s (%d->%d)\n", n, fileObject.String(), fileObject.Type, fileObject.SubType, fileObject.ObjStart, fileObject.ObjEnd)

			object, err := reader.Get(fileObject.Reference, true)

			if err != nil {
				result.WriteString(fmt.Sprintf("❓Failed to get obj: %s", err.Error()))
			}

			var buf bytes.Buffer
			writer := bufio.NewWriter(&buf)
			err = object.PDF(writer)

			if err != nil {
				result.WriteString(fmt.Sprintf("❓Failed to write obj: %s", err.Error()))
			}

			// need to flush the writer to get the bytes in to the buffer
			writer.Flush()
			objectBody := buf.String()

			// the body header sits between << and >>, if present. E.g.:
			// "<<\n/B 1709\n/Filter /FlateDecode\n/I 1733\n/L 1693\n/Length 1173\n/O 1297\n/S 613\n/V 1313\n>>\nstream\nx
			if strings.Contains(objectBody, "<<") && strings.Contains(objectBody, ">>") {

				header := objectBody[strings.Index(objectBody, "<<"):strings.Index(objectBody, ">>")]

				// replace ASCII hex equivalents
				if strings.Contains(header, "#") {
					header, err = replaceEscaped(header)

					if err != nil {
						result.WriteString(fmt.Sprintf("Error decoding body for object %d: %s\n", n, err.Error()))
					}
				}

				// Check all the known keywords
				for i, keyword := range keywords {
					if strings.Contains(header, keyword[0]) {
						//log.Printf("Found %q in object %d\n", keyword, n)
						counts[i]++
					}
				}
			}

		}
	}

	for i, c := range counts {
		if c > 0 {
			result.WriteString(fmt.Sprintf("☠️ Found %d references to the active content %q in the file's objects. %s.\n\n", c, keywords[i][0], keywords[i][1]))
		}
	}

	return result.String(), nil
}

// Decode ASCII hex obfuscation, e.g.
// Bad: "/#4Aava#53cript". Fixed: "/JavaScript"
// Bad: "/J#61vaScrip#74". Fixed: "/JavaScript"
// Bad: "/#4a#61#76#61#53#63#72#69#70#74". Fixed: "/JavaScript"
func replaceEscaped(objectBody string) (string, error) {
	for {
		// keep replacing them until we have none left
		if !strings.Contains(objectBody, "#") {
			break
		}

		// start of three chars e.g. #4A
		index := strings.Index(objectBody, "#")

		// e.g. 4A
		asciiChars := objectBody[index+1 : index+3]

		// decode it
		ascii, err := hex.DecodeString(asciiChars)

		if err == nil {
			objectBody = strings.ReplaceAll(objectBody, objectBody[index:index+3], string(ascii))

			// do we care if it's an error? most likely an # in some other body part
			//return objectBody, fmt.Errorf("Failed to convert %q: %s", asciiChars, err.Error())
		} else {
			break
		}

	}

	return objectBody, nil
}
