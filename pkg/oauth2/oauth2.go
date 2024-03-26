// Package oauth2 represents of social oauth2 implementation
package oauth2

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

// Client is a struct that holds the configuration and factory for a given OAuth2 client.
//
// Fields:
//   - factory: Factory - an instance of Factory, implementing the OAuth2 authentication flow.
//   - cfg: *oauth2.Config - the configuration for the OAuth2 client.
type Client struct {
	factory Factory        // An instance of Factory, implementing the OAuth2 authentication flow.
	cfg     *oauth2.Config // The configuration for the OAuth2 client.
}

// NewClient creates a new client for the given social platform.
// It takes a social string, clientID string, clientSecret string, and redirectURL string as parameters.
// It returns a pointer to a Client.
func NewClient(social string, clientID, clientSecret, redirectURL string) *Client {
	var (
		cfg     *oauth2.Config
		factory Factory
	)

	switch social {
	case "google":
		cfg = &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.profile",
				"https://www.googleapis.com/auth/userinfo.email",
			},
			Endpoint: google.Endpoint,
		}
		factory = NewGoogleFactory(cfg)

	case "facebook":
		cfg = &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"public_profile", "email"},
			Endpoint:     facebook.Endpoint,
		}
		factory = NewFacebookFactory(cfg)

	case "discord":
		cfg = &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"identify"},
			Endpoint: oauth2.Endpoint{
				AuthURL:   "https://discord.com/api/oauth2/authorize",
				TokenURL:  "https://discord.com/api/oauth2/token",
				AuthStyle: oauth2.AuthStyleInParams,
			},
		}
		factory = NewDiscordFactory(cfg)

	case "github":
		// endpoint = github.Endpoint
		// scopes = []string{github.ScopeUser}
		// factory = NewGithubFactory(clientID, clientSecret, redirectURL)
	default:

	}
	return &Client{
		factory: factory,
		cfg:     cfg,
	}
}

// Redirect redirects the client to the OAuth2 authentication URL.
// It takes an http.ResponseWriter and an http.Request as parameters and does not return anything.
func (c *Client) Redirect(w http.ResponseWriter, r *http.Request) {
	oauth2state := generateStateOauthCookie(w)
	url := c.cfg.AuthCodeURL(oauth2state)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// generateStateOauthCookie generates a state OAuth cookie for the provided http.ResponseWriter.
// It returns the generated state string.
func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}
