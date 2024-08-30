package middleware

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/vector-ops/ticketeer/models"
	"gorm.io/gorm"
)

func AuthProtected(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			log.Warn("missing headers")
			return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"status":  "fail",
				"message": "unauthorized",
			})
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			log.Warn("invalid token parts")
			return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"status":  "fail",
				"message": "unauthorized",
			})
		}

		tokenStr := tokenParts[1]
		secret := []byte(os.Getenv("JWT_SECRET"))
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.GetSigningMethod("HS256").Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"status":  "fail",
				"message": "unauthorized",
			})
		}

		userId := token.Claims.(jwt.MapClaims)["id"]

		if err := db.Model(&models.User{}).Where("id = ?", userId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warnf("user not in database")
			return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"status":  "fail",
				"message": "unauthorized",
			})
		}

		c.Locals("userId", userId)

		return c.Next()
	}
}
