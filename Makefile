default: linux

linux:
	go build -o file-inspector-elf ui/*.go

windows:
	env GOOS=windows GOARCH=amd64 go build -o file-inspector.exe ui/*.go

win-cross:
	fyne-cross windows