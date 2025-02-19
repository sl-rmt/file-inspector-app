package files

import (
	"bytes"
	"file-inspector/files/docx"
	"fmt"
	"log"
	"strings"
	"text/tabwriter"
)

func processDocxFile(result *ProcessResult) {
	var _, metadata bytes.Buffer

	// get metadata
	coreProps, customProps, err := docx.GetDocProperties(result.FilePath)

	if err != nil && !strings.Contains(err.Error(), "docProps/custom.xml not found") {
		result.Completed = false
		result.Error = err
		return
	}

	if coreProps != nil {
		coreMap := map[string]string{
			"Title":         coreProps.Title,
			"Creator":       coreProps.Creator,
			"Created":       coreProps.Created,
			"Last Modified": coreProps.LastModified,
			"Keywords":      coreProps.Keywords,
			"Revision":      coreProps.Revision,
			"Category":      coreProps.Category,
			"Language":      coreProps.Language,
		}

		printWithTabs(&metadata, "Core Properties\n", coreMap)
	}

	if customProps != nil {
		customMap := make(map[string]string)
		
		if err == nil {
			for _, prop := range customProps.Properties {
				customMap[prop.Name] = prop.Value
			}
			printWithTabs(&metadata, "Custom Properties\n", customMap)
		}
	}

	log.Println("Docx processing done")
	result.Completed = true

	//result.Analysis = analysis.String()
	result.Metadata = metadata.String()
	
}

// printWithTabs outputs properties using a tab writer
func printWithTabs(metadata *bytes.Buffer, title string, properties map[string]string) {
	w := tabwriter.NewWriter(metadata, 0, 8, 1, '\t', 0)
	
	fmt.Fprintln(w, title)
	
	for key, value := range properties {
		if len(value) > 0 {
			fmt.Fprintf(w, "%s\t%s\n", key, value)
		}
	}
	
	w.Flush()
}
