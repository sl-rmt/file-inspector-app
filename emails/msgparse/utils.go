package msgparse

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func decodeUTF16LE(rawBytes []byte) (string, error) {
	// Make an transformer that converts MS-Win default to UTF8:
	win16be := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)

	// Make a transformer that is like win16be, but abides by BOM:
	utf16bom := unicode.BOMOverride(win16be.NewDecoder())

	// Make a Reader that uses utf16bom:
	unicodeReader := transform.NewReader(bytes.NewReader(rawBytes), utf16bom)
	decoded, err := ioutil.ReadAll(unicodeReader)

	if err != nil {
		return "", fmt.Errorf("error decoding string: %s", err.Error())
	}

	return string(decoded), nil
}

func decodeACSII(rawBytes []byte) (string, error) {
	read, err := charset.NewReader(bytes.NewReader(rawBytes), "ISO-8859-1")

	if err != nil {
		return "", err
	}

	if read == nil {
		return "", fmt.Errorf("ASCII decoder failed to read any bytes")
	}

	decoded, err := ioutil.ReadAll(read)

	if err != nil {
		return "", fmt.Errorf("ASCII decoder failed to read all bytes")
	}
	
	return string(decoded), nil
}
