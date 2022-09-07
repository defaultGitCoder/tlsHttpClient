package tests

import (
	"strings"
	"testing"
)

type T struct {
	Args    map[string]string `json:"args"`
	Data    string            `json:"data"`
	Files   map[string]string `json:"files"`
	Form    map[string]string `json:"form"`
	Headers map[string]string `json:"headers"`
	Json    interface{}       `json:"json"`
	Origin  string            `json:"origin"`
	Url     string            `json:"url"`
}

func TestPostHeaderQueryBody(t *testing.T) {
	resp, err := client.R().
		SetHeader("hello", "world").
		SetQueryParam("hello", "world").
		SetBody("Hello World").
		Post("https://httpbin.org/post")
	if err != nil {
		t.Error(err)
	} else {
		sResp, _ := resp.ToStruct(&T{})
		if sResp.(*T).Data != "Hello World" {
			t.Error("text failed")
		} else {
			t.Log("- text success")
		}

		if sResp.(*T).Args["hello"] != "world" {
			t.Error("args failed")
		} else {
			t.Log("- args success")
		}

		if sResp.(*T).Headers["Hello"] != "world" {
			t.Error("header failed")
		} else {
			t.Log("- header success")
		}
	}
}
func TestPostMultipart(t *testing.T) {
	resp, err := client.R().SetMultipartBoundary("twitter").
		SetMultipartField("hello", "world.jpg", "world", strings.NewReader("world")).
		SetMultipartField("hello1", "world1.jpg", "world1", strings.NewReader("world1")).
		Post("https://httpbin.org/post")
	if err != nil {
		t.Error(err)
	} else {
		sResp, _ := resp.ToStruct(&T{})
		for k, v := range map[string]string{"hello": "world", "hello1": "world1"} {
			if sResp.(*T).Files[k] != v {
				t.Error("multipart", k, v, "failed")
			} else {
				t.Log("- multipart", k, v, "success")
			}
		}
	}
}
func TestPostJson(t *testing.T) {
	jExample := map[string]any{"hello": "world", "number": "one"}
	resp, err := client.R().SetJsonData(jExample).Post("https://httpbin.org/post")
	if err != nil {
		t.Error(err)
	} else {
		sResp, _ := resp.ToStruct(&T{})
		for k, v := range jExample {
			if sResp.(*T).Json.(map[string]any)[k] != v {
				t.Error("json", k, v, "failed")
			} else {
				t.Log("- json", k, v, "success")
			}
		}
	}
}

func TestPostForm(t *testing.T) {
	fExample := map[string]string{"hello": "world", "number": "one"}
	resp, err := client.R().SetFormData(fExample).Post("https://httpbin.org/post")
	if err != nil {
		t.Error(err)
	} else {
		sResp, _ := resp.ToStruct(&T{})
		for k, v := range fExample {
			if sResp.(*T).Form[k] != v {
				t.Error("form", k, v, "failed")
			} else {
				t.Log("- form", k, v, "success")
			}
		}
	}
}
