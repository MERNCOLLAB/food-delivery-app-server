package oauth

import (
	"encoding/json"
	"net/http"

	appErr "food-delivery-app-server/pkg/errors"
)

type GoogleResponseData struct {
	Email      string `json:"email"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Picture    string `json:"picture"`
}

type FacebookResponseData struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Picture   struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
}

type UserInfo struct {
	Email          string
	FirstName      string
	LastName       string
	ProfilePicture string
	Provider       string
}

func VerifyGoogleToken(accessToken string) (*UserInfo, error) {
	var data GoogleResponseData
	resp, err := http.Get("https://www.googleapis.com/oauth2/v3/userinfo?access_token=" + accessToken)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, appErr.NewInternal("failed to verify Google token", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &UserInfo{
		Email:          data.Email,
		FirstName:      data.GivenName,
		LastName:       data.FamilyName,
		ProfilePicture: data.Picture,
		Provider:       "google",
	}, nil
}

func VerifyFacebookToken(accessToken string) (*UserInfo, error) {
	var data FacebookResponseData
	url := "https://graph.facebook.com/me?fields=id,first_name,last_name,email,picture&access_token=" + accessToken
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, appErr.NewInternal("failed to verify Facebook token", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &UserInfo{
		Email:          data.Email,
		FirstName:      data.FirstName,
		LastName:       data.LastName,
		ProfilePicture: data.Picture.Data.URL,
		Provider:       "facebook",
	}, nil
}
