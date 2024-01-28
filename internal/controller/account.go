package controller

import (
	"echo-sqlc-template/internal/database"
	"echo-sqlc-template/internal/services/mail"
	"echo-sqlc-template/internal/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type AccountController struct{}

func NewAccountController(e *echo.Echo) {
	c := AccountController{}

	group := e.Group("/account")
	group.POST("/login", c.Login)
	group.POST("/register", c.Register)
	group.POST("/activate", c.Activate)
	group.POST("/activate/resend", c.ActivateResend)
	group.POST("/password", c.Password)
	group.POST("/password/resend", c.PasswordResend)
	group.POST("/email", c.Email)
}

type LoginRequest struct {
	EmailOrUsername string `json:"emailOrUsername" validate:"required,min=6,max=64"`
	Password        string `json:"password" validate:"required,min=6,max=64"`
}

func (AccountController) Login(c echo.Context) error {
	body, err := Body[LoginRequest](c)
	if err != nil {
		return err
	}

	user, err := database.Q.GetUser(ctx, database.GetUserParams{
		Email:    body.EmailOrUsername,
		Username: body.EmailOrUsername,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid username or e-mail")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid password")
	}

	token, err := database.Q.CreateToken(ctx, database.CreateTokenParams{
		AccessToken:  utils.GenerateRandomString(48),
		RefreshToken: utils.GenerateRandomString(64),
		UserID:       user.ID,
		ExpiredAt: pgtype.Timestamp{
			Time:  time.Now().Local().Add(time.Hour * time.Duration(1)),
			Valid: true,
		},
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to create tokens")
	}

	return c.JSON(200, token)
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=6,max=64"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=64"`
}

func (AccountController) Register(c echo.Context) error {
	body, err := Body[RegisterRequest](c)
	if err != nil {
		return err
	}

	exists, err := database.Q.UserExistsByUsername(ctx, body.Username)
	if err != nil || exists {
		return echo.NewHTTPError(409, "Conflict, user with this username exists")
	}

	exists, err = database.Q.UserExistsByEmail(ctx, body.Email)
	if err != nil || exists {
		return echo.NewHTTPError(409, "Conflict, user with this email exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(400, "Failed to hash password")
	}

	user, err := database.Q.CreateUser(ctx, database.CreateUserParams{
		Username: body.Username,
		Password: string(hashedPassword),
		Email:    body.Email,
	})
	if err != nil {
		return echo.NewHTTPError(400, "Failed to create user")
	}

	code, err := database.Q.CreateConfirmationCode(ctx, database.CreateConfirmationCodeParams{
		Recipient: user.Email,
		Code:      utils.GenerateRandomString(32),
		Type:      database.ConfirmationCodeTypeACTIVATE,
		UserID:    user.ID,
	})
	if err != nil {
		return echo.NewHTTPError(400, "Failed to create confirmation code")
	}

	err = mail.SendCode(user.Email, code.Code, "Активация аккаунта")
	if err != nil {
		return echo.NewHTTPError(400, "Failed to send code to recipient email")
	}

	return c.JSON(201, HideKeys(user, "Password"))
}

type CodeRequest struct {
	Code string `json:"code" validate:"required,len=32"`
}

func (AccountController) Activate(c echo.Context) error {
	body, err := Body[CodeRequest](c)
	if err != nil {
		return err
	}

	code, err := database.Q.GetConfirmationCodeByTypeAndCode(ctx, database.GetConfirmationCodeByTypeAndCodeParams{
		Code: body.Code,
		Type: database.ConfirmationCodeTypeACTIVATE,
	})
	if err != nil {
		return echo.NewHTTPError(404, "Code not found")
	}

	user, err := database.Q.GetUser(ctx, database.GetUserParams{ID: code.UserID})
	if err != nil || user.Enabled {
		return echo.NewHTTPError(400, "User with this email doesn't exists or account activated")
	}

	user, err = database.Q.UpdateUser(ctx, database.UpdateUserParams{
		ID:              user.ID,
		EnabledDoUpdate: true,
		Enabled:         true,
	})
	if err != nil {
		return err
	}

	err = database.Q.DeleteConfirmationCodeByID(ctx, code.ID)
	if err != nil {
		return err
	}

	token, err := database.Q.CreateToken(ctx, database.CreateTokenParams{
		AccessToken:  utils.GenerateRandomString(48),
		RefreshToken: utils.GenerateRandomString(64),
		UserID:       user.ID,
		ExpiredAt: pgtype.Timestamp{
			Time:  time.Now().Local().Add(time.Hour * time.Duration(1)),
			Valid: true,
		},
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to create tokens")
	}

	return c.JSON(200, token)
}

type ResendRequest struct {
	Email string `json:"email" validate:"email"`
}

func (AccountController) ActivateResend(c echo.Context) error {
	body, err := Body[ResendRequest](c)
	if err != nil {
		return err
	}

	user, err := database.Q.GetUser(ctx, database.GetUserParams{Email: body.Email})
	if err != nil || user.Enabled {
		return echo.NewHTTPError(400, "User with this email doesn't exists or account activated")
	}

	exists, err := database.Q.ConfirmationCodeRecentlyExists(ctx, database.ConfirmationCodeRecentlyExistsParams{
		UserID: user.ID,
		Type:   database.ConfirmationCodeTypeACTIVATE,
	})
	if err != nil || exists {
		return echo.NewHTTPError(400, "Wait 15 minutes for create new code")
	}

	code, err := database.Q.CreateConfirmationCode(ctx, database.CreateConfirmationCodeParams{
		Recipient: user.Email,
		Code:      utils.GenerateRandomString(32),
		UserID:    user.ID,
		Type:      database.ConfirmationCodeTypeACTIVATE,
	})
	if err != nil {
		return echo.NewHTTPError(400, "Failed to create confirmation code")
	}

	err = mail.SendCode(user.Email, code.Code, "Активация аккаунта")
	if err != nil {
		return echo.NewHTTPError(400, "Failed to send code to recipient email")
	}

	return c.NoContent(204)
}

type PasswordRequest struct {
	Code     string `json:"code" validate:"required,len=32"`
	Password string `json:"password" validate:"required,min=6,max=64"`
}

func (AccountController) Password(c echo.Context) error {
	body, err := Body[PasswordRequest](c)
	if err != nil {
		return err
	}

	code, err := database.Q.GetConfirmationCodeByTypeAndCode(ctx, database.GetConfirmationCodeByTypeAndCodeParams{
		Code: body.Code,
		Type: database.ConfirmationCodeTypePASSWORDRESET,
	})
	if err != nil {
		return echo.NewHTTPError(404, "Code not found")
	}

	user, err := database.Q.GetUser(ctx, database.GetUserParams{ID: code.UserID})
	if err != nil || !user.Enabled {
		return echo.NewHTTPError(400, "User with this email doesn't exists or account unactivated")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(400, "Failed to hash password")
	}

	user, err = database.Q.UpdateUser(ctx, database.UpdateUserParams{
		ID:               user.ID,
		PasswordDoUpdate: true,
		Password:         string(hashedPassword),
	})
	if err != nil {
		return err
	}

	err = database.Q.DeleteConfirmationCodeByID(ctx, code.ID)
	if err != nil {
		return err
	}

	token, err := database.Q.CreateToken(ctx, database.CreateTokenParams{
		AccessToken:  utils.GenerateRandomString(48),
		RefreshToken: utils.GenerateRandomString(64),
		UserID:       user.ID,
		ExpiredAt: pgtype.Timestamp{
			Time:  time.Now().Local().Add(time.Hour * time.Duration(1)),
			Valid: true,
		},
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to create tokens")
	}

	return c.JSON(200, token)
}

func (AccountController) PasswordResend(c echo.Context) error {
	body, err := Body[ResendRequest](c)
	if err != nil {
		return err
	}

	user, err := database.Q.GetUser(ctx, database.GetUserParams{Email: body.Email})
	if err != nil || !user.Enabled {
		return echo.NewHTTPError(400, "User with this email doesn't exists or account unactivated")
	}

	exists, err := database.Q.ConfirmationCodeRecentlyExists(ctx, database.ConfirmationCodeRecentlyExistsParams{
		UserID: user.ID,
		Type:   database.ConfirmationCodeTypePASSWORDRESET,
	})
	if err != nil || exists {
		return echo.NewHTTPError(400, "Wait 15 minutes for create new code")
	}

	code, err := database.Q.CreateConfirmationCode(ctx, database.CreateConfirmationCodeParams{
		Recipient: user.Email,
		Code:      utils.GenerateRandomString(32),
		UserID:    user.ID,
		Type:      database.ConfirmationCodeTypePASSWORDRESET,
	})
	if err != nil {
		return echo.NewHTTPError(400, "Failed to create confirmation code")
	}

	err = mail.SendCode(user.Email, code.Code, "Сброс пароля")
	if err != nil {
		return echo.NewHTTPError(400, "Failed to send code to recipient email")
	}

	return c.NoContent(204)
}

func (AccountController) Email(c echo.Context) error {
	body, err := Body[CodeRequest](c)
	if err != nil {
		return err
	}

	code, err := database.Q.GetConfirmationCodeByTypeAndCode(ctx, database.GetConfirmationCodeByTypeAndCodeParams{
		Code: body.Code,
		Type: database.ConfirmationCodeTypeEMAILVERIFICATION,
	})
	if err != nil {
		return echo.NewHTTPError(404, "Code not found")
	}

	user, err := database.Q.GetUser(ctx, database.GetUserParams{ID: code.UserID})
	if err != nil || !user.Enabled {
		return echo.NewHTTPError(400, "User with this email doesn't exists or account unactivated")
	}

	user, err = database.Q.UpdateUser(ctx, database.UpdateUserParams{
		ID:            user.ID,
		EmailDoUpdate: true,
		Email:         code.Recipient,
	})
	if err != nil {
		return err
	}

	err = database.Q.DeleteConfirmationCodeByID(ctx, code.ID)
	if err != nil {
		return err
	}

	return c.NoContent(204)
}
