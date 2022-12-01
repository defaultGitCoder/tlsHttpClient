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
		Scheme:   &scheme,
		Host:     &listProxy[0],
		Port:     &listProxy[1],
		Username: nil,
		Password: nil,
	}

	if len(listProxy) >= 4 {
		p.Username = &listProxy[2]
		p.Password = &listProxy[3]
	}

	return p
}

func (p Proxy) haveAuth() bool {
	return p.Username != nil && p.Password != nil
}

func (p Proxy) IsValid() bool {
	return p.Scheme != nil && p.Host != nil && p.Port != nil
}

func (p Proxy) ToUrl() string {
	if !p.IsValid() {
		panic("Proxy is not valid")
	}

	url := *p.Scheme + "://"

	if p.haveAuth() {
		url += *p.Username + ":" + *p.Password + "@"
	}

	url += *p.Host + ":" + *p.Port

	return url
}
