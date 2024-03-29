// Package picketapi is a Go client for the Picket API
//
// The package is a simple wrapper around the Picket API for
// Go applications
package picketapi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
)

// TODO: Code generate this from the Open API spec

type PicketClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

const (
	APIHost = "picketapi.com"
	// APIVersion is the version for the Picket API used by this client
	APIVersion = "v1"
	// APIBasePath is the base path for the Picket API used by this client
	APIBasePath = "/api/" + APIVersion
	APIBaseURL  = "https://" + APIHost + APIBasePath
)

// NewClient creates a new PicketClient
func NewClient(apiKey string) *PicketClient {
	return &PicketClient{
		apiKey:     apiKey,
		baseURL:    APIBaseURL,
		httpClient: &http.Client{},
	}
}

// SetHTTPClient is the http client to use for requests.
func (p *PicketClient) SetHTTPClient(client *http.Client) {
	p.httpClient = client
}

// SetBaseURL sets the base URL of the Picket api.
// Typically this is used for testing non-production deployments of the Picket API.
func (p *PicketClient) SetBaseURL(baseURL string) {
	p.baseURL = baseURL
}

// HTTPHeaders returns the default HTTP headers for the Picket client
func (p PicketClient) HTTPHeaders() http.Header {
	return http.Header{
		"User-Agent":   []string{"picket-go"},
		"Content-Type": []string{"application/json"},
	}
}

func (p PicketClient) DoRequest(method, apiPath string, body interface{}) (*http.Response, error) {
	reqURL, err := url.JoinPath(p.baseURL, apiPath)
	if err != nil {
		return nil, err
	}

	// set body if provided
	var reqBody io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, reqURL, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header = p.HTTPHeaders()
	req.SetBasicAuth(p.apiKey, "")

	return p.httpClient.Do(req)
}

// decodeResponse decodes the response body into an ErrorResponse or the given interface
func decodeResponse[T interface{}](resp *http.Response, body T) (T, error) {
	if !is2xxStatusCode(resp.StatusCode) {
		var errResp ErrorResponse
		if decodeErr := json.NewDecoder(resp.Body).Decode(&errResp); decodeErr != nil {
			return body, decodeErr
		}
		return body, errResp
	}

	err := json.NewDecoder(resp.Body).Decode(&body)
	return body, err
}

// Nonce returns a nonce for given user's wallet address
func (p PicketClient) Nonce(args NonceArgs) (NonceResponse, error) {
	resp, err := p.DoRequest("POST", "/auth/nonce", args)
	var body NonceResponse

	if err != nil {
		return body, err
	}

	return decodeResponse(resp, body)
}

// Auth authenticates and authorizes a user.
// It returns a access token that can be for the rest of the user's session
func (p PicketClient) Auth(args AuthArgs) (AuthResponse, error) {
	resp, err := p.DoRequest("POST", "/auth", args)
	var body AuthResponse

	if err != nil {
		return body, err
	}

	return decodeResponse(resp, body)
}

// Authz authorizes an authenticated user's access token for the given requirements.
// On success, it returns an updated access token
func (p PicketClient) Authz(args AuthzArgs) (AuthResponse, error) {
	resp, err := p.DoRequest("POST", "/authz", args)
	var body AuthResponse

	if err != nil {
		return body, err
	}

	return decodeResponse(resp, body)
}

// Validate validates an access token, optionally, for the given requirements.
// If the access token is valid, then the decoded user payload is returned.
func (p PicketClient) Validate(args ValidateArgs) (AuthorizedUser, error) {
	resp, err := p.DoRequest("POST", "/auth/validate", args)
	var body AuthorizedUser

	if err != nil {
		return body, err
	}

	return decodeResponse(resp, body)
}

// TokenOwnership checks if the wallet address has a given token balance.
// Similar to Authz, but does not require nor return an access token
func (p PicketClient) TokenOwnership(args TokenOwnershipArgs) (TokenOwnershipResponse, error) {
	apiPath := path.Join("chains", args.Chain, "wallets", args.WalletAddress, "tokenOwnership")
	resp, err := p.DoRequest("POST", apiPath, args.Requirements)
	var body TokenOwnershipResponse

	if err != nil {
		return body, err
	}

	return decodeResponse(resp, body)
}
