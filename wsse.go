package wsse

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
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

var nonceSize int = 16

// RoundTrip executes the HTTP request with an X-WSSE header.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	nonce := make([]byte, nonceSize)
	_, err := rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	created := time.Now().Format("2006-01-02T15:04:05Z")

	// Shallow copy request
	var reqWithHeader http.Request = *req
	reqWithHeader.Header = http.Header{}
	for k, v := range req.Header {
		reqWithHeader.Header[k] = v
	}

	reqWithHeader.Header.Set(
		"X-WSSE",
		fmt.Sprintf(
			`UsernameToken Username="%s", PasswordDigest="%s", Nonce="%s", Created="%s"`,
			t.Username,
			base64.StdEncoding.EncodeToString(createPasswordDigest(string(nonce), created, t.Password)),
			base64.StdEncoding.EncodeToString(nonce),
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
