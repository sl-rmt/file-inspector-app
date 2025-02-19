package docx

import (
	"fmt"
	"os"

	reader "github.com/fumiama/go-docx"
)

func parseFile(filePath string) error {
	readFile, err := os.Open(filePath)

	if err != nil {
		return fmt.Errorf("error reading file %q: %s", filePath, err.Error())
	}

	fileinfo, err := readFile.Stat()

	if err != nil {
		panic(err)
	}

	size := fileinfo.Size()
	doc, err := reader.Parse(readFile, size)

	if err != nil {
		panic(err)
	}

	fmt.Println("Plain text:")

	for i, it := range doc.Document.Body.Items {
		switch it.(type) {
		case *reader.Paragraph, *reader.Table: // printable
			fmt.Println(it)
		default:
			fmt.Printf("Item %d has unknown type", i)
		}
	}

	return nil
}
