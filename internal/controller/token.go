package controller

import (
	"echo-sqlc-template/internal/database"
	"echo-sqlc-template/internal/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type TokenController struct{}

func NewTokenController(e *echo.Echo) {
	c := TokenController{}

	group := e.Group("/tokens")
	group.POST("/refresh", c.Refresh)
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func (TokenController) Refresh(c echo.Context) error {
	body, err := Body[RefreshTokenRequest](c)
	if err != nil {
		return err
	}

	token, err := database.Q.GetToken(ctx, database.GetTokenParams{RefreshToken: body.RefreshToken})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Token")
	}

	if token.CreatedAt.Time.Before(time.Now().Add(-720 * time.Hour)) {
		return echo.NewHTTPError(http.StatusUnauthorized, "Token has fully expired, re-login")
	}

	token, err = database.Q.UpdateToken(ctx, database.UpdateTokenParams{
		AccessTokenDoUpdate:  true,
		RefreshTokenDoUpdate: true,
		ExpiredAtDoUpdate:    true,
		AccessToken:          utils.GenerateRandomString(48),
		RefreshToken:         utils.GenerateRandomString(64),
		ExpiredAt: pgtype.Timestamp{
			Time:  time.Now().Local().Add(time.Hour * time.Duration(1)),
			Valid: true,
		},
		ID: token.ID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to create tokens")
	}

	return c.JSON(200, token)
}
