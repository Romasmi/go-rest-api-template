package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/Romasmi/go-rest-api-template/internal/config"
	"github.com/go-chi/jwtauth/v5"
)

var TokenAuth *jwtauth.JWTAuth

func InitAuth(config *config.Config) {
	TokenAuth = jwtauth.New("HS256", []byte(config.JWT.Secret), nil)
}

func Authenticator(next http.Handler) http.Handler {
	return jwtauth.Verifier(TokenAuth)(
		jwtauth.Authenticator(next),
	)
}

func GenerateJWT(userID string, role string) (string, error) {
	expiration := time.Now().Add(24 * time.Hour)

	claims := map[string]interface{}{
		"user_id": userID,
		"role":    role,
		"exp":     expiration.Unix(),
	}

	_, tokenString, err := TokenAuth.Encode(claims)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ExtractBearerToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")

	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}

	return ""
}

func GetUserIDFromToken(claims map[string]interface{}) string {
	if userID, ok := claims["user_id"].(string); ok {
		return userID
	}
	return ""
}

func GetRoleFromToken(claims map[string]interface{}) string {
	if role, ok := claims["role"].(string); ok {
		return role
	}
	return ""
}

func GetClaimsFromRequest(r *http.Request) map[string]interface{} {
	_, claims, _ := jwtauth.FromContext(r.Context())
	return claims
}
