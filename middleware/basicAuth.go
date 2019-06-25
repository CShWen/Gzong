package middleware

import (
	"encoding/base64"
	"net/http"
)

type BaseUser struct {
	Name string
	Pwd  string
}

func (u BaseUser) BasicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		baName, baPwd, ok := r.BasicAuth()
		if ok == true && u.Name == baName && u.Pwd == baPwd {
			h.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`Unauthorized`))
		}
	}
}

// 将username和password进行base64加密编码
func Base64Encode(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
