-- name: GetUploadedImage :one
select *
from uploaded_images
where (id = @id or @id is null)
   or (hash = @hash or @hash is null)
limit 1;

-- name: CreateUploadedImage :one
insert into uploaded_images(hash, key, size, extension, height, width, user_id)
values ($1, $2, $3, $4, $5, $6, $7)
returning *;

-- name: DeleteUploadedImageById :exec
delete
from uploaded_images
where id = $1;
