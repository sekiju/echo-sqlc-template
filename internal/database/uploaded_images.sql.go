// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: uploaded_images.sql

package database

import (
	"context"
)

const createUploadedImage = `-- name: CreateUploadedImage :one
insert into uploaded_images(hash, key, size, extension, height, width, user_id)
values ($1, $2, $3, $4, $5, $6, $7)
returning id, created_at, hash, key, size, extension, height, width, user_id
`

type CreateUploadedImageParams struct {
	Hash      string `json:"hash"`
	Key       string `json:"key"`
	Size      int32  `json:"size"`
	Extension string `json:"extension"`
	Height    int32  `json:"height"`
	Width     int32  `json:"width"`
	UserID    int32  `json:"userId"`
}

func (q *Queries) CreateUploadedImage(ctx context.Context, arg CreateUploadedImageParams) (UploadedImage, error) {
	row := q.db.QueryRow(ctx, createUploadedImage,
		arg.Hash,
		arg.Key,
		arg.Size,
		arg.Extension,
		arg.Height,
		arg.Width,
		arg.UserID,
	)
	var i UploadedImage
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Hash,
		&i.Key,
		&i.Size,
		&i.Extension,
		&i.Height,
		&i.Width,
		&i.UserID,
	)
	return i, err
}

const deleteUploadedImageById = `-- name: DeleteUploadedImageById :exec
delete
from uploaded_images
where id = $1
`

func (q *Queries) DeleteUploadedImageById(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteUploadedImageById, id)
	return err
}

const getUploadedImage = `-- name: GetUploadedImage :one
select id, created_at, hash, key, size, extension, height, width, user_id
from uploaded_images
where (id = $1 or $1 is null)
   or (hash = $2 or $2 is null)
limit 1
`

type GetUploadedImageParams struct {
	ID   int32  `json:"id"`
	Hash string `json:"hash"`
}

func (q *Queries) GetUploadedImage(ctx context.Context, arg GetUploadedImageParams) (UploadedImage, error) {
	row := q.db.QueryRow(ctx, getUploadedImage, arg.ID, arg.Hash)
	var i UploadedImage
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Hash,
		&i.Key,
		&i.Size,
		&i.Extension,
		&i.Height,
		&i.Width,
		&i.UserID,
	)
	return i, err
}