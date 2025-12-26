CREATE IF NOT EXISTS users (
    id bigint generated always as identity primary key,
    email varchar(255) not null unique,
    password_hash text not null,
    created_at TIMESTAMPTZ not null default now(),
);