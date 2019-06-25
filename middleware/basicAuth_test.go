package middleware

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBaseUser_BasicAuth(t *testing.T) {
	name, pwd := "ssname", "sspwd"
	u := BaseUser{Name: name, Pwd: pwd}
	base64code := Base64Encode(name, pwd)

	hdf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("success"))
	})
	basicAuthFunc := u.BasicAuth(hdf)

	req, err := http.NewRequest("GET", "http://127.0.0.1:9876/", nil)
	w := httptest.NewRecorder()

	if err != nil {
		t.Error("构造请求失败")
	}
	basicAuthFunc.ServeHTTP(w, req)

	bodyBytes, err := ioutil.ReadAll(w.Body)
	if w.Code != http.StatusUnauthorized || string(bodyBytes) != "Unauthorized" {
		t.Error("basicAuth认证失败，请求header未包含认证信息却未能阻断请求")
	}

	req.Header.Add("Authorization", "Basic "+base64code)
	basicAuthFunc.ServeHTTP(w, req)

	bodyBytes, err = ioutil.ReadAll(w.Body)
	if string(bodyBytes) != "success" {
		t.Error("basicAuth认证失败，理应认证通过却没有，应顺利走到业务后续逻辑")
	}
}

func TestBase64Encode(t *testing.T) {
	name, pwd := "ssname", "sspwd"
	authCode := Base64Encode(name, pwd)
	if len(authCode) == 0 {
		t.Error("name和pwd进行base64编码后为空")
	}

	deCodeByte, err := base64.StdEncoding.DecodeString(authCode)
	if err != nil {
		t.Error("认证内容base64解码失败")
	}

	deCode := string(deCodeByte)
	deCodeArray := strings.Split(deCode, ":")
	if name != deCodeArray[0] || pwd != deCodeArray[1] {
		t.Error("解码后的内容与编码前内容不一致")
	}
}
