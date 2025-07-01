package utils

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func getCloudinaryURL() string {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")
	return fmt.Sprintf("cloudinary://%s:%s@%s", apiKey, apiSecret, cloudName)
}

func UploadImage(file multipart.File, fileHeader *multipart.FileHeader, folder string) (string, string, error) {
	cld, err := cloudinary.NewFromURL(getCloudinaryURL())
	if err != nil {
		return "", "", err
	}

	ctx := context.Background()
	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: fileHeader.Filename,
		Folder:   folder,
	})
	if err != nil {
		return "", "", err
	}

	return uploadResult.SecureURL, uploadResult.PublicID, nil
}

func ExtractCloudinaryPublicID(url, folderName string) string {
	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return ""
	}

	publicIDWithExt := parts[len(parts)-1]
	ext := path.Ext(publicIDWithExt)

	publicID := strings.TrimSuffix(publicIDWithExt, ext)

	return folderName + "/" + publicID
}

func DeleteImage(publicID string) error {
	cld, err := cloudinary.NewFromURL(getCloudinaryURL())
	if err != nil {
		return err
	}

	ctx := context.Background()
	_, err = cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}
