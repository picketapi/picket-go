// Package picketapi is a Go client for the Picket API
//
// The package is a simple wrapper around the Picket API for
// Go applications
package picketapi

type NonceArgs struct {
	Chain         string `json:"chain"`
	WalletAddress string `json:"wallet_address"`
	Locale        string `json:"locale"`
}

type SigningMessageFormat string

const (
	SigningMessageFormatSimple SigningMessageFormat = "simple"
	SigningMessageFormatSIWE   SigningMessageFormat = "siwe"
)

type NonceResponse struct {
	Nonce     string               `json:"nonce"`
	Statement string               `json:"statement"`
	Format    SigningMessageFormat `json:"format"`
}

type AuthorizationRequirements struct {
	ContractAddress string   `json:"contractAddress,omitempty"`
	MinTokenBalance string   `json:"minTokenBalance,omitempty"`
	AllowedWallets  []string `json:"allowedWallets,omitempty"`
	// Solana specific requirements
	Collection     string   `json:"collection,omitempty"`
	CreatorAddress string   `json:"creatorAddress,omitempty"`
	TokenIds       []string `json:"token_ids,omitempty"`
}

type ChainType = string

const (
	ChainTypeEthereum ChainType = "ethereum"
	ChainTypeSolana   ChainType = "solana"
	ChainTypeFlow     ChainType = "flow"
)

type SigningMessageContext struct {
	Domain    string    `json:"domain"`
	URI       string    `json:"uri"`
	ChainID   string    `json:"chainId"`
	IssuedAt  string    `json:"issuedAt"`
	ChainType ChainType `json:"chainType"`
	Locale    string    `json:"locale"`
}

type AuthArgs struct {
	Chain         string                    `json:"chain,omitempty"`
	WalletAddress string                    `json:"walletAddress"`
	Signature     string                    `json:"signature"`
	Requirements  AuthorizationRequirements `json:"requirements,omitempty"`
	Context       SigningMessageContext     `json:"context,omitempty"`
}

type TokenBalances map[string]string

type AuthorizedUser struct {
	Chain          string        `json:"chain"`
	WalletAddress  string        `json:"walletAddress"`
	DisplayAddress string        `json:"displayAddress"`
	TokenBalances  TokenBalances `json:"tokenBalances"`
}

type AuthResponse struct {
	User        AuthorizedUser `json:"user"`
	AccessToken string         `json:"accessToken"`
}

type AuthzArgs struct {
	AccessToken  string                    `json:"accessToken"`
	Requirements AuthorizationRequirements `json:"requirements,omitempty"`
	Revalidate   bool                      `json:"revalidate"`
}

type ValidateArgs struct {
	AccessToken  string                    `json:"accessToken"`
	Requirements AuthorizationRequirements `json:"requirements,omitempty"`
}

type TokenOwnershipArgs struct {
	Chain         string                    `json:"chain"`
	WalletAddress string                    `json:"walletAddress"`
	Requirements  AuthorizationRequirements `json:"requirements,omitempty"`
}

type TokenOwnershipResponse struct {
	Allowed       bool          `json:"allowed"`
	TokenBalances TokenBalances `json:"tokenBalances"`
}

type Picket interface {
	// Nonce returns a nonce for given user's wallet address
	Nonce(args NonceArgs) (NonceResponse, error)
	// Auth authenticates and authorizes a user.
	// It returns a access token that can be for the rest of the user's session
	Auth(args AuthArgs) (AuthResponse, error)
	// Authz authorizes an authenticated user's access token for the given requirements.
	// On success, it returns an updated access token
	Authz(args AuthzArgs) (AuthResponse, error)
	// Validate validates an access token, optionally, for the given requirements.
	// If the access token is valid, then the decoded user payload is returned.
	Validate(args ValidateArgs) (AuthorizedUser, error)
	// TokenOwnership checks if the wallet address has a given token balance.
	// Similar to Authz, but does not require nor return an access token
	TokenOwnership(args TokenOwnershipArgs) (TokenOwnershipResponse, error)
}
