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
	BaseURL               string
	AdAccountID           string
	AdAccountIDs          []string // v1.1.0: For AGENCY role with multiple ad accounts
	AdAccountTitle        string
	AdManagerAccountID    string // v1.1.0: For Ad Manager Account roles
	AdManagerAccountTitle string // v1.1.0: For Ad Manager Account title
	Email                 string
	ExternalUserID        string
	Name                  string
	Path                  string
	PlatformID            string
	Role                  string
	Secret                string
	ColorMode             string
	Language              string
	Version               string
}

func (p *SSOParams) createSignedRmpPortalURL() string {
	var (
		nonce     = generateNonce(16)
		timestamp = strconv.FormatInt((time.Now()).Unix(), 10)
		signature string
	)

	// Generate signature based on version
	if p.Version == "1.1.0" {
		signature = p.generateSignatureV110(nonce, timestamp)
	} else {
		signature = p.generateSignatureV100(nonce, timestamp)
	}

	// Build query parameters
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

	// Add v1.1.0 fields if present
	if p.Version == "1.1.0" {
		if len(p.AdAccountIDs) > 0 {
			for _, id := range p.AdAccountIDs {
				queryParams.Add("ad_account_ids", id)
			}
		}
		if p.AdManagerAccountID != "" {
			queryParams.Add("ad_manager_account_id", p.AdManagerAccountID)
		}
		if p.AdManagerAccountTitle != "" {
			queryParams.Add("ad_manager_account_title", p.AdManagerAccountTitle)
		}
	}

	// Add optional config parameters
	queryParams.Add("config:color_mode", p.ColorMode)
	queryParams.Add("config:language", p.Language)

	return p.BaseURL + "/sso?" + queryParams.Encode()
}

// generateSignatureV100 creates the signature for v1.0.0
func (p *SSOParams) generateSignatureV100(nonce, timestamp string) string {
	params := []string{
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
	concatenatedString := strings.Join(params, "\n")

	hash := hmac.New(sha256.New, []byte(p.Secret))
	hash.Write([]byte(concatenatedString))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

// generateSignatureV110 creates the signature for v1.1.0
// Parameters are in alphabetical order with new fields included
func (p *SSOParams) generateSignatureV110(nonce, timestamp string) string {
	// Join ad_account_ids array with semicolon separator if present
	adAccountIDs := ""
	if len(p.AdAccountIDs) > 0 {
		adAccountIDs = strings.Join(p.AdAccountIDs, ";")
	}

	// Parameters in alphabetical order
	params := []string{
		p.AdAccountID,           // ad_account_id
		adAccountIDs,            // ad_account_ids
		p.AdAccountTitle,        // ad_account_title
		p.AdManagerAccountID,    // ad_manager_account_id
		p.AdManagerAccountTitle, // ad_manager_account_title
		p.Email,                 // email
		p.ExternalUserID,        // external_user_id
		p.Name,                  // name
		nonce,                   // nonce
		p.Path,                  // path
		p.PlatformID,            // platform_id
		p.Role,                  // role
		timestamp,               // timestamp
		p.Version,               // version
	}
	concatenatedString := strings.Join(params, "\n")

	hash := hmac.New(sha256.New, []byte(p.Secret))
	hash.Write([]byte(concatenatedString))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

// exampleAdAccount demonstrates creating a signed URL for a single ad account.
// Works with AD_ACCOUNT_OWNER, AD_ACCOUNT_USER, or AD_ACCOUNT_VIEWER roles.
func exampleAdAccount() string {
	const (
		baseURL        = "https://{YOUR-RMP-PORTAL_URL}" // Please use the url provided by your account manager
		adAccountID    = "my-ad-account-id"
		adAccountTitle = "My Ad Account"
		email          = "test@example.com"
		externalUserID = "user-id"
		name           = "Example User Name"
		platformID     = "RMP_PLATFORM_ID"
		role           = "AD_ACCOUNT_OWNER"
		secret         = "super-secret"
		colorMode      = "light" // light or dark
		language       = "en"    // en or ko
		version        = "1.1.0"
	)

	params := &SSOParams{
		BaseURL:        baseURL,
		AdAccountID:    adAccountID,
		AdAccountTitle: adAccountTitle,
		Email:          email,
		ExternalUserID: externalUserID,
		Name:           name,
		Path:           "/embed/sponsored-ads/cm/a/" + adAccountID,
		PlatformID:     platformID,
		Role:           role,
		Secret:         secret,
		ColorMode:      colorMode,
		Language:       language,
		Version:        version,
	}

	return params.createSignedRmpPortalURL()
}

// exampleAdAccountAgency demonstrates creating a signed URL for multiple ad accounts with AD_ACCOUNT_AGENCY role (v1.1.0 feature).
// This allows an agency user to access multiple ad accounts.
func exampleAdAccountAgency() string {
	const (
		baseURL        = "https://{YOUR-RMP-PORTAL_URL}"
		email          = "agency@example.com"
		externalUserID = "agency-user-id"
		name           = "Agency User"
		platformID     = "RMP_PLATFORM_ID"
		role           = "AD_ACCOUNT_AGENCY"
		secret         = "super-secret"
		colorMode      = "light"
		language       = "en"
		version        = "1.1.0"
	)

	params := &SSOParams{
		BaseURL:        baseURL,
		AdAccountIDs:   []string{"account-1", "account-2", "account-3"},
		Email:          email,
		ExternalUserID: externalUserID,
		Name:           name,
		Path:           "/embed/sponsored-ads/cm",
		PlatformID:     platformID,
		Role:           role,
		Secret:         secret,
		ColorMode:      colorMode,
		Language:       language,
		Version:        version,
	}

	return params.createSignedRmpPortalURL()
}

// exampleAdManagerAccount demonstrates creating a signed URL for an ad manager account (v1.1.0 feature).
// Works with AD_MANAGER_ACCOUNT_OWNER or AD_MANAGER_ACCOUNT_USER roles.
func exampleAdManagerAccount() string {
	const (
		baseURL               = "https://{YOUR-RMP-PORTAL_URL}"
		adManagerAccountID    = "my-ad-manager-account-id"
		adManagerAccountTitle = "My Ad Manager Account"
		email                 = "manager@example.com"
		externalUserID        = "manager-user-id"
		name                  = "Manager User"
		platformID            = "RMP_PLATFORM_ID"
		role                  = "AD_MANAGER_ACCOUNT_OWNER"
		secret                = "super-secret"
		colorMode             = "light"
		language              = "en"
		version               = "1.1.0"
	)

	params := &SSOParams{
		BaseURL:               baseURL,
		AdManagerAccountID:    adManagerAccountID,
		AdManagerAccountTitle: adManagerAccountTitle,
		Email:                 email,
		ExternalUserID:        externalUserID,
		Name:                  name,
		Path:                  "/embed/sponsored-ads/cm/ama/" + adManagerAccountID,
		PlatformID:            platformID,
		Role:                  role,
		Secret:                secret,
		ColorMode:             colorMode,
		Language:              language,
		Version:               version,
	}

	return params.createSignedRmpPortalURL()
}

func main() {
	// Choose which example to run by uncommenting the desired line:
	signedURL := exampleAdAccount()
	// signedURL := exampleAdAccountAgency()
	// signedURL := exampleAdManagerAccount()

	fmt.Println(signedURL)
}
