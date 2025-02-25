// backend/src/claims_handler.go
package handlers

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
)

func getUintFromClaim(claims jwt.MapClaims, key string) (uint, error) {
	val, exists := claims[key]
	if !exists {
		return 0, fmt.Errorf("claim %s not found", key)
	}
	switch v := val.(type) {
	case float64:
		return uint(v), nil
	case int:
		return uint(v), nil
	case int64:
		return uint(v), nil
	default:
		return 0, fmt.Errorf("invalid type for claim %s", key)
	}
}


