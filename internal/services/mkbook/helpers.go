package mkbook

import "encoding/base64"

func writeBinaryFile(
	filename string, byteData []byte,
	callback func(source string, internalFilename string) (string, error),
) (location string, err error) {
	encoded := base64.StdEncoding.EncodeToString(byteData)
	location, err = callback(
		"data:text/plain;charset=utf-8;base64,"+encoded,
		filename,
	)
	return location, err
}
