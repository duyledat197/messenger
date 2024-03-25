package oauth2

import (
	"log/slog"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const oauth2state = "something"

type GoogleFactory struct {
	cfg *oauth2.Config
}

func NewGoogleFactory(clientID, clientSecret, redirectURL string) Factory {
	return &GoogleFactory{
		cfg: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
			},
			Endpoint: google.Endpoint,
		},
	}
}
func (f *GoogleFactory) Redirect(w http.ResponseWriter, r *http.Request) {
	url := f.cfg.AuthCodeURL(oauth2state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (f *GoogleFactory) CallBack(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	state := r.FormValue("state")
	code := r.FormValue("code")

	if state != oauth2state {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	token, err := f.cfg.Exchange(ctx, code)
	if err != nil {
		slog.Error("Exchange failed with " + err.Error() + "\n")
		return
	}
	client := f.cfg.Client(ctx, token)

	client.Get("https://www.googleapis.com/auth/userinfo.email")
}
