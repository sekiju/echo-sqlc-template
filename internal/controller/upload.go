package controller

import (
	"bytes"
	"echo-sqlc-template/internal/controller/middlewares"
	"echo-sqlc-template/internal/database"
	"echo-sqlc-template/internal/services/storage"
	"echo-sqlc-template/internal/utils"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"path"
)

type UploadController struct{}

func NewUploadController(e *echo.Echo) {
	c := UploadController{}

	group := e.Group("/upload")
	group.Use(middlewares.Authorization())
	group.POST("/image", c.Image)
}

func (UploadController) Image(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}

	defer src.Close()

	buf := &bytes.Buffer{}
	tee := io.TeeReader(src, buf)

	user := c.Get("user").(database.User)

	img, ext, err := image.Decode(tee)
	if err != nil {
		return err
	}

	filePath := fmt.Sprintf("/temporary/%d/%s.%s", user.ID, uuid.New().String(), ext)

	err = storage.CreateFile(filePath, buf)
	if err != nil {
		return err
	}

	uploadedImage, err := database.Q.CreateUploadedImage(ctx, database.CreateUploadedImageParams{
		Hash:      utils.GenerateRandomString(32),
		Key:       filePath,
		Size:      int32(file.Size),
		Extension: ext[1:],
		Height:    int32(img.Bounds().Dy()),
		Width:     int32(img.Bounds().Dx()),
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}

	return c.JSON(201, uploadedImage)
}
