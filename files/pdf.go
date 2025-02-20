package files

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"file-inspector/files/pdf"
)

func processPDFFile(result *ProcessResult) {
	var analysis bytes.Buffer
	var metadata [][]string

	// check encryption
	encrypted, err := pdf.IsEncrypted(result.FilePath)

	if err != nil {
		result.Parsed = false
		result.Completed = false
		result.Error = err
		return
	}

	if encrypted {
		result.Analysis = fmt.Sprintf("File %q is encrypted and password protected, so cannot be inspected.\n", result.FilePath)
		result.Completed = false
		result.Dangerous = true
		return
	}

	// get metadata
	md, err := pdf.GetMetadata(result.FilePath)

	// don't care if errors because there's no metadata
	if err != nil && !strings.Contains(err.Error(), "Failed to get any metadata") {
		result.Completed = false
		result.Error = err
		return
	}

	if len(md) > 0 {
		for field, value := range md {
			metadata = append(metadata, []string{strings.TrimSpace(field), value})
		}
	}

	// check for active content
	activeResult, err := pdf.CheckForActiveContent(result.FilePath)

	if err != nil {
		result.Completed = false
		result.Error = err
		return
	}

	analysis.WriteString("Active content in file:\n\n")

	if len(activeResult) == 0 {
		analysis.WriteString("\t✅ None found")
	} else {
		analysis.WriteString(activeResult)
		result.Dangerous = true
	}

	log.Println("PDF processing done")
	result.Analysis = analysis.String()
	result.Metadata = metadata
	result.Completed = true
}
