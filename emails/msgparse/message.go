package msgparse

import (
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"time"
)

const (
	PropertyUnknown = "unknown"
)

func (m Message) GetPropertyByName(name string) string {
	return m.Properties[name]
}

// Found in the following:
// https://isc.sans.edu/diary/Nested+MSGs+Turtles+All+The+Way+Down/26668
// https://www.fileformat.info/format/outlookmsg/
// https://github.com/libyal/libfmapi/blob/main/documentation/MAPI%20definitions.asciidoc
// https://github.com/DidierStevens/DidierStevensSuite/blob/98c7aa67d1ac92a5ea79b37fa7734b183c16bd64/plugin_msg.py#L28
// https://github.com/echo-devim/pyjacktrick/blob/main/mapi_constants.py
// https://github.com/shaniacht1/content/blob/master/automation-ParseEmailFiles.yml
// https://learn.microsoft.com/en-us/office/client-developer/outlook/mapi/mapi-constants#mapi-mime-conversion-api
// TODO more in https://github.com/libyal/libfmapi/blob/main/documentation/MAPI%20definitions.asciidoc#3-the-property-identifiers

// TODO: also seen a date in 0x800d, 0x802d, 0x8012 and 0x8019
func GetPropertyName(intID int64) string {
	allProps := map[int64]string{
		// 0x0001 – 0x0bff | Message envelope properties (defined by MAPI)
		0x001A: "MessageClass",
		0x0037: "Subject", // https://learn.microsoft.com/en-us/office/client-developer/outlook/mapi/pidtagsubject-canonical-property
		0x003D: "Subject Prefix",
		0x003A: "Report Name", // Canonical Property https://learn.microsoft.com/en-us/office/client-developer/outlook/mapi/pidtagreportname-canonical-property
		0x0040: "Received by name",
		0x0042: "Sent Representing name", // Canonical Property
		0x0044: "Received Representing name",
		0x0045: "Report Entry",
		0x004D: "Org Author Name",
		0x004F: "Reply Recipient Entries", // https://learn.microsoft.com/en-us/office/client-developer/outlook/mapi/pidtagreplyrecipiententries-canonical-property
		0x0050: "Reply Recipient Names",   // https://learn.microsoft.com/en-us/office/client-developer/outlook/mapi/pidtagreplyrecipientnames-canonical-property
		0x005A: "Org Sender Name",
		0x0064: "Sent Representing Address Type",
		0x0065: "Sent Representing email",
		0x0070: "Topic",
		0x0075: "Received by address type",
		0x0076: "Received by email",
		0x0077: "Representing address type",
		0x0078: "Representing email",
		0x007d: "Message Headers", // https://learn.microsoft.com/en-us/office/client-developer/outlook/mapi/pidtagtransportmessageheaders-canonical-property
		0x007F: "TNEF Correlation Key", // https://learn.microsoft.com/en-us/office/client-developer/outlook/mapi/pidtagtnefcorrelationkey-canonical-property

		// 0x0c00 – 0x0dff | Recipient properties (defined by MAPI)
		0x0C1A: "Sender name",    // Canonical Property
		0x0C15: "Recipient Type", // Canonical Property: https://learn.microsoft.com/en-us/office/client-developer/outlook/mapi/pidtagrecipienttype-canonical-property
		0x0C1E: "Sender address type",
		0x0C1F: "Sender Email 2",

		// 0x0e00 – 0x0fff | Non-transmittable message properties (defined by MAPI)
		0x0E02: "Display BCC",
		0x0E03: "Display CC",
		0x0E04: "Display To",
		0x0E05: "Parent Display",
		0x0E06: "Message Delivery Time", // Canonical Property: https://learn.microsoft.com/en-us/office/client-developer/outlook/mapi/pidtagmessagedeliverytime-canonical-property
		0x0E1D: "Subject Normalized",
		0x0E28: "Received account1",
		0x0E29: "Received account2",

		// 0x1000 – 0x2fff | Message content properties (defined by MAPI)
		0x1000: "Message body",
		0x1008: "RTF sync body tag",
		0x1009: "Body RTF",
		0x1013: "Body HTML",
		0x1015: "BodyContentId",
		0x1035: "MessageID",
		0x1046: "Sender Email",

		// 0x3000 - 0x33ff | Common object properties that appear on multiple objects (defined by MAPI)
		0x3001: "Display name",  // Canonical Property
		0x3002: "Address type",  // Canonical Property
		0x3003: "Email address", // Canonical Property
		0x3007: "CreationTime",

		// 0x3400 - 0x35ff | Message store properties (defined by MAPI)

		// 0x3700 – 0x38ff | Attachment properties (defined by MAPI)
		0x3701: "Attachment data",
		0x3703: "Attachment file extension",
		0x3704: "Attachment Filename",
		0x3707: "Attachment long filename",
		0x370E: "Attachment MIME tag",
		0x3712: "Attachment ID",

		// 0x3600 - 0x36ff | Folder and address book container properties (defined by MAPI)

		// 0x3900 – 0x39ff | Address book properties (defined by MAPI)
		0x39FE: "Seven Bit Email",
		0x39FF: "Seven Bit Display Name",

		// 0x3a00 – 0x3bff | Messaging user properties (defined by MAPI)
		0x3A00: "Account",
		0x3A02: "Callback Phone number",
		0x3A05: "Generation",
		0x3A06: "Given name",
		0x3A08: "Business phone",
		0x3A09: "Home phone",
		0x3A0A: "Initials",
		0x3A0B: "Keyword",
		0x3A0C: "Language",
		0x3A0D: "Location",
		0x3A11: "Surname",
		0x3A15: "Postal address",
		0x3A16: "Company name",
		0x3A17: "Title",
		0x3A18: "Department",
		0x3A19: "Office location",
		0x3A1A: "Primary phone",
		0x3A1B: "Business phone2",
		0x3A1C: "Mobile phone",
		0x3A1D: "Radio phone number",
		0x3A1E: "Car phone number",
		0x3A1F: "Other phone",
		0x3A20: "Transmit display name",
		0x3A21: "Pager",
		0x3A22: "User certificate",
		0x3A23: "PrimaryFax",
		0x3A24: "BusinessFax",
		0x3A25: "Home Fax",
		0x3A26: "Country",
		0x3A27: "Locality",
		0x3A28: "State Or Province",
		0x3A29: "Street address",
		0x3A2A: "PostalCode",
		0x3A2B: "PostOfficeBox",
		0x3A2C: "Telex",
		0x3A2D: "ISDN",
		0x3A2E: "Assistant phone",
		0x3A2F: "Home phone 2",
		0x3A40: "Sender Rich Info", // https://learn.microsoft.com/en-us/office/client-developer/outlook/mapi/pidtagsendrichinfo-canonical-property
		0x3A44: "Middle name",
		0x3A45: "Display name prefix",
		0x3A46: "Profession",
		0x3A48: "Spouse name",
		0x3A4B: "TTY TTD radio phone",
		0x3A4C: "FTP site",
		0x3A4E: "Manager name",
		0x3A4F: "Nickname",
		0x3A51: "Business homepage",
		0x3A57: "Company main phone",
		0x3A58: "Children's names",
		0x3A59: "Home City",
		0x3A5A: "Home Country",
		0x3A5B: "Home Postal Code",
		0x3A5C: "Home State Or Province",
		0x3A5D: "Home Street",
		0x3A5F: "Other address city",
		0x3A60: "Other address country",
		0x3A61: "Other address post code",
		0x3A62: "Other address province",
		0x3A63: "Other address street",
		0x3A64: "Other address PO Box",

		// 0x3c00 – 0x3cff | Distribution list properties (defined by MAPI)
		// 0x3d00 – 0x3dff | Profile properties (defined by MAPI)

		// 0x3e00 – 0x3fff | Status object properties (defined by MAPI)
		0x3FF7: "Server",
		0x3FF8: "Creator1",
		0x3FFA: "Creator2",
		0x3FFC: "To Email",

		// 0x4000 - 0x57ff | Message envelope properties (defined by transport providers)
		0x4022: "Creator Address Type",
		0x4023: "Creator Email Address",
		0x4024: "Last Modifier Address Type",
		0x4025: "Last Modifier Address",
		0x4030: "Sender Simple Display Name",
		0x4031: "Sent Representing Simple DisplayName",
		0x4034: "Received By Simple Display Name",
		0x4035: "Received By Representing Simple Display Name",
		0x4038: "Creator Simple Display Name",
		0x4039: "Last Modifier Simple Display Name",
		0x403D: "To address type",
		0x403E: "To Email2",

		// 0x5800 – 0x5fff | Recipient properties (defined by transport and address book providers)
		0x5d01: "Sender SMTP Address",
		0x5d02: "Sent Representing SMTP email",
		0x5d07: "Received By SMTP Address",
		0x5d08: "Received By Representing SMTP Address",
		0x5d0a: "Creator SMTP Address",
		0x5d0b: "Last Modifier SMTP Address",
		0x5FF6: "To",

		// 0x6000 - 0x65ff | Non-transmittable message properties (defined by clients)
		// 0x6600 – 0x67ff | Non-transmittable properties (defined by a service provider). These properties can be visible or invisible to users.
		// 0x67f0 – 0x67ff | Secure profile properties. These properties can be hidden and encrypted.
		// 0x6800 – 0x7bff | Message content properties for custom message classes (defined by creators of those classes)
		//  0x7c00 – 0x7fff | Non-transmittable properties for custom message classes (defined by creators of those classes)

		//  0x8000 – 0xfffe | Named properties (defined by clients and occasionally service providers). These properties are identified by name through the IMAPIProp::GetNamesFromIDs and IMAPIProp::GetIDsFromNames methods.
		
		//0x800a: "Authentication Results",
		//0x8010: "Creation Date Time",
		0x8015: "Microsoft Information Protection (MSIP) Label",
		//0x8017: "Address Entry Display Table",
		//0x8034: "Creation Date Time 2",

		//  0xffff | Special error value PROP_ID_INVALID (reserved by MAPI)

	}

	name, known := allProps[intID]

	if known {
		return name
	} else {
		return PropertyUnknown
	}
}

// Add a property as known or unknown
func (msg *Message) addPropertyToMessage(msgProps EntryProperty, verbose bool) error {

	propertyTypeInt, err := strconv.ParseInt(msgProps.PropertyType, 16, 32)

	if err != nil {
		return fmt.Errorf("error parsing class %s into an int: %s", msgProps.PropertyType, err.Error())
	}

	propertyName := GetPropertyName(propertyTypeInt)

	var dataString string

	// for all text-encoded properties
	if msgProps.Encoding == AsciiEncoding || msgProps.Encoding == UnicodeEncoding {
		dataString = msgProps.Data.(string)
	} else {
		// cover BinaryEncoding and other unknowns
		bytes, err := getInfAsBytes(msgProps.Data)

		if err != nil {
			return fmt.Errorf("failed to get data bytes")
		}

		// convert it to base64 string
		dataString = base64.StdEncoding.EncodeToString(bytes)
	}

	// skip empty fields
	if len(dataString) == 0 {
		return nil
	}

	// if we recognise the property, store it
	if propertyName != PropertyUnknown {
		msg.Properties[propertyName] = dataString
	} else {
		// if we don't recognise it, store it in the other map using the b64 encoded data value
		prop := UnknownProperty{PropertyType: msgProps.PropertyType, Encoding: msgProps.Encoding, Data: dataString}

		// add it to the map
		msg.UnknownProperties[propertyTypeInt] = prop

		// print its details if we're being verbose
		if verbose {
			if msgProps.Encoding != AsciiEncoding && msgProps.Encoding != UnicodeEncoding && msgProps.Encoding != BinaryEncoding {
				log.Printf("Field 0x%s uses unknown encoding type: 0x%s\n", msgProps.PropertyType, msgProps.Encoding)
			}

			// check if it has a time
			decoded, err := GetTimeFromString(dataString, propertyTypeInt)

			if err == nil {
				log.Printf("Found time in unknown property 0x%x: %s", propertyTypeInt, decoded.Format(time.RFC3339Nano))
			}
		}
	}

	return nil
}

func getInfAsBytes(key interface{}) ([]byte, error) {
	buf, ok := key.([]byte)

	if !ok {
		return nil, fmt.Errorf("error decoding bytes from interface")
	}

	return buf, nil
}
