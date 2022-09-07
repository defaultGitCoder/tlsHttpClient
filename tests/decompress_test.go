package tests

import (
	"fmt"
	"strings"
	"testing"
	"tlsHttpClient/tlsHttpClient"
)

var (
	client = tlsHttpClient.New()
)

type JsonUnmarshalTest struct {
	Slideshow struct {
		Author string `json:"author"`
		Date   string `json:"date"`
		Slides []struct {
			Title string   `json:"title"`
			Type  string   `json:"type"`
			Items []string `json:"items,omitempty"`
		} `json:"slides"`
		Title string `json:"title"`
	} `json:"slideshow"`
}

func TestBrotli(t *testing.T) {
	resp, err := client.R().Get("https://httpbin.org/brotli")
	if err != nil {
		t.Error(err)
	} else {
		if val, ok := resp.Json()["brotli"]; !ok || !val.(bool) {
			t.Error("Brotli not supported")
		} else {
			fmt.Println("- Brotli supported")
		}
	}
}
func TestDeflate(t *testing.T) {
	resp, err := client.R().Get("https://httpbin.org/deflate")
	if err != nil {
		t.Error(err)
	} else {
		if val, ok := resp.Json()["deflated"]; !ok || !val.(bool) {
			t.Error("Deflate not supported")
		} else {
			fmt.Println("- Deflate supported")
		}
	}
}
func TestGzip(t *testing.T) {
	resp, err := client.R().Get("https://httpbin.org/gzip")
	if err != nil {
		t.Error(err)
	} else {
		if val, ok := resp.Json()["gzipped"]; !ok || !val.(bool) {
			t.Error("Gzip not supported")
		} else {
			fmt.Println("- Gzip supported")
		}
	}
}
func TestJson(t *testing.T) {
	resp, err := client.R().Get("https://httpbin.org/json")
	if err != nil {
		t.Error(err)
	} else {
		jsonTest, jsonErr := resp.ToStruct(&JsonUnmarshalTest{})
		if jsonErr != nil {
			t.Error(err)
		} else {
			fmt.Println("- Json unmarshalling supported; Author:", jsonTest.(*JsonUnmarshalTest).Slideshow.Author)
		}
	}
}
func TestUTF8(t *testing.T) {
	resp, err := client.R().Get("https://httpbin.org/encoding/utf8")
	if err != nil {
		t.Error(err)
	} else {
		if !strings.Contains(resp.Text, "<h1>Unicode Demo</h1>") {
			t.Error("UTF-8 not supported")
		} else {
			fmt.Println("- UTF-8 supported")
		}
	}
}
func TestUserAgent(t *testing.T) {
	resp, err := client.R().Get("https://httpbin.org/user-agent")
	if err != nil {
		t.Error(err)
	} else {
		jResp := resp.Json()
		if jResp["user-agent"] != tlsHttpClient.ChromeUserAgent {
			t.Error("User-Agent not invalid")
		} else {
			fmt.Println("- User-Agent valid")
		}
	}
}
