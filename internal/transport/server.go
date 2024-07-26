package transport

import "net/http"

type Auth interface {
	Login() http.HandlerFunc
	Logout() http.HandlerFunc
	Register() http.HandlerFunc
	Update() http.HandlerFunc
	GetUserById() http.HandlerFunc
}
