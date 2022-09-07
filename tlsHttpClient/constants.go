package tlsHttpClient

var (
	// Proxy constants
	AvailableSchemas = []string{"http", "https"}

	// Client constants
	defaultHeaders = map[string]string{
		"User-Agent":      ChromeUserAgent,
		"Accept-Encoding": "gzip, deflate, br",
		"Accept":          "*/*",
		"Connection":      "keep-alive",
	}
	defaultTimeout         = 10
	defaultAttempts        = 1
	defaultDisableRedirect = false

	// Fingerprints of browsers
	ChromeUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36"
	ChromeJA3       = "771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-13-18-51-45-43-27-17513-21,29-23-24,0"
)
