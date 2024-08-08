package qimage

import (
	"errors"
	"mime/multipart"
	"net/http"
	"slices"
	"strconv"
)

func ReadFromMultipart(fileHeader *multipart.FileHeader, i *Imager, allowedFileTypes []string) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}

	defer func(file multipart.File) {
		err = file.Close()
		if err != nil {
			return
		}
	}(file)

	raw := make([]byte, fileHeader.Size)
	_, err = file.Read(raw)
	if err != nil {
		return err
	}

	sortIndex, err := strconv.Atoi(fileHeader.Filename)
	if err != nil {
		return err
	}

	filetype := http.DetectContentType(raw)

	if len(allowedFileTypes) > 0 && !slices.Contains(allowedFileTypes, filetype) {
		return errors.New("the provided file format is not allowed")
	}

	(*i).SetRaw(raw)
	(*i).SetSortIndex(sortIndex)
	(*i).SetFileType(filetype)

	return nil
}
