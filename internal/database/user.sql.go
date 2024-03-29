// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: user.sql

package database

import (
	"context"
)

const createUser = `-- name: CreateUser :one
insert into users(email, username, password) values ($1, $2, $3) returning id, enabled, email, username, password, role, avatar, created_at, updated_at, version
`

type CreateUserParams struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.Email, arg.Username, arg.Password)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Enabled,
		&i.Email,
		&i.Username,
		&i.Password,
		&i.Role,
		&i.Avatar,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Version,
	)
	return i, err
}

const deleteUserByID = `-- name: DeleteUserByID :exec
delete
from users
where id = $1
`

func (q *Queries) DeleteUserByID(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteUserByID, id)
	return err
}

const getUser = `-- name: GetUser :one
select id, enabled, email, username, password, role, avatar, created_at, updated_at, version
from users
where (username = $1 or $1 is null)
   or (email = $2 or $2 is null)
   or (id = $3 or $3 is null)
    limit 1
`

type GetUserParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	ID       int32  `json:"id"`
}

func (q *Queries) GetUser(ctx context.Context, arg GetUserParams) (User, error) {
	row := q.db.QueryRow(ctx, getUser, arg.Username, arg.Email, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Enabled,
		&i.Email,
		&i.Username,
		&i.Password,
		&i.Role,
		&i.Avatar,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Version,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
update users
set enabled    = case when $1::boolean then $2::boolean else enabled end,
    email      = case when $3::boolean then $4::varchar(320) else email end,
    username   = case when $5::boolean then $6::varchar(64) else username end,
    password   = case when $7::boolean then $8::varchar else password end,
    avatar     = case when $9::boolean then $10::varchar else avatar end,
    updated_at = now(),
    version    = version + 1
where id = $11
returning id, enabled, email, username, password, role, avatar, created_at, updated_at, version
`

type UpdateUserParams struct {
	EnabledDoUpdate  bool   `json:"enabledDoUpdate"`
	Enabled          bool   `json:"enabled"`
	EmailDoUpdate    bool   `json:"emailDoUpdate"`
	Email            string `json:"email"`
	UsernameDoUpdate bool   `json:"usernameDoUpdate"`
	Username         string `json:"username"`
	PasswordDoUpdate bool   `json:"passwordDoUpdate"`
	Password         string `json:"password"`
	AvatarDoUpdate   bool   `json:"avatarDoUpdate"`
	Avatar           string `json:"avatar"`
	ID               int32  `json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser,
		arg.EnabledDoUpdate,
		arg.Enabled,
		arg.EmailDoUpdate,
		arg.Email,
		arg.UsernameDoUpdate,
		arg.Username,
		arg.PasswordDoUpdate,
		arg.Password,
		arg.AvatarDoUpdate,
		arg.Avatar,
		arg.ID,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Enabled,
		&i.Email,
		&i.Username,
		&i.Password,
		&i.Role,
		&i.Avatar,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Version,
	)
	return i, err
}

const userExistsByEmail = `-- name: UserExistsByEmail :one
select exists (
    select 1
    from users
    where email = $1
)
`

func (q *Queries) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	row := q.db.QueryRow(ctx, userExistsByEmail, email)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const userExistsByUsername = `-- name: UserExistsByUsername :one
select exists (
    select 1
    from users
    where username = $1
)
`

func (q *Queries) UserExistsByUsername(ctx context.Context, username string) (bool, error) {
	row := q.db.QueryRow(ctx, userExistsByUsername, username)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}
