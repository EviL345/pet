package utils

import (
	"blog/internal/config"
	"encoding/json"
	"net/http"
)

// ReadRequest reads a JSON request and validates it against a given struct
func ReadRequest(r *http.Request, i interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(i); err != nil {
		return err
	}

	if err := validate.Struct(i); err != nil {
		return err
	}

	return nil
}

func CreateSessionCookie(cfg *config.Config, session string) *http.Cookie {
	return &http.Cookie{
		Name:       cfg.Session.Name,
		Value:      session,
		Path:       "/",
		MaxAge:     cfg.Cookie.MaxAge,
		Secure:     cfg.Cookie.Secure,
		HttpOnly:   cfg.Cookie.HTTPOnly,
		SameSite:   0,
		RawExpires: "",
	}
}

func DeleteSessionCookie(w http.ResponseWriter, cfg *config.Config) {
	cookie := &http.Cookie{
		Name:   cfg.Session.Name,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)
}
