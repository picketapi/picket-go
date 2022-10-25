package picket_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	picket "github.com/picketapi/picket-go"
)

const apiKey = "YOUR_API_KEY"

func TestPicketDoRequest(t *testing.T) {
	method := "GET"
	path := "/test"
	body := picket.NonceArgs{
		Chain:         "ethereum",
		WalletAddress: "0x1234567890",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request method and path
		if r.Method != method {
			t.Errorf("Expected %s, got %s", method, r.Method)
		}
		if r.URL.Path != path {
			t.Errorf("Expected %s, got %s", path, r.URL.Path)
		}
		// check auth header
		username, _, _ := r.BasicAuth()
		if username != apiKey {
			t.Errorf("Expected api key in basic auth header %s, got %s", apiKey, username)
		}
		// check json header
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected application/json, got %s", contentType)
		}

		reqBody := picket.NonceArgs{}
		json.NewDecoder(r.Body).Decode(&reqBody)

		if reqBody.Chain != body.Chain {
			t.Errorf("Expected %s, got %s", body.Chain, reqBody.Chain)
		}
		if reqBody.WalletAddress != body.WalletAddress {
			t.Errorf("Expected %s, got %s", body.WalletAddress, reqBody.WalletAddress)
		}

		// check body
		w.Write([]byte("success"))
	}))
	defer ts.Close()

	client := picket.NewClient(apiKey)
	client.SetBaseURL(ts.URL)
	client.SetHTTPClient(ts.Client())

	_, err := client.DoRequest(method, path, body)

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestPicketNonce(t *testing.T) {
	want := picket.NonceResponse{
		Nonce:     "abcdefghijklmnop",
		Statement: "Woo",
		Format:    "siwe",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	}))
	defer ts.Close()

	client := picket.NewClient(apiKey)
	client.SetBaseURL(ts.URL)
	client.SetHTTPClient(ts.Client())

	args := picket.NonceArgs{
		Chain:         "ethereum",
		WalletAddress: "0x1234567890",
	}
	got, err := client.Nonce(args)

	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	if got.Nonce != want.Nonce {
		t.Errorf("got %s, want %s", got.Nonce, want.Nonce)
	}
	if got.Statement != want.Statement {
		t.Errorf("got %s, want %s", got.Statement, want.Statement)
	}
	if got.Format != want.Format {
		t.Errorf("got %s, want %s", got.Format, want.Format)
	}
}

func TestPicketAuth(t *testing.T) {
	want := picket.AuthResponse{
		User: picket.AuthorizedUser{
			Chain:          "ethereum",
			WalletAddress:  "0x1234567890",
			DisplayAddress: "my.name.eth",
		},
		AccessToken: "xxx.yyy.zzz",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	}))
	defer ts.Close()

	client := picket.NewClient(apiKey)
	client.SetBaseURL(ts.URL)
	client.SetHTTPClient(ts.Client())

	args := picket.AuthArgs{
		Chain:         "ethereum",
		WalletAddress: "0x1234567890",
		Signature:     "abcdefghijklmnop",
	}
	got, err := client.Auth(args)

	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	if got.AccessToken != want.AccessToken {
		t.Errorf("got %s, want %s", got.AccessToken, want.AccessToken)
	}
	if got.User.DisplayAddress != want.User.DisplayAddress {
		t.Errorf("got %s, want %s", got.User.DisplayAddress, want.User.DisplayAddress)
	}
}

func TestPicketValidate(t *testing.T) {
	want := picket.AuthorizedUser{
		Chain:          "ethereum",
		WalletAddress:  "0x1234567890",
		DisplayAddress: "my.name.eth",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	}))
	defer ts.Close()

	client := picket.NewClient(apiKey)
	client.SetBaseURL(ts.URL)
	client.SetHTTPClient(ts.Client())

	args := picket.ValidateArgs{
		AccessToken: "xxx.yyy.zzz",
	}
	got, err := client.Validate(args)

	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	if got.Chain != want.Chain {
		t.Errorf("got %s, want %s", got.Chain, want.Chain)
	}
	if got.WalletAddress != want.WalletAddress {
		t.Errorf("got %s, want %s", got.WalletAddress, want.WalletAddress)
	}
	if got.DisplayAddress != want.DisplayAddress {
		t.Errorf("got %s, want %s", got.DisplayAddress, want.DisplayAddress)
	}
}

func TestPicketAuthorize(t *testing.T) {
	want := picket.AuthResponse{
		User: picket.AuthorizedUser{
			Chain:          "ethereum",
			WalletAddress:  "0x1234567890",
			DisplayAddress: "my.name.eth",
		},
		AccessToken: "xxx.yyy.zzz",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	}))
	defer ts.Close()

	client := picket.NewClient(apiKey)
	client.SetBaseURL(ts.URL)
	client.SetHTTPClient(ts.Client())

	args := picket.AuthzArgs{
		AccessToken: "aaa.bbb.ccc",
		Requirements: picket.AuthorizationRequirements{
			ContractAddress: "0xContract",
			MinTokenBalance: "100",
		},
	}
	got, err := client.Authz(args)

	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	if got.AccessToken != want.AccessToken {
		t.Errorf("got %s, want %s", got.AccessToken, want.AccessToken)
	}
}

func TestPicketTokenOwnership(t *testing.T) {
	want := picket.TokenOwnershipResponse{
		Allowed: true,
		TokenBalances: map[string]string{
			"0x1234567890": "100",
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(want)
	}))
	defer ts.Close()

	client := picket.NewClient(apiKey)
	client.SetBaseURL(ts.URL)
	client.SetHTTPClient(ts.Client())

	args := picket.TokenOwnershipArgs{
		Chain:         "solana",
		WalletAddress: "0x1234567890",
		Requirements: picket.AuthorizationRequirements{
			ContractAddress: "0xContract",
			MinTokenBalance: "100",
		},
	}
	got, err := client.TokenOwnership(args)

	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	if got.Allowed != want.Allowed {
		t.Errorf("got %t, want %t", got.Allowed, want.Allowed)
	}

	for k, v := range got.TokenBalances {
		if v != want.TokenBalances[k] {
			t.Errorf("got %s, want %s", v, want.Balances[k])
		}
	}
}

func TestPicketErrorResponse(t *testing.T) {
	want := picket.ErrorResponse{
		Msg:  "Oh no! An error",
		Code: "testing_error",
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(want)
	}))
	defer ts.Close()

	client := picket.NewClient(apiKey)
	client.SetBaseURL(ts.URL)
	client.SetHTTPClient(ts.Client())

	args := picket.NonceArgs{
		Chain:         "ethereum",
		WalletAddress: "0x1234567890",
	}
	resp, err := client.Nonce(args)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if resp.Nonce != "" {
		t.Errorf("Expected empty response, got nonce: %s", resp.Nonce)
	}
	if resp.Statement != "" {
		t.Errorf("Expected empty response, got statement: %s", resp.Statement)
	}
	if resp.Format != "" {
		t.Errorf("Expected empty response, got format: %s", resp.Format)
	}

	var got picket.ErrorResponse
	if !errors.As(err, &got) {
		t.Errorf("Expected ErrorRespose, got: %v", err)
	}

	if got.Msg != want.Msg {
		t.Errorf("got %s, want %s", got.Msg, want.Msg)
	}
	if got.Code != want.Code {
		t.Errorf("got %s, want %s", got.Code, want.Code)
	}
}
