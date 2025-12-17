package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Using 16 characters long random nonce would be sufficient in most use cases
func generateNonce(length int) string {
	const possibleCharacters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = possibleCharacters[seededRand.Intn(len(possibleCharacters))]
	}
	return string(b)
}

// SSOParams represent the input params that are required to create a signed embedded campaign manager URL.
type SSOParams struct {
	BaseURL        string
	AdAccountID    string
	AdAccountTitle string
	Email          string
	ExternalUserID string
	Name           string
	Path           string
	PlatformID     string
	Role           string
	Secret         string
	ColorMode      string
	Language       string
	Version        string
}

func (p *SSOParams) createSignedRmpPortalURL() string {
	var (
		nonce     = generateNonce(16)
		timestamp = strconv.FormatInt((time.Now()).Unix(), 10)
		params    = []string{
			p.AdAccountID,
			p.AdAccountTitle,
			p.Email,
			p.ExternalUserID,
			p.Name,
			nonce,
			p.Path,
			p.PlatformID,
			p.Role,
			timestamp,
			p.Version,
		}
		concatenatedString = strings.Join(params, "\n")
		signature          string
	)

	signature = func() string {
		hash := hmac.New(sha256.New, []byte(p.Secret))
		hash.Write([]byte(concatenatedString))

		return base64.StdEncoding.EncodeToString(hash.Sum(nil))
	}()

	queryParams := url.Values{}
	queryParams.Add("ad_account_id", p.AdAccountID)
	queryParams.Add("ad_account_title", p.AdAccountTitle)
	queryParams.Add("email", p.Email)
	queryParams.Add("external_user_id", p.ExternalUserID)
	queryParams.Add("name", p.Name)
	queryParams.Add("nonce", nonce)
	queryParams.Add("path", p.Path)
	queryParams.Add("platform_id", p.PlatformID)
	queryParams.Add("role", p.Role)
	queryParams.Add("timestamp", timestamp)
	queryParams.Add("version", p.Version)
	queryParams.Add("signature", signature)
	queryParams.Add("config:color_mode", p.ColorMode)
	queryParams.Add("config:language", p.Language)

	return p.BaseURL + "/sso?" + queryParams.Encode()
}

func main() {
	const (
		baseURL        = "https://{YOUR-RMP-PORTAL_URL}" // Please use the url provided by your account manager
		adAccountID    = "my-ad-account-id"
		adAccountTitle = "My Ad Account"
		email          = "test@example.com"
		externalUserID = "user-id"
		name           = "Example User Name"
		path           = "/embed/sponsored-ads/cm/a/" + adAccountID
		platformID     = "RMP_PLATFORM_ID"
		role           = "AD_ACCOUNT_OWNER"
		secret         = "super-secret"
		colorMode      = "light" // light or dark
		language       = "en"    // en or ko
		version        = "1.0.0"
	)

	params := &SSOParams{
		BaseURL:        baseURL,
		AdAccountID:    adAccountID,
		AdAccountTitle: adAccountTitle,
		Email:          email,
		ExternalUserID: externalUserID,
		Name:           name,
		Path:           path,
		PlatformID:     platformID,
		Role:           role,
		Secret:         secret,
		ColorMode:      colorMode,
		Language:       language,
		Version:        version,
	}

	signedURL := params.createSignedRmpPortalURL()

	fmt.Println(signedURL)
}
