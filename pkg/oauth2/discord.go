package oauth2

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

// DiscordFactory represents a factory struct for Discord OAuth2 authentication.
// It contains a pointer to an oauth2.Config struct.
type DiscordFactory struct {
	cfg *oauth2.Config // The OAuth2 configuration. This holds the client ID, client secret, redirect URL, and other
	// parameters needed for the OAuth2 authentication flow.
}

// NewDiscordFactory creates a new DiscordFactory instance.
func NewDiscordFactory(cfg *oauth2.Config) Factory {
	return &DiscordFactory{
		cfg: cfg,
	}
}

// CallBack is a function that retrieves user information from Discord API.
//
// Parameters:
//
//	w: http.ResponseWriter for writing response back.
//	r: *http.Request for incoming request data.
//
// Return:
//
//	*Info: Struct containing user information.
//	error: Any error that occurred during the process.
func (f *DiscordFactory) CallBack(w http.ResponseWriter, r *http.Request) (*Info, error) {
	userInfoURL := "https://discord.com/api/users/@me"
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
