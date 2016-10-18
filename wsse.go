package wsse

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

const (
	nonceSize = 16
)

// Transport implements http.Transport. It adds X-WSSE header
// to client requests.
type Transport struct {
	Username string
	Password string

	Transport http.RoundTripper
}

func (t *Transport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
}

// RoundTrip executes the HTTP request with an X-WSSE header.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	nonceByte := make([]byte, nonceSize)
	_, err := rand.Read(nonceByte)
	if err != nil {
		return nil, err
	}
	nonce := fmt.Sprintf("%x", nonceByte)
	// ISO 8601
	created := time.Now().Format("2006-01-02T15:04:05-07:00")
	// Shallow copy request
	reqWithHeader := *req
	reqWithHeader.Header = http.Header{}
	for k, v := range req.Header {
		reqWithHeader.Header[k] = v
	}
	reqWithHeader.Header.Set(
		"X-WSSE",
		fmt.Sprintf(
			`UsernameToken Username="%s", PasswordDigest="%s", Nonce="%s", Created="%s"`,
			t.Username,
			base64.StdEncoding.EncodeToString(createPasswordDigest(nonce, created, t.Password)),
			base64.StdEncoding.EncodeToString([]byte(nonce)),
			created,
		),
	)
	return t.transport().RoundTrip(&reqWithHeader)
}

func createPasswordDigest(nonce, created, password string) []byte {
	digest := sha1.New()
	digest.Write([]byte(nonce))
	digest.Write([]byte(created))
	digest.Write([]byte(password))
	return digest.Sum(nil)
}
