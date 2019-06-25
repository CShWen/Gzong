package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewSessionManager(t *testing.T) {
	sm := NewSessionManager("gzCookie", 30)
	if sm == nil {
		t.Error("创建sessionManager不应为空")
	}
}

func TestSessionManager_NewSessionId(t *testing.T) {
	sm := NewSessionManager("gzCookie", 30)
	sessionId := sm.NewSessionId()
	if len(sessionId) <= 0 {
		t.Error("newSessionId不应为空")
	}
}

func TestSessionManager_NewSession_and_CheckCookieValid(t *testing.T) {
	sm := NewSessionManager("gzCookie", 30)
	r, _ := http.NewRequest("GET", "http://127.0.0.1:9876/", nil)
	w := httptest.NewRecorder()
	tw := httptest.NewRecorder()
	value := make(map[interface{}]interface{})
	sm.NewSession(w, r, value)
	sm.NewSession(w, r, value)
	sessionId := sm.NewSession(tw, r, value)

	if len(sm.sessionMap) != 3 {
		t.Error("sessionManager包含的session数目不符合预期")
	}

	sign := false
	cookieStrArray := strings.Split(tw.Header().Get("Set-Cookie"), "; ")
	for i, cookieStr := range cookieStrArray {
		if strings.Contains(cookieStrArray[i], "=") == true {
			cookie := strings.Split(cookieStr, "=")
			key, value := cookie[0], cookie[1]
			if key == "gzCookie" && value == sessionId {
				sign = true
			}
		}
	}
	if sign == false {
		t.Log("response的cookie中不包含预设的session信息")
	}

	checkSessionId, err := sm.CheckCookieValid(w, r)
	if err == nil || sessionId == checkSessionId {
		t.Error("无session理应认证失败")
	}

	r.Header.Add("Cookie", "gzCookie=test")
	if err == nil || sessionId == checkSessionId {
		t.Error("错误的session理应认证失败")
	}

	r.Header.Del("Cookie")
	r.Header.Add("Cookie", "gzCookie="+sessionId)
	checkSessionId, err = sm.CheckCookieValid(w, r)
	if err != nil || sessionId != checkSessionId {
		t.Error("session未通过校验")
	}

}

func TestSessionManager_EndSession(t *testing.T) {
	sm := NewSessionManager("gzCookie", 30)
	r, _ := http.NewRequest("GET", "http://127.0.0.1:9876/", nil)
	w := httptest.NewRecorder()
	value := make(map[interface{}]interface{})
	sm.NewSession(w, r, value)
	sm.NewSession(w, r, value)
	sessionId := sm.NewSession(w, r, value)

	r.Header.Add("Cookie", "gzCookie="+sessionId)

	checkSessionId, err := sm.CheckCookieValid(w, r)
	if err != nil || sessionId != checkSessionId {
		t.Error("session未通过校验")
	}

	sm.EndSession(w, r)

	if len(sm.sessionMap) != 2 {
		t.Error("sessionManager包含的session数目不符合预期")
	}

	checkSessionId, err = sm.CheckCookieValid(w, r)
	if err == nil || sessionId == checkSessionId {
		t.Error("session未被清除")
	}
}

func TestSessionManager_EndSessionById(t *testing.T) {
	sm := NewSessionManager("gzCookie", 30)
	r, _ := http.NewRequest("GET", "http://127.0.0.1:9876/", nil)
	w := httptest.NewRecorder()
	value := make(map[interface{}]interface{})
	sessionId := sm.NewSession(w, r, value)

	r.Header.Add("Cookie", "gzCookie="+sessionId)

	checkSessionId, err := sm.CheckCookieValid(w, r)
	if err != nil || sessionId != checkSessionId {
		t.Error("session未通过校验")
	}

	sm.EndSessionById(sessionId)

	checkSessionId, err = sm.CheckCookieValid(w, r)
	if err == nil || sessionId == checkSessionId {
		t.Error("session未被清除")
	}
}

func TestSessionManager_GetSessionValue(t *testing.T) {
	sm := NewSessionManager("gzCookie", 30)
	r, _ := http.NewRequest("GET", "http://127.0.0.1:9876/", nil)
	w := httptest.NewRecorder()

	valueMap := make(map[interface{}]interface{})
	valueMap["a1"] = "b1"
	valueMap["a2"] = "b2"
	sessionId := sm.NewSession(w, r, valueMap)

	value, err := sm.GetSessionValue(sessionId, "a1")
	if value != "b1" || err != nil {
		t.Error("有效的session应取得符合预期的存储内容")
	}
	value, err = sm.GetSessionValue("test", "a1")
	if value != nil || err == nil {
		t.Error("无效session不应可获取到存储内容")
	}
}

func TestSessionManager_SessionGC(t *testing.T) {
	sm := NewSessionManager("gzCookie", 1)
	r, _ := http.NewRequest("GET", "http://127.0.0.1:9876/", nil)
	w := httptest.NewRecorder()
	value := make(map[interface{}]interface{})
	sm.NewSession(w, r, value)
	sm.NewSession(w, r, value)
	sm.NewSession(w, r, value)

	if len(sm.sessionMap) != 3 {
		t.Error("sessionManager中有效的session数目不符合预期")
	}
	time.Sleep(1234 * time.Millisecond)
	if len(sm.sessionMap) != 0 {
		t.Error("sessionManager中有效的session数目不符合预期")
	}
}
