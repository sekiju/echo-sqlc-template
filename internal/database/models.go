// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package database

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type ConfirmationCodeType string

const (
	ConfirmationCodeTypeACTIVATE          ConfirmationCodeType = "ACTIVATE"
	ConfirmationCodeTypeEMAILVERIFICATION ConfirmationCodeType = "EMAIL_VERIFICATION"
	ConfirmationCodeTypePASSWORDRESET     ConfirmationCodeType = "PASSWORD_RESET"
)

func (e *ConfirmationCodeType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = ConfirmationCodeType(s)
	case string:
		*e = ConfirmationCodeType(s)
	default:
		return fmt.Errorf("unsupported scan type for ConfirmationCodeType: %T", src)
	}
	return nil
}

type NullConfirmationCodeType struct {
	ConfirmationCodeType ConfirmationCodeType `json:"confirmationCodeType"`
	Valid                bool                 `json:"valid"` // Valid is true if ConfirmationCodeType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullConfirmationCodeType) Scan(value interface{}) error {
	if value == nil {
		ns.ConfirmationCodeType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.ConfirmationCodeType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullConfirmationCodeType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.ConfirmationCodeType), nil
}

type UserRole string

const (
	UserRoleUSER          UserRole = "USER"
	UserRoleMODERATOR     UserRole = "MODERATOR"
	UserRoleADMINISTRATOR UserRole = "ADMINISTRATOR"
)

func (e *UserRole) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = UserRole(s)
	case string:
		*e = UserRole(s)
	default:
		return fmt.Errorf("unsupported scan type for UserRole: %T", src)
	}
	return nil
}

type NullUserRole struct {
	UserRole UserRole `json:"userRole"`
	Valid    bool     `json:"valid"` // Valid is true if UserRole is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullUserRole) Scan(value interface{}) error {
	if value == nil {
		ns.UserRole, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.UserRole.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullUserRole) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.UserRole), nil
}

type ConfirmationCode struct {
	ID        int32                `json:"id"`
	CreatedAt pgtype.Timestamp     `json:"createdAt"`
	Recipient string               `json:"recipient"`
	Code      string               `json:"code"`
	Type      ConfirmationCodeType `json:"type"`
	UserID    int32                `json:"userId"`
}

type Token struct {
	ID           int32            `json:"id"`
	AccessToken  string           `json:"accessToken"`
	RefreshToken string           `json:"refreshToken"`
	UserID       int32            `json:"userId"`
	ExpiredAt    pgtype.Timestamp `json:"expiredAt"`
	CreatedAt    pgtype.Timestamp `json:"createdAt"`
	UpdatedAt    pgtype.Timestamp `json:"updatedAt"`
	Version      int32            `json:"version"`
}

type UploadedImage struct {
	ID        int32            `json:"id"`
	CreatedAt pgtype.Timestamp `json:"createdAt"`
	Hash      string           `json:"hash"`
	Key       string           `json:"key"`
	Size      int32            `json:"size"`
	Extension string           `json:"extension"`
	Height    int32            `json:"height"`
	Width     int32            `json:"width"`
	UserID    int32            `json:"userId"`
}

type User struct {
	ID        int32            `json:"id"`
	Enabled   bool             `json:"enabled"`
	Email     string           `json:"email"`
	Username  string           `json:"username"`
	Password  string           `json:"password"`
	Role      UserRole         `json:"role"`
	Avatar    pgtype.Text      `json:"avatar"`
	CreatedAt pgtype.Timestamp `json:"createdAt"`
	UpdatedAt pgtype.Timestamp `json:"updatedAt"`
	Version   int32            `json:"version"`
}
