package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

// SessionManager session管理对象
type SessionManager struct {
	cookieName  string
	maxLifeTime int64
	lock        sync.Mutex
	sessionMap  map[string]Session
}

// Session session原子
type Session struct {
	sessionID string
	lastTime  time.Time
	values    map[interface{}]interface{}
}

// NewSessionManager 新建session管理并将其返回，同时预设定时清理过期session
func NewSessionManager(cookieName string, maxLifeTime int64) *SessionManager {
	manager := &SessionManager{
		cookieName:  cookieName,
		maxLifeTime: maxLifeTime,
		sessionMap:  make(map[string]Session),
	}
	go manager.SessionGC()
	return manager
}

// NewSessionID 返回一个随机构建的sessionID
func (manager *SessionManager) NewSessionID() string {
	bytes := make([]byte, 32)

	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		nano := time.Now().UnixNano()
		return strconv.FormatInt(nano, 10)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

// NewSession 新建一个session并入库，将其设置到response的cookie，返回其sessionID
func (manager *SessionManager) NewSession(w http.ResponseWriter, r *http.Request, sessionValues map[interface{}]interface{}) string {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	newSessionID := url.QueryEscape(manager.NewSessionID())
	session := Session{
		sessionID: newSessionID,
		lastTime:  time.Now(),
		values:    sessionValues,
	}
	manager.sessionMap[newSessionID] = session

	cookie := http.Cookie{
		Name:     manager.cookieName,
		Value:    newSessionID,
		Path:     r.URL.Path,
		HttpOnly: true,
		MaxAge:   int(manager.maxLifeTime),
	}
	http.SetCookie(w, &cookie)
	return newSessionID
}

// EndSession 结束request中包含的session并设置到response的cookie中
func (manager *SessionManager) EndSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	}
	manager.lock.Lock()
	defer manager.lock.Unlock()
	delete(manager.sessionMap, cookie.Value)
	newCookie := http.Cookie{
		Name:     manager.cookieName,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now(),
		MaxAge:   -1,
	}
	http.SetCookie(w, &newCookie)
}

// EndSessionByID 结束指定sessionID的session
func (manager *SessionManager) EndSessionByID(sessionID string) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	delete(manager.sessionMap, sessionID)
}

// GetSessionValue 根据sessionID和key查询存放对应的数据
func (manager *SessionManager) GetSessionValue(sessionID string, key interface{}) (interface{}, error) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	if session, ok := manager.sessionMap[sessionID]; ok == true {
		if val, ok := session.values[key]; ok == true {
			return val, nil
		}
	}
	return nil, errors.New("invalid sessionID")
}

// CheckCookieValid 根据request校验cookie对应的session是否存在或有效
func (manager *SessionManager) CheckCookieValid(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie(manager.cookieName)

	if err != nil || cookie == nil {
		return "", err
	}
	manager.lock.Lock()
	defer manager.lock.Unlock()
	sessionID := cookie.Value
	if session, ok := manager.sessionMap[sessionID]; ok == true {
		session.lastTime = time.Now()
		return sessionID, nil
	}
	return "", errors.New("invalid sessionID")
}

// SessionGC 定期清理过期session
func (manager *SessionManager) SessionGC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	for sessionID, session := range manager.sessionMap {
		if session.lastTime.Unix()+manager.maxLifeTime <= time.Now().Unix() {
			delete(manager.sessionMap, sessionID)
		}
	}

	time.AfterFunc(time.Duration(manager.maxLifeTime)*time.Second, manager.SessionGC)
}
