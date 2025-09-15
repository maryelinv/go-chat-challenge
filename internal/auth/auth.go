package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func Hash(pw string) ([]byte, error) { return bcrypt.GenerateFromPassword([]byte(pw), 12) }

func Check(hash []byte, pw string) error { return bcrypt.CompareHashAndPassword(hash, []byte(pw)) }

func newSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

var sessions = map[string]uint{}

func SetSession(w http.ResponseWriter, userID uint) error {
	sid, err := newSessionID()
	if err != nil {
		return err
	}
	sessions[sid] = userID
	http.SetCookie(w, &http.Cookie{
		Name:     "sid",
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})
	return nil
}

func ClearSession(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie("sid"); err == nil {
		delete(sessions, c.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "sid",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	})
}

func CurrentUserID(r *http.Request) (uint, bool) {
	c, err := r.Cookie("sid")
	if err != nil {
		return 0, false
	}
	uid, ok := sessions[c.Value]
	return uid, ok
}
