package sessions

import (
	"context"
	"strings"
)

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

func (c *PrincipalClaims) ToSession() *Session {
	return &Session{
		PrincipalID: c.PrincipalID,
		Metadata: SessionMetadata{
			"orgId": c.OrgID,
			"roles": strings.Join(c.Roles, ","),
		},
	}
}

func ClaimsFromSession(session *Session) PrincipalClaims {
	if session == nil {
		return PrincipalClaims{IsAuthenticated: false}
	}

	return PrincipalClaims{
		IsAuthenticated: true,
		PrincipalID:     session.PrincipalID,
		OrgID:           session.Metadata["orgId"],
		Roles:           strings.Split(session.Metadata["roles"], ","),
	}
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
