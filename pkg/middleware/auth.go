package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"golang-fiber/pkg/auth"
	"golang-fiber/pkg/common"
	"golang-fiber/pkg/constants"
)

func JWTProtected(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return common.ResponseJsonWithCode(c, fiber.StatusUnauthorized, uuid.New(),
				constants.CodeUnauthorized, constants.MessageENUnauthorized, constants.MessageTHUnauthorized, nil)
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := auth.ValidateToken(tokenStr, jwtSecret)
		if err != nil {
			return common.ResponseJsonWithCode(c, fiber.StatusUnauthorized, uuid.New(),
				constants.CodeUnauthorized, constants.MessageENUnauthorized, constants.MessageTHUnauthorized, nil)
		}

		ctx := auth.SetUserID(c.UserContext(), claims.UserID)
		c.SetUserContext(ctx)
		return c.Next()
	}
}
