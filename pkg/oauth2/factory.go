package oauth2

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
)

// Info contains the information about the authenticated user.
// Email - email of the user.
// AvatarURL - URL of the user's avatar image.
// Name - name of the user.
// Gender - gender of the user.
type Info struct {
	Email     string // Email of the user.
	AvatarURL string // URL of the user's avatar image.
	Name      string // Name of the user.
	Gender    string // Gender of the user.
}

// Factory is an interface that describes a factory for creating OAuth2 authentication instances.
// It has a single method, `CallBack`, which handles the callback from OAuth2 and returns authentication information.
type Factory interface {
	// CallBack handles the callback from OAuth2.
	// It takes in http.ResponseWriter and *http.Request as parameters.
	// Returns *Info and error.
	CallBack(w http.ResponseWriter, r *http.Request) (*Info, error)
}

// retrieveUserInfo retrieves user information from the specified URL using the provided oauth2 configuration,
// http response writer and request. It returns the retrieved user information and any error encountered.
//
// Parameters:
// - cfg: The oauth2 configuration.
// - w: The http response writer.
// - r: The http request.
// - url: The URL to retrieve the user information from.
//
// Returns:
// - []byte: The retrieved user information.
// - error: Any error encountered.
func retrieveUserInfo(cfg *oauth2.Config, w http.ResponseWriter, r *http.Request, url string) ([]byte, error) {
	// Retrieve the oauthstate cookie from the request.
	oauth2state, _ := r.Cookie("oauthstate")

	// Create a context.
	ctx := r.Context()

	// Retrieve the state and code from the request form.
	state := r.FormValue("state")
	code := r.FormValue("code")

	// Check if the retrieved state matches the oauthstate cookie value.
	// If not, redirect the user and return nil values.
	if state != oauth2state.Value {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return nil, nil
	}

	// Exchange the code for an oauth2 token.
	token, err := cfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("unable to exchange token: %w", err)
	}

	// Create an oauth2 client using the token.
	client := cfg.Client(ctx, token)

	// Retrieve the user information from the specified URL.
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("unable to get google url: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body.
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}

	// Return the retrieved user information and any error encountered.
	return b, err
}
