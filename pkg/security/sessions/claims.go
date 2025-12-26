package sessions

import "context"

type ctxKeyAuthClaims int

const PrincipalClaimsKey ctxKeyAuthClaims = 1

// PrincipalClaims represents the claims of a principal, customize it as needed
type PrincipalClaims struct {
	IsAuthenticated bool     `json:"isAuthenticated"`
	PrincipalID     string   `json:"principalId"`
	OrgID           string   `json:"orgId"`
	Roles           []string `json:"roles"`
}

func (c *PrincipalClaims) HasRole(theRole string) bool {
	for _, role := range c.Roles {
		if role == theRole {
			return true
		}
	}
	return false
}

func AddToContext(ctx context.Context, claims PrincipalClaims) context.Context {
	return context.WithValue(ctx, PrincipalClaimsKey, claims)
}

func ExtractClaimsFromContext(ctx context.Context) PrincipalClaims {
	claims, ok := ctx.Value(PrincipalClaimsKey).(PrincipalClaims)
	if !ok {
		return PrincipalClaims{IsAuthenticated: false}
	}

	return claims
}
