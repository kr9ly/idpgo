package idpgo

import "net/url"

type Query map[string]string

func (q Query) String() string {
	buffer := make([]byte, 0, 100)
	isFirst := true
	for key, value := range q {
		if !isFirst {
			buffer = append(buffer, "&"...)
		}
		buffer = append(buffer, url.QueryEscape(key)...)
		buffer = append(buffer, "="...)
		buffer = append(buffer, url.QueryEscape(value)...)
		isFirst = false
	}
	return string(buffer)
}
