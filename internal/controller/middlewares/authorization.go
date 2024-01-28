package middlewares

import (
	"context"
	"echo-sqlc-template/internal/database"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"time"
)

type AuthorizationConfig struct {
	Optional bool
	Role     database.UserRole
}

type OptsFn func(*AuthorizationConfig)

func defaultOpts() AuthorizationConfig {
	return AuthorizationConfig{Role: "USER", Optional: false}
}

func WithRole(role database.UserRole) OptsFn {
	return func(c *AuthorizationConfig) {
		c.Role = role
	}
}

func SetOptional() OptsFn {
	return func(c *AuthorizationConfig) {
		c.Optional = true
	}
}

var roleRankTable = map[database.UserRole]int{
	database.UserRoleUSER:          0,
	database.UserRoleMODERATOR:     1,
	database.UserRoleADMINISTRATOR: 2,
}

func roleRank(r *database.UserRole) int {
	rank, _ := roleRankTable[*r]
	return rank
}

func Authorization(opts ...OptsFn) echo.MiddlewareFunc {
	cfg := defaultOpts()
	for _, fn := range opts {
		fn(&cfg)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("user", nil)

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				if cfg.Optional {
					return next(c)
				}

				return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization Header")
			}

			headersParts := strings.Split(authHeader, " ")
			if len(headersParts) != 2 {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}

			token, err := database.Q.GetToken(context.Background(), database.GetTokenParams{AccessToken: headersParts[1]})
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Token")

			}

			if token.ExpiredAt.Time.Before(time.Now().Add(time.Duration(-1) * time.Hour)) {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token has expired")
			}

			user, err := database.Q.GetUserByTokenID(context.Background(), token.ID)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Token")
			}

			c.Set("user", user)

			if roleRank(&user.Role) < roleRank(&cfg.Role) {
				return echo.ErrForbidden
			}

			return next(c)
		}
	}
}
