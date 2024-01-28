-- name: GetConfirmationCodeByID :one
select *
from confirmation_codes
where id = @id
limit 1;

-- name: GetConfirmationCodeByTypeAndCode :one
select *
from confirmation_codes
where type = @type
  and code = @code
  and created_at >= now() - interval '15 minutes'
limit 1;

-- name: ConfirmationCodeRecentlyExists :one
select exists (
    select 1
    from confirmation_codes
    where user_id = @user_id
      and type = @type
      and created_at > now() - interval '15 minutes'
);

-- name: CreateConfirmationCode :one
insert into confirmation_codes(recipient, code, type, user_id)
values ($1, $2, $3, $4)
returning *;

-- name: DeleteConfirmationCodeByID :exec
delete
from confirmation_codes
where id = $1;
