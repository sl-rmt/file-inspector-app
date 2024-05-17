package files

import (
	"bytes"
	"fmt"
	"log"

	"github.com/RedMapleTech/pdf-parse/pdf"
)

func processPDFFile(result *ProcessResult) {
	var analysis bytes.Buffer

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
	metadata, err := pdf.GetMetadata(result.FilePath)

	if err != nil {
		result.Completed = false
		result.Error = err
		return
	}

	analysis.WriteString("File Metadata:\n\n")

	// w := new(tabwriter.Writer)
	// w.Init(&analysis, 4, 2, 0, '\t', 0)

	for field, value := range metadata {
		analysis.WriteString(fmt.Sprintf("%s\t%s\n", field, value))
	}

	//w.Flush()
	analysis.WriteString("\n")

	// check for active content
	activeResult, err := pdf.CheckForActiveContent(result.FilePath)

	if err != nil {
		result.Completed = false
		result.Error = err
		return
	}

	analysis.WriteString("Inspecting for active content:\n\n")

	if len(activeResult) == 0 {
		analysis.WriteString("None found")
	} else {
		analysis.WriteString(activeResult)
		result.Dangerous = true
	}

	log.Println("PDF processing done")
	result.Analysis = analysis.String()
	result.Completed = true
}
