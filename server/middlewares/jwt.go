package middlewares

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"to-persist/server/global"
)

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token: ")
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")

		if token == "" {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		token = strings.Fields(token)[1]

		standardClaims, err := ParseToken(token)
		if err != nil {
			// If we are not just returning a status code,
			// we need to respond with different messages depending on the type of the error
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}
		// Get user id from subject and set the value to gin's context
		c.Set("user-id", standardClaims.Subject)
		c.Next()
	}
}

func ParseToken(tokenStr string) (*jwt.StandardClaims, error) {
	// Parse the JWT token with the standard claims structure.
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(global.Config.JwtConfig.JwtKey), nil
	})

	if err != nil {
		return nil, evaluateTokenError(err)
	}
	// Assert that the token's claims are of type *jwt.StandardClaims and the token is valid.
	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, TokenInvalid
}

func evaluateTokenError(err error) error {
	var ve *jwt.ValidationError
	if errors.As(err, &ve) {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return TokenMalformed
		}
		if ve.Errors&jwt.ValidationErrorExpired != 0 {
			return TokenExpired
		}
		if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
			return TokenNotValidYet
		}
		// Add additional JWT error checks if necessary.
	}
	return TokenInvalid
}
