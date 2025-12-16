package sessions

import "context"

type ctxKeyAuthClaims int

const PrincipalClaimsKey ctxKeyAuthClaims = 1

type AuthClaimsMetadata struct {
	OrganizationID string `json:"organizationId"`
}

type AuthClaims struct {
	IsAuthenticated bool               `json:"isAuthenticated"`
	Principal       string             `json:"principal"`
	Roles           []string           `json:"roles"`
	Metadata        AuthClaimsMetadata `json:"metadata"`
}

func (c AuthClaims) HasRole(theRole string) bool {
	for _, role := range c.Roles {
		if role == theRole {
			return true
		}
	}
	return false
}

func AddAuthClaimsToContext(ctx context.Context, claims *AuthClaims) context.Context {
	return context.WithValue(ctx, PrincipalClaimsKey, claims)
}

func GetAuthClaimsFromContext(ctx context.Context) *AuthClaims {
	claims, ok := ctx.Value(PrincipalClaimsKey).(*AuthClaims)
	if !ok {
		return &AuthClaims{IsAuthenticated: false}
	}

	return claims
}
