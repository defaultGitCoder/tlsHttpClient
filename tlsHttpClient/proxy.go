package tlsHttpClient

import (
	"strings"
)

type Proxy struct {
	Scheme *string

	Host *string
	Port *string

	Username *string
	Password *string
}

//goland:noinspection GoUnusedExportedFunction
func StringToProxy(strProxy string, scheme string) *Proxy {
	if scheme == "" {
		scheme = AvailableSchemas[0]
	} else {
		if !StringInStringArray(scheme, AvailableSchemas) {
			panic("Invalid scheme, scheme must be one of: " + strings.Join(AvailableSchemas, ", "))
		}
	}

	listProxy := strings.Split(strProxy, ":")
	if len(listProxy) < 2 {
		return nil
	}

	p := &Proxy{
		scheme:   &scheme,
		host:     &listProxy[0],
		port:     &listProxy[1],
		username: nil,
		password: nil,
	}

	if len(listProxy) >= 4 {
		p.username = &listProxy[2]
		p.password = &listProxy[3]
	}

	return p
}

func (p Proxy) haveAuth() bool {
	return p.username != nil && p.password != nil
}

func (p Proxy) IsValid() bool {
	return p.scheme != nil && p.host != nil && p.port != nil
}

func (p Proxy) ToUrl() string {
	if !p.IsValid() {
		panic("Proxy is not valid")
	}

	url := *p.scheme + "://"

	if p.haveAuth() {
		url += *p.username + ":" + *p.password + "@"
	}

	url += *p.host + ":" + *p.port

	return url
}
