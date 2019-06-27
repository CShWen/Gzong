package gzong

import (
	"github.com/cshwen/gzong/middleware"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
	bodyBytes, _ := ioutil.ReadAll(r.Body)
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

	resp, _ := http.Get(tsOk.URL + "/test")
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

	resp, _ = http.Get(tsUn.URL + "/test")
	if resp.StatusCode != http.StatusUnauthorized {
		t.Error("请求添加了中间件basicAuth的服务，未携带认证信息应不通过")
	}

	req, _ := http.NewRequest("GET", tsUn.URL+"/test", nil)
	req.Header.Add("Authorization", "Basic "+middleware.Base64Encode(name, pwd))
	resp, _ = http.DefaultClient.Do(req)

	if resp.StatusCode != http.StatusOK {
		t.Error("请求添加了中间件basicAuth的服务，携带认证信息应顺利通过")
	}
}

func TestRouter_GET(t *testing.T) {
	gz := New()
	gz.GET("/test", testGet)
	gz.Run(":9871")

	resp, _ := http.Get("http://127.0.0.1:9871/test")
	if resp.StatusCode != http.StatusOK {
		t.Error("GET请求存在的地址，响应返回的状态码不符合预期")
	}
	if resp.Header.Get(contentType) != contentTypeTextHTML {
		t.Error("GET请求存在的地址，响应返回的header不符合预期")
	}
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	if string(bodyBytes) != success {
		t.Error("GET请求存在的地址，响应返回的body不符合预期")
	}
}

func TestRouter_POST(t *testing.T) {
	gz := New()
	gz.POST("/test", testPost)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gz.ServeHTTP(w, r)
	}))
	defer ts.Close()

	resp, _ := http.Post(ts.URL+"/test", contentTypeApplicationJSON, strings.NewReader(postBody))
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

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

	req, _ := http.NewRequest("PUT", ts.URL+"/test", strings.NewReader(postBody))
	resp, _ := http.DefaultClient.Do(req)
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

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

	resp, _ := http.Get(ts.URL + "/test")
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK || string(bodyBytes) != success {
		t.Error("GET请求存在的地址访问异常")
	}

	resp, _ = http.Get(ts.URL + "/error")
	if resp.StatusCode != http.StatusNotFound {
		t.Error("GET请求不存在的地址应404")
	}
}

func TestRouter_Run(t *testing.T) {
	gz := New()
	gz.GET("/test", testGet)
	gz.Run(":9872")

	resp, _ := http.Get("http://127.0.0.1:9872/test")
	if resp.StatusCode != http.StatusOK {
		t.Error("GET请求存在的地址访问异常")
	}
	resp, _ = http.Get("http://127.0.0.1:9872/error")
	if resp.StatusCode != http.StatusNotFound {
		t.Error("GET请求不存在的地址应404")
	}
	go func() { defer gz.Close() }()
}
