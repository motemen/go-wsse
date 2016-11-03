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

	WSSE string

	Transport http.RoundTripper
}

// NewTransport creates a new transport
func NewTransport(username, password string) *Transport {
	t := &Transport{
		Username: username,
		Password: password,
	}
	return t
}

func (t *Transport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
}

// RoundTrip executes the HTTP request with an X-WSSE header.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	wsse, err := CreateHeader(t.Username, t.Password)
	if err != nil {
		return nil, err
	}
	// Shallow copy request
	reqWithHeader := *req
	reqWithHeader.Header = http.Header{}
	for k, v := range req.Header {
		reqWithHeader.Header[k] = v
	}
	reqWithHeader.Header.Set("Authorization", `WSSE profile="UsernameToken"`)
	reqWithHeader.Header.Set("X-WSSE", wsse)
	return t.transport().RoundTrip(&reqWithHeader)
}

// Nonce cannot be any byte because php (or orocrm) cannot handle it
func createNonce() (string, error) {
	nonceByte := make([]byte, nonceSize)
	_, err := rand.Read(nonceByte)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", nonceByte), nil
}

// createCreatedDate creates an ISO 8601 date
func createCreatedDate(t time.Time) string {
	return t.Format("2006-01-02T15:04:05-07:00")
}

func createPasswordDigest(nonce, created, password string) []byte {
	digest := sha1.New()
	digest.Write([]byte(nonce))
	digest.Write([]byte(created))
	digest.Write([]byte(password))
	return digest.Sum(nil)
}

// CreateHeader return the X-WSSE header key as string
func CreateHeader(username, password string) (string, error) {
	nonce, err := createNonce()
	if err != nil {
		return "", err
	}
	created := createCreatedDate(time.Now())
	s := fmt.Sprintf(
		`UsernameToken Username="%s", PasswordDigest="%s", Nonce="%s", Created="%s"`,
		username,
		base64.StdEncoding.EncodeToString(createPasswordDigest(nonce, created, password)),
		base64.StdEncoding.EncodeToString([]byte(nonce)),
		created,
	)
	return s, nil
}
