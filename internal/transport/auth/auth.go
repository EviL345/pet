package auth

import (
	"blog/internal/config"
	"blog/internal/models"
	"blog/internal/repository"
	"blog/internal/service"
	"blog/pkg/httperrors"
	"blog/pkg/utils"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"log"

	"net/http"
)

type Server struct {
	authService    service.Auth
	sessionService service.Session
	cfg            *config.Config
}

func NewHandlers(authService service.Auth, sessionService repository.RedisSession, cfg config.Config) *Server {
	return &Server{
		authService:    authService,
		sessionService: sessionService,
		cfg:            &cfg,
	}
}

func (h *Server) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := &models.User{}
		if err := utils.ReadRequest(r, user); err != nil {
			log.Println(err)

			httperrors.ErrorResponse(w, r, err)

			return
		}

		createdUser, err := h.authService.Register(user)
		if err != nil {
			log.Println(err)

			httperrors.ErrorResponse(w, r, err)

			return
		}
		session, err := h.sessionService.CreateSession(&models.Session{
			UserID: createdUser.User.ID,
		}, h.cfg.Session.Expire)
		if err != nil {
			log.Println(err)

			httperrors.ErrorResponse(w, r, err)

			return
		}

		cookie := utils.CreateSessionCookie(h.cfg, session)
		http.SetCookie(w, cookie)

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, createdUser)
	}

}

func (h *Server) Login() http.HandlerFunc {
	/*type Login struct {
		Email    string `json:"email" db:"email" validate:"omitempty,lte=60,email"`
		Password string `json:"password,omitempty" db:"password" validate:"required, gte=6"`
	}*/
	return func(w http.ResponseWriter, r *http.Request) {
		//		login := &Login{}
		user := &models.User{}
		if err := utils.ReadRequest(r, user); err != nil {
			log.Println(err)

			httperrors.ErrorResponse(w, r, err)

			return
		}

		userWithToken, err := h.authService.Login(user)
		if err != nil {
			log.Println(err)

			httperrors.ErrorResponse(w, r, err)

			return
		}

		session, err := h.sessionService.CreateSession(&models.Session{UserID: userWithToken.User.ID}, h.cfg.Session.Expire)
		if err != nil {
			log.Println(err)

			httperrors.ErrorResponse(w, r, err)

			return
		}

		cookie := utils.CreateSessionCookie(h.cfg, session)
		http.SetCookie(w, cookie)

		render.Status(r, http.StatusOK)
		render.JSON(w, r, userWithToken)
	}

}

func (h *Server) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(h.cfg.Session.Name)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				log.Println(err)

				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, err)

				return
			}

			log.Println(err)

			httperrors.ErrorResponse(w, r, err)

			return
		}

		if err = h.sessionService.DeleteSessionByID(cookie.Value); err != nil {
			log.Println(err)

			httperrors.ErrorResponse(w, r, err)
		}

		utils.DeleteSessionCookie(w, h.cfg)

		render.Status(r, http.StatusOK)
		render.JSON(w, r, nil)

	}

}

func (h *Server) ChangePassword() http.HandlerFunc {
	type Passwords struct {
		OldPassword string `json:"old_password" validate:"required,omitempty,gte=6"`
		NewPassword string `json:"new_password" validate:"required,omitempty,gte=6"`
	}
	return func(w http.ResponseWriter, r *http.Request) {

		uID, err := uuid.Parse(r.URL.Query().Get("user_id"))
		if err != nil {
			log.Println(err)

			httperrors.ErrorResponse(w, r, err)

			return
		}

		passwords := &Passwords{}
		if err = utils.ReadRequest(r, passwords); err != nil {
			log.Println(err)

			httperrors.ErrorResponse(w, r, err)

			return
		}

		user, err := h.authService.GetUserById(uID.String())
		if err != nil {
			log.Println(err)

			httperrors.ErrorResponse(w, r, err)

			return
		}

		updatedUser, err := h.authService.ChangePassword(user, passwords.OldPassword, passwords.NewPassword)
		if err != nil {
			log.Println(err)

			httperrors.ErrorResponse(w, r, err)

			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, updatedUser)
	}

}

func (h *Server) GetUserById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uID, err := uuid.Parse(r.URL.Query().Get("user_id"))
		if err != nil {
			log.Println(err)

			httperrors.ErrorResponse(w, r, err)

			return
		}

		user, err := h.authService.GetUserById(uID.String())
		if err != nil {
			log.Println(err)

			httperrors.ErrorResponse(w, r, err)

			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, user)
	}

}
