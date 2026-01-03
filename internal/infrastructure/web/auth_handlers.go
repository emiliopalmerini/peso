package web

import (
	"encoding/json"
	"errors"
	"html/template"
	"log/slog"
	"net/http"

	assets "peso"
	"peso/internal/application"
	"peso/internal/domain/user"
	"peso/internal/infrastructure/middleware"
)

type AuthHandlers struct {
	authService *application.AuthService
	templates   *template.Template
	logger      *slog.Logger
}

func NewAuthHandlers(authService *application.AuthService, logger *slog.Logger) *AuthHandlers {
	return &AuthHandlers{
		authService: authService,
		templates:   loadAuthTemplates(),
		logger:      logger,
	}
}

func (h *AuthHandlers) LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	if middleware.UserFromContext(r.Context()) != nil {
		u := middleware.UserFromContext(r.Context())
		http.Redirect(w, r, "/users/"+u.ID().String(), http.StatusSeeOther)
		return
	}

	data := struct {
		Title string
		Error string
	}{
		Title: "Login - Peso",
	}

	if err := h.templates.ExecuteTemplate(w, "login.html", data); err != nil {
		writeError(h.logger, w, r, http.StatusInternalServerError, "Template error", err)
	}
}

func (h *AuthHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(h.logger, w, r, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	u, sess, err := h.authService.Login(email, password)
	if err != nil {
		if err == application.ErrNoPassword {
			http.SetCookie(w, &http.Cookie{
				Name:     "pending_email",
				Value:    email,
				Path:     "/",
				MaxAge:   300, // 5 minutes
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			})
			http.Redirect(w, r, "/set-password", http.StatusSeeOther)
			return
		}

		data := struct {
			Title string
			Error string
		}{
			Title: "Login - Peso",
			Error: "Email o password non validi",
		}
		w.WriteHeader(http.StatusUnauthorized)
		h.templates.ExecuteTemplate(w, "login.html", data)
		return
	}

	middleware.SetSessionCookie(w, sess.Token())
	http.Redirect(w, r, "/users/"+u.ID().String(), http.StatusSeeOther)
}

func (h *AuthHandlers) RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	if middleware.UserFromContext(r.Context()) != nil {
		u := middleware.UserFromContext(r.Context())
		http.Redirect(w, r, "/users/"+u.ID().String(), http.StatusSeeOther)
		return
	}

	data := struct {
		Title string
		Error string
	}{
		Title: "Registrati - Peso",
	}

	if err := h.templates.ExecuteTemplate(w, "register.html", data); err != nil {
		writeError(h.logger, w, r, http.StatusInternalServerError, "Template error", err)
	}
}

func (h *AuthHandlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(h.logger, w, r, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")

	if password != confirmPassword {
		data := struct {
			Title string
			Error string
			Name  string
			Email string
		}{
			Title: "Registrati - Peso",
			Error: "Le password non coincidono",
			Name:  name,
			Email: email,
		}
		w.WriteHeader(http.StatusBadRequest)
		h.templates.ExecuteTemplate(w, "register.html", data)
		return
	}

	u, sess, err := h.authService.Register(name, email, password)
	if err != nil {
		errMsg := "Errore durante la registrazione"
		switch {
		case err == application.ErrEmailAlreadyExists:
			errMsg = "Email già registrata"
		case err == application.ErrInvalidEmail:
			errMsg = "Email non valida"
		case errors.Is(err, user.ErrPasswordTooShort):
			errMsg = "La password deve essere di almeno 8 caratteri"
		}

		data := struct {
			Title string
			Error string
			Name  string
			Email string
		}{
			Title: "Registrati - Peso",
			Error: errMsg,
			Name:  name,
			Email: email,
		}
		w.WriteHeader(http.StatusBadRequest)
		h.templates.ExecuteTemplate(w, "register.html", data)
		return
	}

	middleware.SetSessionCookie(w, sess.Token())
	http.Redirect(w, r, "/users/"+u.ID().String(), http.StatusSeeOther)
}

func (h *AuthHandlers) SetPasswordPageHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("pending_email")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := struct {
		Title string
		Email string
		Error string
	}{
		Title: "Imposta Password - Peso",
		Email: cookie.Value,
	}

	if err := h.templates.ExecuteTemplate(w, "set_password.html", data); err != nil {
		writeError(h.logger, w, r, http.StatusInternalServerError, "Template error", err)
	}
}

func (h *AuthHandlers) SetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(h.logger, w, r, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	cookie, err := r.Cookie("pending_email")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	email := cookie.Value

	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")

	if password != confirmPassword {
		data := struct {
			Title string
			Email string
			Error string
		}{
			Title: "Imposta Password - Peso",
			Email: email,
			Error: "Le password non coincidono",
		}
		w.WriteHeader(http.StatusBadRequest)
		h.templates.ExecuteTemplate(w, "set_password.html", data)
		return
	}

	u, sess, err := h.authService.SetPassword(email, password)
	if err != nil {
		errMsg := "Errore durante l'impostazione della password"
		if errors.Is(err, user.ErrPasswordTooShort) {
			errMsg = "La password deve essere di almeno 8 caratteri"
		}
		data := struct {
			Title string
			Email string
			Error string
		}{
			Title: "Imposta Password - Peso",
			Email: email,
			Error: errMsg,
		}
		w.WriteHeader(http.StatusBadRequest)
		h.templates.ExecuteTemplate(w, "set_password.html", data)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "pending_email",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	middleware.SetSessionCookie(w, sess.Token())
	http.Redirect(w, r, "/users/"+u.ID().String(), http.StatusSeeOther)
}

func (h *AuthHandlers) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	token := middleware.SessionTokenFromContext(r.Context())
	if token != "" {
		h.authService.Logout(token)
	}

	middleware.ClearSessionCookie(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func loadAuthTemplates() *template.Template {
	tmpl := template.New("").Funcs(template.FuncMap{
		"title": func(s string) string {
			if len(s) == 0 {
				return s
			}
			return string(s[0]-32) + s[1:]
		},
		"toJson": func(v interface{}) template.JS {
			b, err := json.Marshal(v)
			if err != nil {
				return template.JS("null")
			}
			return template.JS(string(b))
		},
	})
	template.Must(tmpl.ParseFS(assets.FS, "templates/*.html"))
	return tmpl
}
