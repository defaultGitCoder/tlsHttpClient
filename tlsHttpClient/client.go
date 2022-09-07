package tlsHttpClient

import (
	"errors"
	"github.com/quotpw/tlsHttpClient/tlsHttpClient/cycletls"
)

const (
	MethodGet     = "GET"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodDelete  = "DELETE"
	MethodPatch   = "PATCH"
	MethodHead    = "HEAD"
	MethodOptions = "OPTIONS"
)

type RequestProps struct {
	QueryParam      map[string]string
	Headers         map[string]string
	Cookies         []cycletls.Cookie
	DisableRedirect bool
}

type Client struct {
	CycleTLS *cycletls.CycleTLS
	Ja3      string
	Attempts int
	Timeout  int
	props    RequestProps
	proxy    *Proxy
}

//goland:noinspection ALL
func New() *Client {
	CycleTLS := cycletls.New()
	return &Client{
		CycleTLS: &CycleTLS,
		Ja3:      ChromeJA3,
		Attempts: defaultAttempts,
		Timeout:  defaultTimeout,
		props: RequestProps{
			QueryParam:      map[string]string{},
			Headers:         newDefaultHeaders(),
			Cookies:         []cycletls.Cookie{},
			DisableRedirect: defaultDisableRedirect,
		},
		proxy: nil,
	}
}

func (c *Client) SetDisableRedirect(value bool) *Client {
	c.props.DisableRedirect = value
	return c
}

func (c *Client) SetJA3(ja3 string) *Client {
	c.Ja3 = ja3
	return c
}

func (c *Client) SetProxy(proxy *Proxy) error {
	if proxy == nil {
		return errors.New("proxy is nil")
	}
	if !proxy.IsValid() {
		return errors.New("proxy is invalid")
	}
	c.proxy = proxy
	return nil
}

func (c *Client) SetHeader(header, value string) *Client {
	c.props.Headers[header] = value
	return c
}

func (c *Client) SetHeaders(headers map[string]string) *Client {
	for key, value := range headers {
		c.SetHeader(key, value)
	}
	return c
}

func (c *Client) ReplaceHeaders(headers map[string]string) *Client {
	c.props.Headers = make(map[string]string)
	c.SetHeaders(headers)
	return c
}

func (c *Client) SetTimeout(timeout int) *Client {
	c.Timeout = timeout
	return c
}

func (c *Client) SetQueryParams(queryParam map[string]string) {
	for key, value := range queryParam {
		c.props.QueryParam[key] = value
	}
}

func (c *Client) R() *Request {
	return &Request{
		Client:                 c,
		Method:                 "",
		URL:                    "",
		QueryParam:             nil,
		Headers:                map[string]string{},
		Body:                   "",
		Forms:                  nil,
		Json:                   nil,
		Multipart:              nil,
		MultipartBody:          nil,
		DisableRedirect:        c.props.DisableRedirect,
		SetContentTypeDirectly: false,
		Proxy:                  c.proxy,
		Attempts:               c.Attempts,
		Timeout:                c.Timeout,
	}
}

func (c *Client) execute(r *Request) (*Response, error) {
	url := r.ExportUrl()
	body := r.ExportBody()
	headers, userAgent := r.ExportHeaders()
	response, err := c.CycleTLS.Do(
		url,
		cycletls.Options{
			URL:             url,
			Method:          r.Method,
			Headers:         headers,
			Body:            body,
			Ja3:             c.Ja3,
			UserAgent:       userAgent,
			Proxy:           r.ExportProxy(),
			Timeout:         r.Timeout,
			DisableRedirect: r.DisableRedirect,
			HeaderOrder:     nil,
			OrderAsProvided: false,
		},
		r.Method,
	)
	if err != nil {
		return nil, err
	}

	responseObj := &Response{
		Bytes:      response.Bytes,
		Text:       response.Text,
		StatusCode: response.StatusCode,
		Headers:    response.Headers,
		Cookies:    response.Cookies,
	}
	if len(response.Cookies) > 0 {
		c.props.Cookies = append(c.props.Cookies, responseObj.Cookies...)
	}

	return responseObj, nil
}
