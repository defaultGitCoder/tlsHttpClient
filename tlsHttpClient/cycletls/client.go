package cycletls

import (
	http "github.com/Danny-Dasilva/fhttp"
	"github.com/Danny-Dasilva/fhttp/cookiejar"

	"time"

	"golang.org/x/net/proxy"
)

type browser struct {
	JA3       string
	UserAgent string
}

var disabledRedirect = func(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

func clientBuilder(browser browser, dialer proxy.ContextDialer, timeout int, disableRedirect bool, jar *cookiejar.Jar) http.Client {
	//if timeout is not set in call default to 15
	if timeout == 0 {
		timeout = 15
	}
	client := http.Client{
		Transport: newRoundTripper(browser, dialer),
		Timeout:   time.Duration(timeout) * time.Second,
		Jar:       jar,
	}
	//if disableRedirect is set to true httpclient will not redirect
	if disableRedirect {
		client.CheckRedirect = disabledRedirect
	}
	return client
}

// newClient creates a new http client
func newClient(browser browser, timeout int, disableRedirect bool, UserAgent string, proxyURL string, jar *cookiejar.Jar) (http.Client, error) {
	if len(proxyURL) > 0 && len(proxyURL) > 0 {
		dialer, err := newConnectDialer(proxyURL, UserAgent)
		if err != nil {
			return http.Client{
				Timeout:       time.Duration(timeout) * time.Second,
				CheckRedirect: disabledRedirect,
				Jar:           jar,
			}, err
		}
		return clientBuilder(
			browser,
			dialer,
			timeout,
			disableRedirect,
			jar,
		), nil
	}

	return clientBuilder(
		browser,
		proxy.Direct,
		timeout,
		disableRedirect,
		jar,
	), nil

}
