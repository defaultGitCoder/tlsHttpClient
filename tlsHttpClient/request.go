package tlsHttpClient

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/url"
	"strings"
)

type Request struct {
	Client *Client

	Method string

	URL        string
	QueryParam map[string]string

	Headers map[string]string

	Body          string
	Forms         *map[string]string
	Json          *map[string]any
	Multipart     *multipart.Writer
	MultipartBody *bytes.Buffer

	DisableRedirect bool

	SetContentTypeDirectly bool

	Proxy *Proxy

	Attempts int
	Timeout  int
}

func (r *Request) ExportProxy() string {
	if r.Proxy != nil {
		return r.Proxy.ToUrl()
	}
	return ""
}

func (r *Request) ExportHeaders() (map[string]string, string) {
	headers := make(map[string]string)

	for k, v := range r.Client.Props.Headers {
		headers[k] = v
	}
	for k, v := range r.Headers {
		headers[k] = v
	}

	userAgent := ChromeUserAgent
	for k, v := range headers {
		if strings.ToLower(k) == "user-agent" {
			userAgent = v
			break
		}
	}

	return headers, userAgent
}

func (r *Request) ExportUrl() string {
	result := r.URL
	if len(r.QueryParam) > 0 || len(r.Client.Props.QueryParam) > 0 {
		if strings.Contains(result, "?") {
			if !strings.HasSuffix(result, "&") {
				result += "&"
			}
		} else {
			result += "?"
		}

		params := url.Values{}
		if r.Client.Props.QueryParam != nil {
			for k, v := range r.Client.Props.QueryParam {
				params.Set(k, v)
			}
		}
		if r.QueryParam != nil {
			for k, v := range r.QueryParam {
				params.Set(k, v)
			}
		}
		result += params.Encode()
	}
	return result
}

func (r *Request) ExportBody() string {
	body, contentType, setBody := prepareBody(r)
	if setBody {
		if !r.SetContentTypeDirectly {
			r.SetHeader("Content-Type", contentType)
		}
	}
	return body
}

func (r *Request) SetHeader(header, value string) *Request {
	r.Headers[header] = value
	return r
}

func (r *Request) SetHeaders(headers map[string]string) *Request {
	for h, v := range headers {
		r.SetHeader(h, v)
	}
	return r
}

func (r *Request) SetQueryParam(param, value string) *Request {
	if r.QueryParam == nil {
		r.QueryParam = map[string]string{}
	}
	r.QueryParam[param] = value
	return r
}

func (r *Request) SetQueryParams(params map[string]string) *Request {
	for p, v := range params {
		r.SetQueryParam(p, v)
	}
	return r
}
func (r *Request) SetQueryString(query string) *Request {
	params, err := url.ParseQuery(strings.TrimSpace(query))
	if err == nil {
		for p, v := range params {
			for _, pv := range v {
				r.SetQueryParam(p, pv)
			}
		}
	}
	return r
}

func (r *Request) SetFormData(data map[string]string) *Request {
	if r.Forms == nil {
		r.Forms = &map[string]string{}
	}
	for k, v := range data {
		(*r.Forms)[k] = v
	}
	return r
}
func (r *Request) SetBody(body string) *Request {
	r.Body = body
	return r
}

func (r *Request) SetMultipartField(param, fileName, contentType string, reader io.Reader) *Request {
	if r.Multipart == nil {
		r.MultipartBody = &bytes.Buffer{}
		r.Multipart = multipart.NewWriter(r.MultipartBody)
	}
	header := NewMultipartFieldHeader(param, fileName, contentType)
	part, _ := r.Multipart.CreatePart(header)
	_, _ = io.Copy(part, reader)
	return r
}

func (r *Request) SetMultipartFormData(data map[string]string) *Request {
	if r.Multipart == nil {
		r.MultipartBody = &bytes.Buffer{}
		r.Multipart = multipart.NewWriter(r.MultipartBody)
	}
	for k, v := range data {
		_ = r.Multipart.WriteField(k, v)
	}

	return r
}

func (r *Request) SetMultipartBoundary(boundary string) *Request {
	if r.Multipart == nil {
		r.MultipartBody = &bytes.Buffer{}
		r.Multipart = multipart.NewWriter(r.MultipartBody)
	}
	err := r.Multipart.SetBoundary(boundary)
	if err != nil {
		log.Println("SetMultipartBoundary:", err)
	}
	return r
}

func (r *Request) SetJsonData(data map[string]any) *Request {
	if r.Json == nil {
		r.Json = &map[string]any{}
	}
	for k, v := range data {
		(*r.Json)[k] = v
	}
	return r
}

func (r *Request) SetContentType(contentType string) *Request {
	r.SetContentTypeDirectly = true
	r.SetHeader("Content-Type", contentType)
	return r
}

func (r *Request) SetProxy(proxy Proxy) *Request {
	r.Proxy = &proxy
	return r
}

func (r *Request) Execute(method, url string) (*Response, error) {
	r.Method = method
	r.URL = url

	if r.Multipart != nil && !(method == MethodPost || method == MethodPut || method == MethodPatch) {
		return nil, fmt.Errorf("multipart content is not allowed in HTTP verb [%v]", method)
	}

	var err error
	var resp *Response

	if r.Client.Attempts == 0 {
		r.Attempts = 1
	}

	for i := 0; i < r.Attempts; i++ {
		resp, err = r.Client.execute(r)
		if err == nil {
			break
		}
	}

	return resp, err
}

func (r *Request) Get(url string) (*Response, error) {
	return r.Execute(MethodGet, url)
}

func (r *Request) Head(url string) (*Response, error) {
	return r.Execute(MethodHead, url)
}

func (r *Request) Post(url string) (*Response, error) {
	return r.Execute(MethodPost, url)
}

func (r *Request) Put(url string) (*Response, error) {
	return r.Execute(MethodPut, url)
}

func (r *Request) Delete(url string) (*Response, error) {
	return r.Execute(MethodDelete, url)
}

func (r *Request) Options(url string) (*Response, error) {
	return r.Execute(MethodOptions, url)
}

func (r *Request) Patch(url string) (*Response, error) {
	return r.Execute(MethodPatch, url)
}

func (r *Request) Send() (*Response, error) {
	return r.Execute(r.Method, r.URL)
}
