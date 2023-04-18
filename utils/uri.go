package utils

import "net/url"

type uriLib struct{}

var UriLib = uriLib{}

func (uriLib) GetDomain(urlPath string) (string, error) {
	u, err := url.Parse(urlPath)
	if err != nil {
		return "", err
	}
	return u.Hostname(), nil
}
