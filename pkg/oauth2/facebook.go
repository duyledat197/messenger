package oauth2

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

// FacebookFactory represents a factory struct for Facebook OAuth2 authentication.
// It contains a pointer to an oauth2.Config struct.
type FacebookFactory struct {
	cfg *oauth2.Config // The OAuth2 configuration.
}

// NewFacebookFactory creates a new FacebookFactory instance.
func NewFacebookFactory(cfg *oauth2.Config) Factory {
	return &FacebookFactory{
		cfg: cfg,
	}
}

// CallBack handles the callback from Google OAuth.
// It takes in http.ResponseWriter and *http.Request as parameters.
// Returns *Info and error.
func (f *FacebookFactory) CallBack(w http.ResponseWriter, r *http.Request) (*Info, error) {
	userInfoURL := "https://graph.facebook.com/me?fields=id,name,first_name,last_name,picture,email"
	b, err := retrieveUserInfo(f.cfg, w, r, userInfoURL)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve social user info: %w", err)
	}
	type info struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		Birthday  string `json:"birthday"`
		Gender    string `json:"gender"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Picture   struct {
			Height       int    `json:"height"`
			IsSilhouette bool   `json:"is_silhouette"`
			URL          string `json:"url"`
			Width        int    `json:"width"`
		} `json:"picture,omitempty"`
	}
	data := &info{}
	if err := json.Unmarshal(b, data); err != nil {
		return nil, fmt.Errorf("unable to unmarshal body: %w", err)
	}

	return &Info{
		Email:     data.Email,
		AvatarURL: data.Picture.URL,
		Name:      data.Name,
		Gender:    data.Gender,
	}, nil
}
