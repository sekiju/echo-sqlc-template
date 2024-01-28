create table uploaded_images
(
    id         serial primary key,
    created_at timestamp   not null default now(),

    hash       varchar(32) not null unique,
    key        varchar     not null,
    size       int         not null,
    extension  varchar(6)  not null,
    height     int         not null,
    width      int         not null,

    user_id    int         not null references users (id)
);

---- create above / drop below ----

drop table uploaded_images;