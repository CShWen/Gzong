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

type SessionManager struct {
	cookieName  string
	maxLifeTime int64
	lock        sync.Mutex
	sessionMap  map[string]Session
}

type Session struct {
	sessionId string
	lastTime  time.Time
	values    map[interface{}]interface{}
}

func NewSessionManager(cookieName string, maxLifeTime int64) *SessionManager {
	manager := &SessionManager{
		cookieName:  cookieName,
		maxLifeTime: maxLifeTime,
		sessionMap:  make(map[string]Session),
	}
	go manager.SessionGC()
	return manager
}

func (manager *SessionManager) NewSessionId() string {
	bytes := make([]byte, 32)

	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		nano := time.Now().UnixNano()
		return strconv.FormatInt(nano, 10)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

func (manager *SessionManager) NewSession(w http.ResponseWriter, r *http.Request, sessionValues map[interface{}]interface{}) string {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	newSessionId := url.QueryEscape(manager.NewSessionId())
	session := Session{
		sessionId: newSessionId,
		lastTime:  time.Now(),
		values:    sessionValues,
	}
	manager.sessionMap[newSessionId] = session

	cookie := http.Cookie{
		Name:     manager.cookieName,
		Value:    newSessionId,
		Path:     r.URL.Path,
		HttpOnly: true,
		MaxAge:   int(manager.maxLifeTime),
	}
	http.SetCookie(w, &cookie)
	return newSessionId
}

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

func (manager *SessionManager) EndSessionById(sessionId string) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	delete(manager.sessionMap, sessionId)
}

func (manager *SessionManager) GetSessionValue(sessionId string, key interface{}) (interface{}, error) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	if session, ok := manager.sessionMap[sessionId]; ok == true {
		if val, ok := session.values[key]; ok == true {
			return val, nil
		}
	}
	return nil, errors.New("invalid sessionId")
}

func (manager *SessionManager) CheckCookieValid(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie(manager.cookieName)

	if err != nil || cookie == nil {
		return "", err
	}
	manager.lock.Lock()
	defer manager.lock.Unlock()
	sessionId := cookie.Value
	if session, ok := manager.sessionMap[sessionId]; ok == true {
		session.lastTime = time.Now()
		return sessionId, nil
	}
	return "", errors.New("invalid sessionId")
}

func (manager *SessionManager) SessionGC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	for sessionId, session := range manager.sessionMap {
		if session.lastTime.Unix()+manager.maxLifeTime <= time.Now().Unix() {
			delete(manager.sessionMap, sessionId)
		}
	}

	time.AfterFunc(time.Duration(manager.maxLifeTime)*time.Second, manager.SessionGC)
}
