# builds convert for mac
build-mac:
	GOOS=darwin go build cmd/convert.go

# builds convert for windows
build-windows:
	GOOS=windows go build cmd/convert.go
