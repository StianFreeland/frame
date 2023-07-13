package comm

import "errors"

var (
	base64DecodeFailed = errors.New("base64 decode failed")
	pemDecodeFailed    = errors.New("pem decode failed")
	rsaInvalidPvtKey   = errors.New("rsa invalid private key")
	rsaDecryptFailed   = errors.New("rsa decrypt failed")
)
