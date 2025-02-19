package docx

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
)

// CoreProperties represents the core properties XML structure
type CoreProperties struct {
	XMLName      xml.Name `xml:"coreProperties"`
	Title        string   `xml:"title"`
	Subject      string   `xml:"subject"`
	Creator      string   `xml:"creator"`
	Description  string   `xml:"description"`
	Keywords     string   `xml:"keywords"`
	LastModified string   `xml:"modified"`
	Created      string   `xml:"created"`
	Revision     string   `xml:"revision"`
	Category     string   `xml:"category"`
	Language     string   `xml:"language"`
}

// CustomProperty represents a single custom property
type CustomProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

// CustomProperties represents the custom properties XML structure
type CustomProperties struct {
	Properties []CustomProperty `xml:"property"`
}

func GetDocProperties(filePath string) (*CoreProperties, *CustomProperties, error) {
	r, err := zip.OpenReader(filePath)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file for zip reader: %s", err.Error())
	}

	defer r.Close()

	// Core properties
	coreProps, err := extractCoreProperties(&r.Reader)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to extract core properties: %s", err.Error())
	}

	// Custom Properties
	customProps, err := extractCustomProperties(&r.Reader)

	if err != nil {
		return coreProps, nil, fmt.Errorf("error getting custom properties: %s", err)
	}

	return coreProps, customProps, nil
}

// extractXML reads and decodes XML from a .docx ZIP entry
func extractXML(r *zip.Reader, filePath string, v interface{}) error {
	for _, f := range r.File {
		if f.Name == filePath {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()
			return xml.NewDecoder(rc).Decode(v)
		}
	}
	return fmt.Errorf("%s not found", filePath)
}

// extractCoreProperties retrieves core properties
func extractCoreProperties(r *zip.Reader) (*CoreProperties, error) {
	var coreProps CoreProperties
	err := extractXML(r, "docProps/core.xml", &coreProps)
	if err != nil {
		return nil, err
	}
	return &coreProps, nil
}

// extractCustomProperties retrieves custom properties
func extractCustomProperties(r *zip.Reader) (*CustomProperties, error) {
	var customProps CustomProperties
	err := extractXML(r, "docProps/custom.xml", &customProps)
	if err != nil {
		return nil, err
	}
	return &customProps, nil
}
