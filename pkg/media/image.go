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

func DeleteRestaurantImage(imageURL string, folderName string) {
	if imageURL != "" {
		publicID := utils.ExtractCloudinaryPublicID(imageURL, folderName)
		if publicID != "" {
			err := utils.DeleteImage(publicID)
			if err != nil {
				fmt.Println("Failed to delete restaurant image from Cloudinary:", err)
			}
		}
	}
}
