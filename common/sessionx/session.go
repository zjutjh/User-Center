package sessionx

import (
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/zjutjh/User-Center/common/errorsx"
)

type Config struct {
	Name     string `json:",optional"`
	Secret   string `json:",optional"`
	Path     string `json:",optional"`
	Domain   string `json:",optional"`
	MaxAge   int    `json:",optional"`
	Secure   bool   `json:",optional"`
	HTTPOnly bool   `json:",optional"`
	SameSite string `json:",optional"`
}

type payload struct {
	UserID    int64 `json:"userId"`
	ExpiredAt int64 `json:"expiredAt"`
}

type Manager struct {
	conf   Config
	codec  *securecookie.SecureCookie
	cookie string
}

func NewManager(conf Config) *Manager {
	name := conf.Name
	if name == "" {
		name = "user-center-session"
	}

	secret := conf.Secret
	if secret == "" {
		secret = "user-center-change-me"
	}

	if conf.Path == "" {
		conf.Path = "/"
	}
	if conf.MaxAge == 0 {
		conf.MaxAge = 7 * 24 * 3600
	}

	codec := securecookie.New([]byte(secret), nil)
	codec.MaxAge(conf.MaxAge)

	return &Manager{
		conf:   conf,
		codec:  codec,
		cookie: name,
	}
}

func (m *Manager) Set(w http.ResponseWriter, userID int64) error {
	encoded, err := m.codec.Encode(m.cookie, payload{
		UserID:    userID,
		ExpiredAt: time.Now().Add(time.Duration(m.conf.MaxAge) * time.Second).Unix(),
	})
	if err != nil {
		return errorsx.ErrUnknown
	}

	http.SetCookie(w, &http.Cookie{
		Name:     m.cookie,
		Value:    encoded,
		Path:     m.conf.Path,
		Domain:   m.conf.Domain,
		MaxAge:   m.conf.MaxAge,
		Secure:   m.conf.Secure,
		HttpOnly: m.conf.HTTPOnly,
		SameSite: parseSameSite(m.conf.SameSite),
	})
	return nil
}

func (m *Manager) Clear(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     m.cookie,
		Value:    "",
		Path:     m.conf.Path,
		Domain:   m.conf.Domain,
		MaxAge:   -1,
		Secure:   m.conf.Secure,
		HttpOnly: m.conf.HTTPOnly,
		SameSite: parseSameSite(m.conf.SameSite),
	})
}

func (m *Manager) UserID(r *http.Request) (int64, error) {
	cookie, err := r.Cookie(m.cookie)
	if err != nil {
		return 0, errorsx.ErrNotLoggedIn
	}

	var data payload
	if err := m.codec.Decode(m.cookie, cookie.Value, &data); err != nil {
		return 0, errorsx.ErrInvalidCookie
	}
	if data.UserID == 0 || data.ExpiredAt <= time.Now().Unix() {
		return 0, errorsx.ErrNotLoggedIn
	}

	return data.UserID, nil
}

func parseSameSite(mode string) http.SameSite {
	switch strings.ToLower(mode) {
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode
	}
}
