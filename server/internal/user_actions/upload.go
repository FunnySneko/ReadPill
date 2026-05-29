package user_actions

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func (ah *ActionsHandler) UploadImage(file multipart.File, header multipart.FileHeader) (string, error) {
	defer file.Close()

	filename := uuid.NewString() + filepath.Ext(header.Filename)

	err := os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		return "", err
	}

	dst, err := os.Create("./uploads/" + filename)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", nil
	}

	imagePath := "/images/" + filename
	return imagePath, nil
}
