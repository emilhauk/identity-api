package util

import (
	"net/http"
	"time"
)

func DeleteCookie(cookie *http.Cookie, w http.ResponseWriter) {
	cookie.Value = ""
	cookie.Expires = time.Unix(0, 0)
	http.SetCookie(w, cookie)
}
