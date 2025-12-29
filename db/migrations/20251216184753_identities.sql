-- migrate:up
create table if not exists organizations (
    id varchar(50) not null,
    name varchar(100) not null,
    slug varchar(20) not null unique,
    status varchar(20) not null default 'active',
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    primary key (id)
);
create table if not exists identities (
    id varchar(50) not null,
    email varchar(255) not null unique,
    first_name varchar(100) not null,
    last_name varchar(100) not null,
    status varchar(20) not null default 'active',
    email_verified_at timestamp,
    failed_login_attempts int not null default 0,
    lock_expires_at timestamp not null default now(),
    last_login_at timestamp not null default now(),
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    primary key (id)
);
create table if not exists user_otts (
    id varchar(255) not null,
    -- this is the hashed code
    kind varchar(100) not null,
    principal varchar(100) not null,
    expires_at timestamp not null default now(),
    created_at timestamp not null default now(),
    primary key (id)
);
create table user_sessions (
    id varchar(255) not null,
    -- this is the hashed opaque token
    principal_id varchar(50) not null,
    ip_address varchar(50) not null default '0.0.0.0',
    metadata jsonb not null default '{}',
    expires_at timestamp not null default now(),
    created_at timestamp not null default now(),
    primary key (id)
);
-- migrate:down
drop table if exists user_sessions;
drop table if exists user_otts;
drop table if exists identities;
drop table if exists organizations;