package pdf

import (
	"fmt"
	"os"

	"seehuhn.de/go/pdf"
)

func GetMetadata(filePath string) (map[string]string, error) {

	reader, err := getReader(filePath)

	if err != nil {
		return nil, err
	}

	metadata := reader.GetMeta()
	fields := make(map[string]string)
	fields["Version"] = metadata.Version.String()

	// TODO doesn't work
	//fields["Pages"] = metadata.Catalog.Pages.String()

	info := metadata.Info

	if info == nil {
		return nil, fmt.Errorf("Failed to get any metadata")
	}

	if info.Title != "" {
		fields["Title"] = info.Title
	}

	if info.Author != "" {
		fields["Author"] = info.Author
	}

	if info.Subject != "" {
		fields["Subject"] = info.Subject
	}

	if info.Keywords != "" {
		fields["Keywords"] = info.Keywords
	}

	if info.Creator != "" {
		fields["Creator"] = info.Creator
	}

	if info.Producer != "" {
		fields["Producer"] = info.Producer
	}

	if info.CreationDate.IsZero() {
		fields["CreationDate"] = info.CreationDate.String()
	}

	if !info.ModDate.IsZero() {
		fields["Modified Date"] = info.ModDate.String()
	}

	return fields, nil
}

func getReader(filePath string) (*pdf.Reader, error) {
	fd, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	opt := &pdf.ReaderOptions{
		ErrorHandling: pdf.ErrorHandlingReport,
	}

	r, err := pdf.NewReader(fd, opt)

	if err != nil {
		return nil, err
	}

	return r, nil
}
