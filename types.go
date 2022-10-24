// Package picket is a Go client for the Picket API
//
// The package is a simple wrapper around the Picket API for
// Go applications
package picket

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
	ContractAddress string `json:"contract_address"`
	MinTokenBalance string `json:"min_token_balance"`
	Collection      string `json:"collection"`
	CreatorAddress  string `json:"creator_address"`
	// TODO: Do we need to support the following?
	TokenIds []string `json:"token_ids"`
}

type AuthArgs struct {
	Chain         string                    `json:"chain"`
	WalletAddress string                    `json:"wallet_address"`
	Signature     string                    `json:"signature"`
	Requirements  AuthorizationRequirements `json:"requirements"`
}

type TokenBalancse map[string]string

type AuthorizedUser struct {
	Chain          string `json:"chain"`
	WalletAddress  string `json:"wallet_address"`
	DisplayAddress string `json:"display_address"`
}

type AuthResponse struct {
	User        AuthorizedUser `json:"user"`
	AccessToken string         `json:"access_token"`
}

type AuthzArgs struct {
	AccessToken  string                    `json:"access_token"`
	Requirements AuthorizationRequirements `json:"requirements"`
}

type TokenOwnershipArgs struct {
	Chain         string                    `json:"chain"`
	WalletAddress string                    `json:"wallet_address"`
	Requirements  AuthorizationRequirements `json:"requirements"`
}

type TokenOwnershipResponse struct {
	Allowed  bool          `json:"allowed"`
	Balances TokenBalancse `json:"balances"`
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
	// TokenOwnership checks if the wallet address has a given token balance.
	// Similar to Authz, but does not require nor return an access token
	TokenOwnership(args TokenOwnershipArgs) (TokenOwnershipResponse, error)
}
