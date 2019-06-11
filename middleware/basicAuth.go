package middleware

import (
	"net/http"
)

type UserInfo struct {
	Name string
	Pwd  string
}

func (u UserInfo) CheckUserIdentity(h http.HandlerFunc) http.HandlerFunc {
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
