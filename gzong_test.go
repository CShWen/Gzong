package gzong

import (
	"github.com/cshwen/gzong/middleware"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"strconv"
	"log"
)

const contentType = "Content-Type"
const contentTypeTextHTML = "text/html;charset=UTF-8"
const contentTypeApplicationJSON = "application/json"
const success = "success"
const postBody = `{"test": "ss"}`

func TestNew(t *testing.T) {
	gz := New()
	if gz == nil {
		t.Error("gzong New()不应返回为空")
	}
}

func testPost(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != err {
		log.Println(err)
	}
	w.Header().Set(contentType, contentTypeApplicationJSON)
	w.WriteHeader(http.StatusOK)
	w.Write(bodyBytes)
}

func testGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(contentType, contentTypeTextHTML)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(success))
}

func TestRouter_Add(t *testing.T) {
	gz := New()
	gz.GET("/get", testGet)
	gz.POST("/post", testGet)
	gz.PUT("/put", testGet)
	gz.Add("/add", "GET", testGet)
	gz.Add("/add", "POST", testGet)

	if len(gz.handlersMap) != 4 {
		t.Error("路由个数不符合预期构建数目")
	}
	routeNum := 0
	for _, handlers := range gz.handlersMap {
		routeNum += len(handlers)
	}
	if routeNum != 5 {
		t.Error("路由总个数不符合预期构建数目")
	}
}

func TestRouter_AddMiddleware(t *testing.T) {
	gz := New()
	gz.GET("/test", testGet)

	tsOk := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gz.ServeHTTP(w, r)
	}))
	defer tsOk.Close()

	resp, err := http.Get(tsOk.URL + "/test")
	if err != err {
		t.Log(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Error("GET请求存在的地址，响应返回的状态码不符合预期")
	}

	name, pwd := "ss", "pwd"
	u := middleware.BaseUser{Name: name, Pwd: pwd}
	gz.AddMiddleware(u.BasicAuth)

	tsUn := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gz.ServeHTTP(w, r)
	}))
	defer tsUn.Close()

	resp, err = http.Get(tsUn.URL + "/test")
	if err != err {
		t.Log(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Error("请求添加了中间件basicAuth的服务，未携带认证信息应不通过")
	}

	req, err := http.NewRequest("GET", tsUn.URL+"/test", nil)
	if err != err {
		t.Log(err)
	}
	req.Header.Add("Authorization", "Basic "+middleware.Base64Encode(name, pwd))
	resp, err = http.DefaultClient.Do(req)
	if err != err {
		t.Log(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Error("请求添加了中间件basicAuth的服务，携带认证信息应顺利通过")
	}
}

func TestRouter_POST(t *testing.T) {
	gz := New()
	gz.POST("/test", testPost)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gz.ServeHTTP(w, r)
	}))
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/test", contentTypeApplicationJSON, strings.NewReader(postBody))
	if err != err {
		t.Log(err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != err {
		t.Log(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Error("POST请求存在的地址，响应返回的状态码不符合预期")
	}
	if resp.Header.Get(contentType) != contentTypeApplicationJSON {
		t.Error("POST请求存在的地址，响应返回的header不符合预期")
	}
	if string(bodyBytes) != postBody {
		t.Error("POST请求存在的地址，响应返回的body不符合预期")
	}
}

func TestRouter_PUT(t *testing.T) {
	gz := New()
	gz.PUT("/test", testPost)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gz.ServeHTTP(w, r)
	}))
	defer ts.Close()

	req, err := http.NewRequest("PUT", ts.URL+"/test", strings.NewReader(postBody))
	if err != err {
		t.Log(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != err {
		t.Log(err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != err {
		t.Log(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Error("PUT请求存在的地址，响应返回的状态码不符合预期")
	}
	if resp.Header.Get(contentType) != contentTypeApplicationJSON {
		t.Error("PUT请求存在的地址，响应返回的header不符合预期")
	}
	if string(bodyBytes) != postBody {
		t.Error("PUT请求存在的地址，响应返回的body不符合预期")
	}
}

func TestRouter_ServeHTTP(t *testing.T) {
	gz := New()
	gz.GET("/test", testGet)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gz.ServeHTTP(w, r)
	}))
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/test")
	if err != err {
		t.Log(err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != err {
		t.Log(err)
	}
	if resp.StatusCode != http.StatusOK || string(bodyBytes) != success {
		t.Error("GET请求存在的地址访问异常")
	}

	resp, err = http.Get(ts.URL + "/error")
	if err != err {
		t.Log(err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Error("GET请求不存在的地址应404")
	}
}

func TestRouter_Run(t *testing.T) {
	port := 9872
	strPort := ":" + strconv.Itoa(port)

	gz := New()
	gz.GET("/test", testGet)

	go func() {
		gz.Run(strPort)
		time.Sleep(1 * time.Second)
		//defer gz.Close()
	}()

	resp, err := http.Get("http://127.0.0.1" + strPort + "/test")
	t.Log("sstest resp:\t", resp)
	if err != err {
		t.Log(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Error("GET请求存在的地址访问异常")
	}

	resp, err = http.Get("http://127.0.0.1" + strPort + "/error")
	if err != err {
		t.Log(err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Error("GET请求不存在的地址应404")
	}
}

func TestRouter_GET(t *testing.T) {
	port := 9871
	strPort := ":" + strconv.Itoa(port)

	gz := New()
	gz.GET("/test", testGet)
	go func() {
		gz.Run(strPort)
		time.Sleep(1 * time.Second)
		//defer gz.Close()
	}()

	resp, err := http.Get("http://127.0.0.1" + strPort + "/test")
	if err != err {
		t.Log(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Error("GET请求存在的地址，响应返回的状态码不符合预期")
	}

	if resp.Header.Get(contentType) != contentTypeTextHTML {
		t.Error("GET请求存在的地址，响应返回的header不符合预期")
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != err {
		t.Log(err)
	}
	if string(bodyBytes) != success {
		t.Error("GET请求存在的地址，响应返回的body不符合预期")
	}

	resp, err = http.Get("http://127.0.0.1" + strPort + "/error")
	if err != err {
		t.Log(err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Error("GET请求不存在的地址应404")
	}
}
