-- name: GetToken :one
select *
from tokens
where (id = @id or @id is null)
   or (access_token = @access_token or @access_token is null)
   or (refresh_token = @refresh_token or @refresh_token is null)
limit 1;

-- name: GetUserByTokenID :one
select u.*
from users u
         join tokens t on u.id = t.user_id
where t.id = @id
limit 1;

-- name: CreateToken :one
insert into tokens(access_token, refresh_token, user_id, expired_at)
values ($1, $2, $3, $4)
returning *;

-- name: UpdateToken :one
update tokens
set access_token  = case when @access_token_do_update::boolean then @access_token::varchar(48) else access_token end,
    refresh_token = case when @refresh_token_do_update::boolean then @refresh_token::varchar(64) else refresh_token end,
    expired_at    = case when @expired_at_do_update::boolean then @expired_at::timestamp else expired_at end,
    updated_at    = now(),
    version       = version + 1
where id = @id
returning *;

-- name: DeleteTokenByID :exec
delete
from tokens
where id = $1;