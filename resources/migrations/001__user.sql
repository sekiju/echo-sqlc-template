create type user_role as enum ('USER', 'MODERATOR', 'ADMINISTRATOR');

create table users
(
    id         serial primary key,

    enabled    boolean      not null default false,
    email      varchar(320) not null unique,
    username   varchar(64)  not null unique,
    password   varchar      not null,
    role       user_role    not null default 'USER',
    avatar     varchar,

    created_at timestamp    not null default now(),
    updated_at timestamp,
    version    int          not null default 0
);

create table tokens
(
    id            serial primary key,

    access_token  varchar(48) not null unique,
    refresh_token varchar(64) not null unique,

    user_id       int         not null references users (id),

    expired_at    timestamp   not null,
    created_at    timestamp   not null default now(),
    updated_at    timestamp,
    version       int         not null default 0
);

create index idx_user_id_on_tokens on tokens (user_id);

create type confirmation_code_type as enum ('ACTIVATE', 'EMAIL_VERIFICATION', 'PASSWORD_RESET');

create table confirmation_codes
(
    id         serial primary key,
    created_at timestamp              not null default now(),

    recipient  varchar(320)           not null,
    code       varchar(32)            not null unique,
    type       confirmation_code_type not null,

    user_id    int                    not null references users (id)
);

create index idx_code_on_confirmation_codes on confirmation_codes (code);

---- create above / drop below ----

drop index idx_user_id_on_tokens;
drop index idx_code_on_confirmation_codes;
drop table users;
drop table tokens;
drop table confirmation_codes;
drop type user_role;
drop type confirmation_code_type;