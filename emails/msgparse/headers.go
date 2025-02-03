package msgparse

import (
	"bufio"
	"io"
	"net/textproto"
	"strings"
)

func GetHeaderByName(headers, name string) (string, error) {
	reader := strings.NewReader(headers)
	tpReader := textproto.NewReader(bufio.NewReader(reader))

	hdr, err := tpReader.ReadMIMEHeader()

	if err != nil && err != io.EOF {
		return "", err
	}

	auth := hdr.Get(name)

	return auth, nil
}
