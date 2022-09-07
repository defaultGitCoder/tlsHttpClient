package tlsHttpClient

import (
	"encoding/json"
	"fmt"
	"net/textproto"
	"net/url"
	"strings"
)

var (
	escapeQuotes = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")
)

func NewMultipartFieldHeader(fieldName, filename string, contentType string) textproto.MIMEHeader {
	h := make(textproto.MIMEHeader)
	if filename == "" {
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"`, escapeQuotes.Replace(fieldName)))
	} else {
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes.Replace(fieldName), escapeQuotes.Replace(filename)))
	}
	if contentType != "" {
		h.Set("Content-Type", contentType)
	}
	return h
}

func prepareFormBody(forms map[string]string) string {
	values := url.Values{}
	for k, v := range forms {
		values.Add(k, v)
	}
	return values.Encode()
}

func prepareBody(r *Request) (string, string, bool) {
	if r.Body != "" {
		return r.Body, "text/plain", true
	}
	if r.Forms != nil {
		return prepareFormBody(*r.Forms), "application/x-www-form-urlencoded", true
	}
	if r.Json != nil {
		j, _ := json.Marshal(*r.Json)
		return string(j), "application/json", true
	}
	if r.Multipart != nil {
		return r.MultipartBody.String() + "\r\n--" + r.Multipart.Boundary() + "--\r\n", r.Multipart.FormDataContentType(), true
	}
	return "", "", false
}

func Unmarshal(text []byte, v interface{}) error {
	return json.Unmarshal(text, v)
}

func StringInStringArray(value string, array []string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

func newDefaultHeaders() map[string]string {
	headers := map[string]string{}
	for k, v := range defaultHeaders {
		headers[k] = v
	}
	return headers
}
