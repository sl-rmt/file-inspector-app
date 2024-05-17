package files

import (
	"bytes"
	"fmt"

	"fyne.io/fyne/v2/data/binding"
	"github.com/RedMapleTech/pdf-parse/pdf"
)

// return true if dangerous
func processPDFFile(filePath string, displayText binding.String) (bool, error) {

	// check encryption
	encrypted, err := pdf.IsEncrypted(filePath)

	if err != nil {
		displayText.Set(fmt.Sprintf("Error opening file: %s", err.Error()))
		return true, nil
	}

	if encrypted {
		displayText.Set(fmt.Sprintf("File %q is encrypted and password protected\n", filePath))
		return true, nil
	}

	// get metadata
	metadata, err := pdf.GetMetadata(filePath)

	if err != nil {
		return true, err
	}

	var analysis bytes.Buffer
	analysis.WriteString("File Metadata:\n\n")

	// w := new(tabwriter.Writer)
	// w.Init(&analysis, 4, 2, 0, '\t', 0)

	for field, value := range metadata {
		analysis.WriteString(fmt.Sprintf("%s\t%s\n", field, value))
	}

	//w.Flush()
	analysis.WriteString("\n")

	// check for active content
	result, err := pdf.CheckForActiveContent(filePath)

	if err != nil {
		return true, err
	}

	analysis.WriteString("Inspecting for active content:\n\n")
	dangerous := false

	if len(result) == 0 {
		analysis.WriteString("None found")
	} else {
		analysis.WriteString(result)
		dangerous = true
	}

	displayText.Set(analysis.String())

	return dangerous, nil

}
