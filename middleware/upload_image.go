package middleware

import (
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

func UploadImageValidator(key string, optional ...bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		isOptional := false
		if len(optional) > 0 {
			isOptional = optional[0]
		}
		file, header, err := c.Request.FormFile(key)
		if err != nil {
			if isOptional {
				c.Next()
				return
			}
			appError := appErr.NewBadRequest("Image file is required", err)
			c.JSON(appError.Code, gin.H{"error": appError.Message})
			c.Abort()
			return
		}

		if header.Size > MaxImageSize {
			appError := appErr.NewBadRequest("File size exceeds 2 MB", nil)
			c.JSON(appError.Code, gin.H{"error": appError.Message})
			c.Abort()
			return
		}

		buff := make([]byte, 512)
		n, err := file.Read(buff)
		if err != nil && err != io.EOF {
			appError := appErr.NewBadRequest("Failed to read the file", err)
			c.JSON(appError.Code, gin.H{"error": appError.Message})
			c.Abort()
			return
		}

		filetype := http.DetectContentType(buff[:n])
		if !allowedTypes[filetype] {
			appError := appErr.NewBadRequest("Invalid file type. Only JPEG, PNG, and WEBP are allowed", err)
			c.JSON(appError.Code, gin.H{"error": appError.Message})
			c.Abort()
			return
		}

		if seeker, ok := file.(io.Seeker); ok {
			seeker.Seek(0, io.SeekStart)
		} else {
			appError := appErr.NewInternal("Cannot seek the file for upload", nil)
			c.JSON(appError.Code, gin.H{"error": appError.Message})
			file.Close()
			return
		}

		c.Set("imageFile", file)
		c.Set("imageHeader", header)
		c.Next()
	}
}
