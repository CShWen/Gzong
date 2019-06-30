package complex

import (
	"fmt"
	"github.com/cshwen/gzong"
	"github.com/cshwen/gzong/middleware"
	"io"
	"log"
	"net/http"
)

var globalSessMgr *middleware.SessionManager

func main() {
	gz := gzong.New()
	gz.GET("/testhello", helloFunc)
	gz.GET("/testjson", jsonFunc)
	gz.POST("/testpost", testPostFunc)
	gz.AddMiddleware(middleware.RequestDetailsLog)
	gz.AddMiddleware(middleware.ServiceConSumeTimeLog)
	// request headers记得添加 Authorization: [Basic c3M6cHdk]，否则请求401
	name, pwd := "ss", "pwd"
	fmt.Println("request headers记得添加 Authorization: ", "Basic "+middleware.Base64Encode(name, pwd))
	u := middleware.BaseUser{Name: name, Pwd: pwd}
	gz.AddMiddleware(u.BasicAuth)

	globalSessMgr = middleware.NewSessionManager("gzCookie", 30)
	gz.GET("/login", businessAndSessionFunc)
	gz.GET("/business", businessAndSessionFunc)
	gz.GET("/logout", businessAndSessionFunc)

	gz.Run(":8080")
}

func helloFunc(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "hello gzong.")
}

func jsonFunc(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"alive": true, "cc": "ss"}`))
}

func testPostFunc(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "test post.")
}

func businessAndSessionFunc(w http.ResponseWriter, r *http.Request) {
	smResult := sessionFilter(w, r)
	if !smResult {
		log.Println("session验证结束，无需执行后续事务")
		return
	}
	io.WriteString(w, `session验证通过，继续后续业务`)
}

func sessionFilter(w http.ResponseWriter, r *http.Request) (result bool) {
	result = false
	sessionID, err := globalSessMgr.CheckCookieValid(w, r)
	if err != nil {
		if r.URL.Path == "/login" {
			sessionID := globalSessMgr.NewSession(w, r, make(map[interface{}]interface{}))
			fmt.Println("new SessionID:\t", sessionID)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("login success."))
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
		}
	} else {
		if r.URL.Path == "/logout" {
			globalSessMgr.EndSessionByID(sessionID)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("logout success."))
		} else {
			result = true
		}
	}
	return
}
