package wsse

import (
	"crypto/sha1"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// http.Header{"X-Wsse":[]string{"UsernameToken Username=\"user\", PasswordDigest=\"76LhUHaHJV7p6/HzaXWJ+wTUSxM=\", Nonce=\"OIKJlY0HGDY4uuW0\", Created=\"2014-12-01T12:38:55Z\""}}

var (
	rxSplitHeader = regexp.MustCompile(`\s*,\s*`)
	rxKeyValue    = regexp.MustCompile(`^(\w+)="(.+)"$`)
)

func TestRoundTrip(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, `WSSE profile="UsernameToken"`, r.Header.Get("Authorization"))

		wsseHeader := r.Header.Get("X-WSSE")
		assert.NotEmpty(t, wsseHeader)
		assert.True(t, strings.HasPrefix(wsseHeader, "UsernameToken "))

		kv := map[string]string{}

		wsseHeader = strings.TrimPrefix(wsseHeader, "UsernameToken ")
		parts := rxSplitHeader.Split(wsseHeader, -1)
		for _, part := range parts {
			m := rxKeyValue.FindStringSubmatch(part)
			if m != nil {
				kv[m[1]] = m[2]
			}
		}

		assert.Equal(t, kv["Username"], "user")

		nonceBytes, err := base64.StdEncoding.DecodeString(kv["Nonce"])
		assert.NoError(t, err)

		digest := sha1.New()
		digest.Write(nonceBytes)
		digest.Write([]byte(kv["Created"]))
		digest.Write([]byte("pass"))

		assert.Equal(t, kv["PasswordDigest"], base64.StdEncoding.EncodeToString(digest.Sum(nil)))

		w.Write([]byte{})
	}

	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	client := http.Client{
		Transport: NewTransport("user", "pass"),
	}
	_, err := client.Get(ts.URL)

	assert.NoError(t, err)
}

func TestNewTransport(t *testing.T) {
	trans := NewTransport("user", "pass")
	assert.Equal(t, "user", trans.Username)
	assert.Equal(t, "pass", trans.Password)
}

type testRT struct{}

func (t *testRT) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, nil
}

func TestTransport(t *testing.T) {
	trans2 := NewTransport("user", "pass")
	ttRT := &testRT{}
	trans2.Transport = ttRT
	trans2.transport()
	assert.Equal(t, ttRT, trans2.Transport)
}

func TestCreatePasswordDigest(t *testing.T) {
	digest := createPasswordDigest(
		"cc244d8c8a54cfda",
		"2014-12-01T03:39:38Z",
		"pass",
	)
	base64Digest := base64.StdEncoding.EncodeToString(digest)
	assert.Equal(t, base64Digest, "8OjkyL8RK7/vse443STJVoOc7hw=")
}

func TestCreateNonce(t *testing.T) {
	n1, err1 := createNonce()
	n2, err2 := createNonce()
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, n1, n2)
}

func TestCreateCreatedDate(t *testing.T) {
	assert.Equal(t, "2001-09-22 11:01:33 +0000 UTC", time.Date(2001, 9, 22, 11, 1, 33, 0, time.UTC).String())
}
