package middleware

import (
	"encoding/base64"
	"net/http"
)

// BaseUser 用户元数据
type BaseUser struct {
	Name string
	Pwd  string
}

// BasicAuth 简单的身份认证，校验通过则继续下层业务逻辑，不通过则401
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

// Base64Encode 将username和password进行base64加密编码
func Base64Encode(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
