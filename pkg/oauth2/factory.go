package oauth2

import (
	"net/http"
)

type Factory interface {
	Redirect(w http.ResponseWriter, r *http.Request)
	CallBack(w http.ResponseWriter, r *http.Request)
}
