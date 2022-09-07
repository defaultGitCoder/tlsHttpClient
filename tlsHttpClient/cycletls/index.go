package cycletls

import (
	http "github.com/Danny-Dasilva/fhttp"
	"github.com/Danny-Dasilva/fhttp/cookiejar"
	"io"
	"log"
	"net/url"
	"strings"
)

// Options sets CycleTLS client options
type Options struct {
	URL             string
	Method          string
	Headers         map[string]string
	Body            string
	Ja3             string
	UserAgent       string
	Proxy           string
	Timeout         int
	DisableRedirect bool
	HeaderOrder     []string
	OrderAsProvided bool
}

type cycleTLSRequest struct {
	RequestID string  `json:"requestId"`
	Options   Options `json:"options"`
	jar       *cookiejar.Jar
}

// rename to request+client+options
type fullRequest struct {
	req     *http.Request
	client  http.Client
	options cycleTLSRequest
}

// Response contains CycleTLS response data
type Response struct {
	Headers    map[string]string
	Cookies    []Cookie
	StatusCode int
	Bytes      []byte
	Text       string
}

// CycleTLS creates full request and response
type CycleTLS struct {
	ReqChan   chan fullRequest
	RespChan  chan Response
	CookieJar *cookiejar.Jar
}

// ready Request
func processRequest(request cycleTLSRequest) (result fullRequest) {
	var browser = browser{
		JA3:       request.Options.Ja3,
		UserAgent: request.Options.UserAgent,
	}

	client, err := newClient(
		browser,
		request.Options.Timeout,
		request.Options.DisableRedirect,
		request.Options.UserAgent,
		request.Options.Proxy,
		request.jar,
	)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(strings.ToUpper(request.Options.Method), request.Options.URL, strings.NewReader(request.Options.Body))
	if err != nil {
		log.Fatal(err)
	}
	var headerOrder []string
	//master header order, all your headers will be ordered based on this list and anything extra will be appended to the end
	//if your site has any custom headers, see the header order chrome uses and then add those headers to this list
	if len(request.Options.HeaderOrder) > 0 {
		//lowercase headers
		for _, v := range request.Options.HeaderOrder {
			lowerCaseKey := strings.ToLower(v)
			headerOrder = append(headerOrder, lowerCaseKey)
		}
	} else {
		headerOrder = append(headerOrder,
			"host",
			"connection",
			"cache-control",
			"device-memory",
			"viewport-width",
			"rtt",
			"downlink",
			"ect",
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"sec-ch-ua-full-version",
			"sec-ch-ua-arch",
			"sec-ch-ua-platform",
			"sec-ch-ua-platform-version",
			"sec-ch-ua-model",
			"upgrade-insecure-requests",
			"user-agent",
			"accept",
			"sec-fetch-site",
			"sec-fetch-mode",
			"sec-fetch-user",
			"sec-fetch-dest",
			"referer",
			"accept-encoding",
			"accept-language",
			"cookie",
		)
	}

	headerMap := make(map[string]string)
	var headerOrderKey []string
	for _, key := range headerOrder {
		for k, v := range request.Options.Headers {
			lowerCaseKey := strings.ToLower(k)
			if key == lowerCaseKey {
				headerMap[k] = v
				headerOrderKey = append(headerOrderKey, lowerCaseKey)
			}
		}

	}

	//ordering the pseudo headers and our normal headers
	req.Header = http.Header{
		http.HeaderOrderKey:  headerOrderKey,
		http.PHeaderOrderKey: {":method", ":authority", ":scheme", ":path"},
	}
	//set our Host header
	u, err := url.Parse(request.Options.URL)
	if err != nil {
		panic(err)
	}

	//append our normal headers
	for k, v := range request.Options.Headers {
		if k != "Content-Length" {
			req.Header.Set(k, v)
		}
	}
	req.Header.Set("Host", u.Host)
	req.Header.Set("user-agent", request.Options.UserAgent)
	return fullRequest{req: req, client: client, options: request}

}

func dispatcher(res fullRequest) (response Response, err error) {
	resp, err := res.client.Do(res.req)
	if err != nil {
		parsedError := parseError(err)
		response := parsedError.ErrorMsg + "-> \n" + err.Error()
		return Response{
			Headers:    map[string]string{},
			Cookies:    []Cookie{},
			StatusCode: parsedError.StatusCode,
			Bytes:      []byte(response),
			Text:       response,
		}, err //normally return error here
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	encoding := resp.Header["Content-Encoding"]

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	text, bytes := DecompressBody(bodyBytes, encoding)
	headers := make(map[string]string)

	for name, values := range resp.Header {
		if name == "Set-Cookie" {
			headers[name] = strings.Join(values, "/,/")
		} else {
			for _, value := range values {
				headers[name] = value
			}
		}
	}

	var cookies []Cookie
	for _, v := range res.client.Jar.Cookies(res.req.URL) {
		cookies = append(cookies, Cookie{
			Name:    v.Name,
			Value:   v.Value,
			Path:    v.Path,
			Domain:  v.Domain,
			Expires: v.Expires,
			JSONExpires: Time{
				Time: v.Expires,
			},
			RawExpires: v.RawExpires,
			MaxAge:     v.MaxAge,
			Secure:     v.Secure,
			HTTPOnly:   v.HttpOnly,
			SameSite:   v.SameSite,
			Raw:        v.Raw,
			Unparsed:   v.Unparsed,
		})
	}

	return Response{
		Headers:    headers,
		Cookies:    cookies,
		StatusCode: resp.StatusCode,
		Bytes:      bytes,
		Text:       text,
	}, nil

}

// Queue queues request in worker pool
func (client *CycleTLS) Queue(URL string, options Options, Method string) {

	options.URL = URL
	options.Method = Method
	//TODO add timestamp to request
	opt := cycleTLSRequest{"Queued Request", options, client.CookieJar}
	response := processRequest(opt)
	client.ReqChan <- response
}

// Do creates a single request
func (client *CycleTLS) Do(URL string, options Options, Method string) (response Response, err error) {

	options.URL = URL
	options.Method = Method
	opt := cycleTLSRequest{"cycleTLSRequest", options, client.CookieJar}

	res := processRequest(opt)
	response, err = dispatcher(res)
	if err != nil {
		return response, err
	}

	return response, nil
}

//TODO rename this

func getNewJar() *cookiejar.Jar {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	return jar
}

func New(workers ...bool) CycleTLS {
	if len(workers) > 0 && workers[0] {
		reqChan := make(chan fullRequest)
		respChan := make(chan Response)
		go workerPool(reqChan, respChan)

		return CycleTLS{ReqChan: reqChan, RespChan: respChan, CookieJar: getNewJar()}
	}
	return CycleTLS{CookieJar: getNewJar()}

}

// Close closes channels
func (client *CycleTLS) Close() {
	close(client.ReqChan)
	close(client.RespChan)

}

// Worker Pool
func workerPool(reqChan chan fullRequest, respChan chan Response) {
	//MAX
	for i := 0; i < 100; i++ {
		go worker(reqChan, respChan)
	}
}

// Worker
func worker(reqChan chan fullRequest, respChan chan Response) {
	for res := range reqChan {
		response, _ := dispatcher(res)
		respChan <- response
	}
}
