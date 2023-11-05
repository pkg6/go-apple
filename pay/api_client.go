package pay

import (
	"context"
	"crypto/ecdsa"
	"crypto/tls"
	"github.com/golang-jwt/jwt"
	"github.com/pkg6/go-apple/util"
	"github.com/zzqqw/gclient"
	"net"
	"net/http"
	"time"
)

const (
	apiUrl    = "https://api.storekit.itunes.apple.com"
	apiBoxUrl = "https://api.storekit-sandbox.itunes.apple.com"
)

type ResponseErrorMessage struct {
	ErrorCode    int    `json:"errorCode,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// ApiClient

type ApiClient struct {
	//Your app’s bundle ID (exp: “com.example.testbundleid2021”)
	Bid string
	//https://appstoreconnect.apple.com/access/api/subs
	// Your issuer ID from the Keys page in App Store Connect (exp: "57246542-96fe-1a63-e053-0824d011072a")
	Iss string
	//Your private key ID from App Store Connect (Ex: GTK54UHD42)
	KeyID        string
	IsProduction bool
	PrivateKey   *ecdsa.PrivateKey
	Client       gclient.ClientInterface
}

func NewApiClient(iss, bid, keyID string, privateKey []byte, isProduction bool) (api *ApiClient, err error) {
	api = &ApiClient{Bid: bid, Iss: iss, KeyID: keyID, IsProduction: isProduction}
	api.PrivateKey, err = util.EcdsaPrivateKey(privateKey)
	if err != nil {
		return
	}
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}
	api.Client = gclient.NewSetHttpClient(&http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			//DialContext: defaultTransportDialContext(),
			DialContext:           dialer.DialContext,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			DisableKeepAlives:     true,
			ForceAttemptHTTP2:     true,
		},
	})
	if api.IsProduction {
		api.Client.SetBaseURL(apiUrl)
	} else {
		api.Client.SetBaseURL(apiBoxUrl)
	}
	return
}

type CustomClaims struct {
	jwt.Claims
	Iss string `json:"iss"`
	Iat int64  `json:"iat"`
	Exp int64  `json:"exp"`
	Aud string `json:"aud"`
	Bid string `json:"bid"`
}

// BuildJwtToken
//https://developer.apple.com/documentation/appstoreserverapi/generating_tokens_for_api_requests
func (a *ApiClient) BuildJwtToken() (string, error) {
	claims := CustomClaims{
		Iss: a.Iss,
		Iat: time.Now().Unix(),
		Exp: time.Now().Add(5 * time.Minute).Unix(),
		Aud: "appstoreconnect-v1",
		Bid: a.Bid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header = map[string]any{
		"alg": "ES256",
		"kid": a.KeyID,
		"typ": "JWT",
	}
	accessToken, err := token.SignedString(a.PrivateKey)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}
func (a *ApiClient) WithTokenGet(path string, data, d any) error {
	token, err := a.BuildJwtToken()
	if err != nil {
		return err
	}
	a.Client.WithToken(token)
	return a.Client.GetUnmarshal(context.Background(), path, data, d)
}
func (a *ApiClient) WithTokenPost(path string, data, d any) error {
	token, err := a.BuildJwtToken()
	if err != nil {
		return err
	}
	a.Client.WithToken(token)
	return a.Client.AsJson().PostUnmarshal(context.Background(), path, data, d)
}
func (a *ApiClient) WithTokenPut(path string, data, d any) error {
	token, err := a.BuildJwtToken()
	if err != nil {
		return err
	}
	a.Client.WithToken(token)
	return a.Client.AsJson().PutUnmarshal(context.Background(), path, data, d)
}
