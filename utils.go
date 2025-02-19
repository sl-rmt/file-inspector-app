package main

// check the file extension is one we support
// return false if not
func fileOkayToProcess(fileExtension string) bool {

	switch fileExtension {
	case ".docx":
		fallthrough
	case ".eml":
		fallthrough
	case ".msg":
		fallthrough
	case ".pdf":
		return true
	default:
		return false
	}
}
