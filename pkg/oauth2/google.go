package oauth2

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

// GoogleFactory is a factory struct for Google OAuth2 authentication.
// It contains a pointer to a Client.
type GoogleFactory struct {
	cfg *oauth2.Config
}

// NewGoogleFactory creates a new GoogleFactory instance.
func NewGoogleFactory(cfg *oauth2.Config) Factory {
	return &GoogleFactory{
		cfg: cfg,
	}
}

// CallBack handles the callback from Google OAuth.
// It takes in http.ResponseWriter and *http.Request as parameters.
// Returns *Info and error.
func (f *GoogleFactory) CallBack(w http.ResponseWriter, r *http.Request) (*Info, error) {
	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo"
	b, err := retrieveUserInfo(f.cfg, w, r, userInfoURL)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve social user info: %w", err)
	}
	type info struct {
		Email      string `json:"email"`
		Name       string `json:"name"`
		GivenName  string `json:"given_name"`
		FamilyName string `json:"family_name"`
		Picture    string `json:"picture"`
		Gender     string `json:"gender"`
		Locale     string `json:"locale"`
	}
	data := &info{}
	if err := json.Unmarshal(b, data); err != nil {
		return nil, fmt.Errorf("unable to unmarshal body: %w", err)
	}

	return &Info{
		Email:     data.Email,
		AvatarURL: data.Picture,
		Name:      data.Name,
		Gender:    data.Gender,
	}, nil
}
