package msgparse

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/richardlehane/mscfb"
)

// Process an attachment entry and add the content to the passed attachment instance
func addEntryToAttachment(entry *mscfb.File, attachment *Attachment) error {
	switch entry.Name {
	case attachmentName:
		rawBytes := make([]byte, entry.Size)
		entry.Read(rawBytes)

		decoded, err := decodeUTF16LE(rawBytes)

		if err != nil {
			return fmt.Errorf("error decoding unicode from AttachmentName: %s", err.Error())
		} else {
			attachment.Filename = string(decoded)
		}
	case attachmentLongName:
		rawBytes := make([]byte, entry.Size)
		entry.Read(rawBytes)

		decoded, err := decodeUTF16LE(rawBytes)

		if err != nil {
			return fmt.Errorf("error decoding unicode from AttachmentLongName: %s", err.Error())
		} else {
			attachment.LongFilename = string(decoded)
		}
	case attachmentUnicodeExtension:
		rawBytes := make([]byte, entry.Size)
		entry.Read(rawBytes)

		decoded, err := decodeUTF16LE(rawBytes)

		if err != nil {
			return fmt.Errorf("error decoding unicode from AttachmentUnicodeExtension: %s", err.Error())
		} else {
			attachment.UnicodeExtension = string(decoded)
		}
	case attachmentMimeTag:
		rawBytes := make([]byte, entry.Size)
		entry.Read(rawBytes)

		decoded, err := decodeUTF16LE(rawBytes)

		if err != nil {
			return fmt.Errorf("error decoding unicode from AttachmentMimeTag: %s", err.Error())
		} else {
			attachment.MimeTag = string(decoded)
		}
	case attachmentFolder:
		// don't care here
	case attachmentData:
		bytes := make([]byte, entry.Size)
		read, err := entry.Read(bytes)

		if err != nil {
			return fmt.Errorf("error reading bytes from entry: %s", err.Error())
		} else if read != int(entry.Size) {
			return fmt.Errorf("read %d bytes from entry, not %d", read, entry.Size)
		}

		// lots of empty entries, for some reason
		if entry.Size > 0 {
			attachment.Bytes = bytes
		}
	case attachmentOtherBinData1:
		fallthrough
	case attachmentOtherBinData2:
		fallthrough
	case attachmentOtherBinData3:
		fallthrough
	case attachmentOtherBinData4:
		bytes := make([]byte, entry.Size)
		read, err := entry.Read(bytes)

		if err != nil && err != io.EOF {
			return fmt.Errorf("error reading bytes from entry: %s", err.Error())
		} else if read != int(entry.Size) {
			return fmt.Errorf("read %d bytes from entry, not %d: %s", read, entry.Size, err.Error())
		}

		// lots of empty entries, for some reason
		if entry.Size > 0 {
			attachment.OtherData = bytes
		}
	default:
		// unknown
		log.Printf("Unknown attachment entry %q. Skipping it.\n", entry.Name)
	}

	return nil
}

// bit risky if we don't trust the attachments
func DumpBinaryAttachment(attachment Attachment) {
	filename := fmt.Sprintf("att_%s", attachment.LongFilename)
	out, err := os.Create(filename)

	if err != nil {
		log.Fatalf(err.Error())
	}

	written, err := out.Write(attachment.Bytes)

	if written != len(attachment.Bytes) {
		log.Printf("Only wrote %d bytes, not full amount of  %d", written, len(attachment.Bytes))
	} else if written != attachment.Size {
		log.Printf("Only wrote %d bytes, not expected size of %d", written, attachment.Size)
	} else if err != nil {
		log.Fatalf(err.Error())
	}

	out.Close()
	log.Printf("Attachment data written to %s\n", filename)
}
