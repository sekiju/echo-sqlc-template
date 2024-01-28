package storage

import (
	"bytes"
	"context"
	"echo-sqlc-template/internal/config"
	"echo-sqlc-template/internal/database"
	"fmt"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"net/http"
	"path"
)

var client = &http.Client{}

func CreateFile(p string, f io.Reader) error {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormFile("file", path.Base(p))
	if err != nil {
		return err
	}

	_, err = io.Copy(fw, f)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", config.Data.Storage.Daemon+p, &b)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("bad status: %s", res.Status)
	}

	return nil
}

func MoveFile(hash, folder string) (string, error) {
	img, err := database.Q.GetUploadedImage(context.Background(), database.GetUploadedImageParams{Hash: hash})
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("/%s/%s.%s", folder, uuid.New().String(), img.Extension)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s?destination=%s", config.Data.Storage.Daemon+img.Key, key), nil)
	if err != nil {
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", res.Status)
	}

	err = database.Q.DeleteUploadedImageById(context.Background(), img.ID)
	if err != nil {
		return "", err
	}

	return key, nil
}
