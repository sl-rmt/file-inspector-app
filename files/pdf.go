package files

import (
	"bytes"
	"fmt"
	"log"

	"github.com/RedMapleTech/pdf-parse/pdf"
)

func processPDFFile(result *ProcessResult) {
	var analysis, metadata bytes.Buffer

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

	if err != nil {
		result.Completed = false
		result.Error = err
		return
	}

	// w := new(tabwriter.Writer)
	// w.Init(&analysis, 4, 2, 0, '\t', 0)

	for field, value := range md {
		metadata.WriteString(fmt.Sprintf("%s\t%s\n", field, value))
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
		analysis.WriteString("None found")
	} else {
		analysis.WriteString(activeResult)
		result.Dangerous = true
	}

	log.Println("PDF processing done")
	result.Analysis = analysis.String()
	result.Metadata = metadata.String()
	result.Completed = true
}
