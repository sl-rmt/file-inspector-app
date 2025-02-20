package files

import (
	"file-inspector/files/docx"
	"log"
	"strings"
)

func processDocxFile(result *ProcessResult) {
	//var analysisText bytes.Buffer
	var metadata [][]string

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

		for key, value := range coreMap {
			if len(value) > 0 {
				metadata = append(metadata, []string{key, value})
			}
		}
	}

	if customProps != nil {
		customMap := make(map[string]string)

		for _, prop := range customProps.Properties {
			customMap[prop.Name] = prop.Value
		}

		for key, value := range customMap {
			if len(value) > 0 {
				metadata = append(metadata, []string{key, value})
			}
		}
	}

	log.Println("Docx processing done")
	result.Completed = true

	//result.Analysis = analysis.String()
	result.Metadata = metadata

}
