package utils

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)



func getCloudinaryURL() string {
    cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
    apiKey := os.Getenv("CLOUDINARY_API_KEY")
    apiSecret := os.Getenv("CLOUDINARY_API_SECRET")
    return fmt.Sprintf("cloudinary://%s:%s@%s", apiKey, apiSecret, cloudName)
}

func UploadImage(file multipart.File, fileHeader *multipart.FileHeader) (string, string, error) {
    cld, err := cloudinary.NewFromURL(getCloudinaryURL())
    if err != nil {
        return "", "", err
    }

    ctx := context.Background()
    uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
        PublicID: fileHeader.Filename,
        Folder:   "profile_pictures",
    })
    if err != nil {
        return "", "", err
    }

    return uploadResult.SecureURL, uploadResult.PublicID, nil
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