# picket-go

The official Picket API Go SDK

## Installation

```bash 
go get -u github.com/picketapi/picket-go
```

## Usage - Quick Start

Use the `picket.NewClient` to create the Picket API client. It takes a _secret API key_ as a parameter.

```go
import (
	picket "github.com/picketapi/picket-go"
)

client := picket.NewClient("YOUR_SECRET_API_KEY")
```

## Nonce

A `nonce` is random value generated by the Picket API to that user must sign to prove ownership a wallet address. The `nonce` function can be used to implement your own wallet authentication flow. 

A nonce is unique to a project and wallet address. If a `nonce` doesn't exist for the project and wallet address, Picket will generate a new nonce; otherwise, Picket will return the existing nonce. A nonce is valid for two minutes before self-destructing.

```go
resp, err := picket.Nonce(picket.NonceArgs{
		Chain:         "solana",
		WalletAddress: "wAllEtadDresS",
})
fmt.Printf("received nonce: %s", resp.Nonce)
```

## Auth

`auth` is the server-side equivalent of login. `auth` should only be used in a trusted server environment. The most common use-case for `auth` is [linking a wallet to an existing application account](https://docs.picketapi.com/picket-docs/tutorials/link-a-wallet-to-a-web-2.0-account).

```go
resp, err := picket.Auth(picket.AuthArgs{
		Chain:         "ethereum",
		WalletAddress: "0x1234567890",
		Signature:     "abcdefghijklmnop",
})
fmt.Printf("received access token: %s", resp.AccessToken)
```

## Authz (Authorize)
`authz` stands for authorization. Unlike Auth, which handles both authentication and authorization, Authz only handles authorization. 
Given an authenticated user's access token and authorization requirements, `authz` will issue a new access token on success (user is authorized) or, on failure, it will return a 4xx HTTP error code.
```go
resp, err := picket.Authz(picket.AuthzArgs{
		AccessToken: "aaa.bbb.ccc",
		Requirements: picket.AuthorizationRequirements{
			ContractAddress: "0xContract",
			MinTokenBalance: "100",
		},
})
fmt.Printf("received updated access token: %s", resp.AccessToken)
```

## Validate
`validate` validates an access token. `validate` should be called, or manually access token validation should be done, server-side before trusting a request's access token. It's common to move access token validation and decoding logic to a shared middleware across API endpoints.
If the access token is valid, validate returns the decoded claims of the access token.

```go
resp, err := picket.Validate(picket.ValidateArgs{
		AccessToken: "aaa.bbb.ccc",
		// Authorizatopn args are optional
		Requirements: picket.AuthorizationRequirements{
			ContractAddress: "0xContract",
			MinTokenBalance: "100",
		},
})
fmt.Printf("decoded access token: %v", resp)
```

## Verify Token Ownership
If you only want to verify token ownership server side for a given wallet, `tokenOwnership` allows you to do just that.

```go
resp, err := picket.TokenOwnership(picket.TokenOwnershipArgs{
		Chain:         "solana",
		WalletAddress: "waLLeTaddrESs",
		Requirements: picket.AuthorizationRequirements{
			Collection: "METAPLEX_COLLECTION",
			MinTokenBalance: "3",
		},
})
fmt.Printf("wallet has sufficient tokens: %s", resp.Allowed)
```
