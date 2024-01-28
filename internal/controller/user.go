package controller

import (
	"echo-sqlc-template/internal/controller/middlewares"
	"echo-sqlc-template/internal/database"
	"echo-sqlc-template/internal/services/storage"
	"echo-sqlc-template/internal/utils"
	"fmt"
	"github.com/labstack/echo/v4"
)

type UserController struct{}

func NewUserController(e *echo.Echo) {
	c := UserController{}

	group := e.Group("/users")
	group.GET("/whoami", c.WhoAmI, middlewares.Authorization())
	group.GET("/:id", c.GetUnique)
	group.PATCH("/:id", c.Update)
}

func (UserController) WhoAmI(c echo.Context) error {
	user := c.Get("user").(database.User)
	return c.JSON(200, HideKeys(user.ResolveKey(), "Password"))
}

func (UserController) GetUnique(c echo.Context) error {
	id, err := RouteID(c)
	if err != nil {
		return err
	}

	user, err := database.Q.GetUser(ctx, database.GetUserParams{ID: id})
	if err != nil {
		return err
	}

	return c.JSON(200, HideKeys(user.ResolveKey(), "Email", "Password", "Enabled"))
}

type UpdateUserRequest struct {
	Username *string `json:"username,omitempty" validate:"min=6,max=64"`
	Email    *string `json:"email,omitempty" validate:"email"`
	Avatar   *string `json:"avatar,omitempty" validate:"len=32"`
}

func (UserController) Update(c echo.Context) error {
	body, err := Body[UpdateUserRequest](c)
	if err != nil {
		return err
	}

	id, err := RouteID(c)
	if err != nil {
		return err
	}

	user, err := database.Q.GetUser(ctx, database.GetUserParams{ID: id})
	if err != nil {
		return err
	}

	params := database.UpdateUserParams{
		ID: user.ID,
	}

	if body.Username != nil {
		exists, err := database.Q.UserExistsByUsername(ctx, *body.Username)
		if err != nil || exists {
			return echo.NewHTTPError(400, "Conflict, user with this username exists")
		}

		params.UsernameDoUpdate = true
		params.Username = *body.Username
	}

	if body.Email != nil {
		exists, err := database.Q.UserExistsByEmail(ctx, *body.Email)
		if err != nil || exists {
			return echo.NewHTTPError(409, "Conflict, user with this email exists")
		}

		_, err = database.Q.CreateConfirmationCode(ctx, database.CreateConfirmationCodeParams{
			Recipient: user.Email,
			Code:      utils.GenerateRandomString(32),
			Type:      database.ConfirmationCodeTypeEMAILVERIFICATION,
			UserID:    user.ID,
		})
		if err != nil {
			return echo.NewHTTPError(400, "Failed to create confirmation code")
		}
	}

	if body.Avatar != nil {
		key, err := storage.MoveFile(*body.Avatar, fmt.Sprintf("user/%d", user.ID))
		if err != nil {
			return err
		}

		params.AvatarDoUpdate = true
		params.Avatar = key
	}

	user, err = database.Q.UpdateUser(ctx, params)
	if err != nil {
		return err
	}

	return c.JSON(200, HideKeys(user.ResolveKey(), "Password"))
}
