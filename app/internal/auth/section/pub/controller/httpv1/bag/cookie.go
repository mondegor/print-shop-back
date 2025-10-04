package bag

import (
	"net/http"
	"time"
)

type (
	// RefreshTokenCookie - comment struct.
	RefreshTokenCookie struct {
		name   string
		domain string
		path   string
		expiry time.Duration
	}
)

// NewRefreshTokenCookie - создаёт объект RefreshTokenCookie.
func NewRefreshTokenCookie(name, domain, path string, expiry time.Duration) *RefreshTokenCookie {
	return &RefreshTokenCookie{
		name:   name,
		domain: domain,
		path:   path,
		expiry: expiry,
	}
}

// GetValue - comment method.
func (c *RefreshTokenCookie) GetValue(r *http.Request) (refreshToken string) {
	cookie, err := r.Cookie(c.name)
	if err != nil {
		return ""
	}

	return cookie.Value
}

// SetValue - comment method.
func (c *RefreshTokenCookie) SetValue(w http.ResponseWriter, refreshToken string) {
	ck := http.Cookie{
		Name:     c.name,
		Value:    refreshToken,
		Path:     c.path,
		Domain:   c.domain,
		Expires:  time.Now().Add(c.expiry),
		MaxAge:   int(c.expiry.Seconds()),
		Secure:   true, // TODO: for test ?
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode, // TODO: for test ?
		// Partitioned: true,                 // TODO: for test ?
	}

	http.SetCookie(w, &ck)
}
