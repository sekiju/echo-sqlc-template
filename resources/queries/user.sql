-- name: GetUser :one
select *
from users
where (username = @username or @username is null)
   or (email = @email or @email is null)
   or (id = @id or @id is null)
    limit 1;

-- name: UserExistsByEmail :one
select exists (
    select 1
    from users
    where email = @email
);

-- name: UserExistsByUsername :one
select exists (
    select 1
    from users
    where username = @username
);

-- name: CreateUser :one
insert into users(email, username, password) values ($1, $2, $3) returning *;

-- name: UpdateUser :one
update users
set enabled    = case when @enabled_do_update::boolean then @enabled::boolean else enabled end,
    email      = case when @email_do_update::boolean then @email::varchar(320) else email end,
    username   = case when @username_do_update::boolean then @username::varchar(64) else username end,
    password   = case when @password_do_update::boolean then @password::varchar else password end,
    avatar     = case when @avatar_do_update::boolean then @avatar::varchar else avatar end,
    updated_at = now(),
    version    = version + 1
where id = @id
returning *;

-- name: DeleteUserByID :exec
delete
from users
where id = $1;
