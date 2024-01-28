package database

import (
	"echo-sqlc-template/internal/config"
	"github.com/jackc/pgx/v5/pgtype"
)

func ResolveTextKey(key pgtype.Text) pgtype.Text {
	if key.Valid {
		key.String = config.Data.Storage.Host + key.String
	} else {
		key.String = config.Data.Storage.Host + "/no_avatar.png"
		key.Valid = true
	}

	return key
}

func ResolveStringKey(key string) string {
	return config.Data.Storage.Host + key
}

func (x *User) ResolveKey() User {
	x.Avatar = ResolveTextKey(x.Avatar)
	return *x
}
