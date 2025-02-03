package msgparse

const (
	// other types - https://www.dimastr.com/redemption/utils.htm
	PT_NULL     = "0001" // null value
	PT_SHORT    = "0002" // signed 16 bit value
	PT_LONG     = "0003" // signed or unsigned 32 bit value
	PT_FLOAT    = "0004" // 32 bit floating point
	PT_DOUBLE   = "0005" // 64 bit floating point
	PT_CURRENCY = "0006" // currency (64 bit integer)
	PT_APPTIME  = "0007" // date type
	PT_ERROR    = "000A" // 32 bit error value
	PT_BOOLEAN  = "000B" // boolean
	PT_OBJECT   = "000D" // embedded object
	PT_LONGLONG = "0014" // 64 bit signed integer
	PT_SYSTIME  = "0040" // date type
	OLEGUID     = "0048" // OLE GUID

	// type we've seen
	AsciiEncoding   = "001E" // aka 8 bit string
	UnicodeEncoding = "001F"
	BinaryEncoding  = "0102"

	propertyStreamPrefix = "__substg1.0_"
	attachmentPrefix     = "__substg1.0_37"
	attachmentData             = "__substg1.0_37010102"
	attachmentUnicodeExtension = "__substg1.0_3703001F"
	attachmentFolder           = "__substg1.0_3701000D"
	attachmentName             = "__substg1.0_3704001F"
	attachmentLongName         = "__substg1.0_3707001F"
	attachmentMimeTag          = "__substg1.0_370E001F"

	attachmentOtherBinData1 = "__substg1.0_37020102"
	attachmentOtherBinData2 = "__substg1.0_371D0102"
	attachmentOtherBinData3 = "__substg1.0_370A0102"
	attachmentOtherBinData4 = "__substg1.0_37090102"
)

type Message struct {
	Properties        map[string]string
	UnknownProperties map[int64]UnknownProperty
	Attachments       []Attachment
}

// EntryProperty holds the type of data and the data itself
type EntryProperty struct {
	PropertyType string
	Encoding     string
	Data         interface{}
}

type UnknownProperty struct {
	PropertyType string
	Encoding     string
	Data         string
}

type Attachment struct {
	Bytes            []byte
	OtherData        []byte
	Size             int
	Filename         string
	LongFilename     string
	MimeTag          string
	UnicodeExtension string
}
