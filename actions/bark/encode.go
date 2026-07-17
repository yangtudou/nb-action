package bark

import (
	"encoding/base64"
	"net/url"
)

func encodeCiphertext(data []byte) string {

	base64Text :=
		base64.StdEncoding.EncodeToString(data)

	urlText :=
		url.QueryEscape(base64Text)

	return urlText
}
