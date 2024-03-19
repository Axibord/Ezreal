package middlewares

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/labstack/echo/v4"
)

type htmxNonceKey string
type responseTargetsNonceKey string
type twNonceKey string

const (
	HtmxNonceKey            = htmxNonceKey("htmxNonce")
	ResponseTargetsNonceKey = responseTargetsNonceKey("responseTargetsNonce")
	TwNonceKey              = twNonceKey("twNonce")
)

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

func CSPMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			htmxNonce := generateRandomString(16)
			responseTargetsNonce := generateRandomString(16)
			twNonce := generateRandomString(16)

			oldContext := c.Request().Context()
			newContext := context.WithValue(oldContext, HtmxNonceKey, htmxNonce)
			newContext = context.WithValue(newContext, ResponseTargetsNonceKey, responseTargetsNonce)
			newContext = context.WithValue(newContext, TwNonceKey, twNonce)

			// update the context with the new values
			c.SetRequest(c.Request().WithContext(newContext))

			// The hash of the CSS that HTMX injects
			htmxCSSHash := "sha256-pgn1TCGZX6O77zDvy0oTODMOxemn0oj0LeCnQTRj7Kg="

			cspHeader := fmt.Sprintf("default-src 'self'; script-src 'nonce-%s' 'nonce-%s'; style-src 'nonce-%s' '%s'; font-src 'self' https://fonts.googleapis.com https://fonts.gstatic.com; script-src-elem 'self';",
				htmxNonce, responseTargetsNonce, twNonce, htmxCSSHash)
			c.Response().Header().Set("Content-Security-Policy", cspHeader)

			// Proceed with the next middleware/handler in the chain
			return next(c)
		}
	}
}
