package media

import (
	"fmt"

	"food-delivery-app-server/internal/auth"
	"food-delivery-app-server/pkg/utils"
)

func DeleteProfilePicIfNotDefault(profilePic string, folderName string) {
	if profilePic != auth.DefaultProfilePic && profilePic != "" {
		publicID := utils.ExtractCloudinaryPublicID(profilePic, folderName)
		if publicID != "" {
			err := utils.DeleteImage(publicID)
			fmt.Println(err)
		}
	}
}

func DeleteImage(resource string, imageURL string, folderName string) {
	if imageURL != "" {
		publicID := utils.ExtractCloudinaryPublicID(imageURL, folderName)
		if publicID != "" {
			err := utils.DeleteImage(publicID)
			if err != nil {
				fmt.Printf("Failed to %s delete image from Cloudinary: %v\n", resource, err)
			}
		}
	}
}
