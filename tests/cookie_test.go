package tests_test

import (
	"testing"
	"tlsHttpClient/tlsHttpClient"
)

var (
	client  = tlsHttpClient.New()
	cookies = map[string]string{"Foo": "Bar", "Baz": "Qux"}
)

func TestCookie(t *testing.T) {
	// RECEIVE COOKIES
	response, err := client.R().SetQueryParams(cookies).Get("https://httpbin.org/cookies/set")
	if err != nil {
		t.Error(err)
		return
	}

	if len(response.Cookies) != len(cookies) {
		t.Error("Received cookies count mismatch; Received: ", len(response.Cookies), " Expected: ", len(cookies))
	}

	for _, cookie := range response.Cookies {
		if cookie.Value != cookies[cookie.Name] {
			t.Error("Received cookie value mismatch; Received: ", cookie.Value, " Expected: ", cookies[cookie.Name])
		}
	}

	// GET COOKIES
	response, err = client.R().Get("https://httpbin.org/cookies")
	if err != nil {
		t.Error(err)
	}
	resp := response.Json()
	for v, k := range cookies {
		if resp["cookies"].(map[string]interface{})[v] != k {
			t.Error("Cookie value mismatch; Received: ", resp["cookies"].(map[string]interface{})[v], " Expected: ", k)
		}
	}
}
