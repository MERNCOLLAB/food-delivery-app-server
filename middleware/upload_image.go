package middleware

import (
	"bytes"
	appErr "food-delivery-app-server/pkg/errors"
	"io"

	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	MaxImageSize = 2 << 20
)

var allowedTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
}

func UploadImageValidator(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, header, err := c.Request.FormFile(key)
		if err != nil {
			appError := appErr.NewBadRequest("Image file is required", err)
			c.JSON(appError.Code, gin.H{"error": appError.Message})
			return
		}

		defer file.Close()

		if header.Size > MaxImageSize {
			appError := appErr.NewBadRequest("File size exceeds 2 MB", nil)
			c.JSON(appError.Code, gin.H{"error": appError.Message})
			return
		}

		buff := bytes.NewBuffer(nil)
		if _, err := io.Copy(buff, file); err != nil {
			appError := appErr.NewBadRequest("Failed to read the file", err)
			c.JSON(appError.Code, gin.H{"error": appError.Message})
			return
		}

		filetype := http.DetectContentType(buff.Bytes()[:512])
		if !allowedTypes[filetype] {
			appError := appErr.NewBadRequest("Invalid file type. Only JPEG, PNG, and WEBP are allowed", err)
			c.JSON(appError.Code, gin.H{"error": appError.Message})
			return
		}

		c.Set("imageFile", buff)
		c.Set("imageHeader", header)
		c.Next()
	}
}
