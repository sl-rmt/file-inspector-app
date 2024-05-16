package files

import (
	"bytes"
	"fmt"

	"fyne.io/fyne/v2/data/binding"
	"github.com/RedMapleTech/pdf-parse/pdf"
)

func processPDFFile(filePath string, displayText binding.String) error {

	// check encryption
	encrypted, err := pdf.IsEncrypted(filePath)

	if err != nil {
		displayText.Set(fmt.Sprintf("Error opening file: %s", err.Error()))
		return nil
	}

	if encrypted {
		displayText.Set(fmt.Sprintf("File %q is encrypted and password protected\n", filePath))
		return nil
	}

	// get metadata
	metadata, err := pdf.GetMetadata(filePath)

	if err != nil {
		return err
	}

	var analysis bytes.Buffer
	analysis.WriteString("File Metadata:\n\n")
	
	for field, value := range metadata {
		analysis.WriteString(fmt.Sprintf("%s\t%s\n", field, value))
	}

	analysis.WriteString("\n")

	// check for active content
	result, err := pdf.CheckForActiveContent(filePath)

	if err != nil {
		return err
	}

	analysis.WriteString("Inspecting for active content:\n\n")

	if len(result) == 0 {
		analysis.WriteString("None found")
	} else {
		analysis.WriteString(result)
	}

	displayText.Set(analysis.String())

	return nil

}
