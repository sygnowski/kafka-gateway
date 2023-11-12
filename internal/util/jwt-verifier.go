package util

import (
	"log"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type JwtCfg struct {
	Alg          string
	Key          string
	RequiredRole []string `json:"roles"`
}

type Claim struct {
	Roles string
	jwt.RegisteredClaims
}

type JwtVerifier struct {
	cfg JwtCfg
}

func MakeJwtVerifier(config *JwtCfg) *JwtVerifier {
	v := JwtVerifier{cfg: *config}
	return &v
}

func (verifier *JwtVerifier) IsValid(token string) bool {
	claims := &Claim{}

	jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return verifier.cfg.Key, nil
	})

	if err != nil {
		log.Printf("jwt chk error: %v", err)
	}

	if jwtToken != nil && jwtToken.Valid {
		roles := strings.Split(claims.Roles, ",")

		return verifier.allRolesAreValid(roles)
	}

	return false
}

func (verifier *JwtVerifier) allRolesAreValid(roles []string) bool {
	chk := func(want string, from []string) bool {
		has := false
		for _, role := range from {
			if role != want {
				continue
			}
			has = true
			break
		}
		return has
	}

	for _, desired := range verifier.cfg.RequiredRole {
		if !chk(desired, roles) {
			return false
		}
	}

	return true
}
