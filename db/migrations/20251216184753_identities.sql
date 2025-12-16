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
    organization_id varchar(50) not null,
    roles jsonb not null default '[]',
    hashed_password varchar(255) not null, -- used only for password authentication
    status varchar(20) not null default 'active',
    email_verified_at timestamp,
    failed_login_attempts int not null default 0,
    lock_expires_at timestamp not null default now(),
    last_login_at timestamp not null default now(),
    comments varchar(255) not null default '',
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    primary key (id)
);

create table if not exists one_time_tokens (
    id varchar(50) not null,
    hashed_code varchar(255) not null,
    category varchar(100) not null,
    identity_id varchar(50) not null,
    used boolean not null default false,
    expires_at timestamp not null default now(),
    created_at timestamp not null default now(),
    primary key (id),
    foreign key (identity_id) references identities(id) on delete cascade
);

create table identity_sessions (
    id varchar(255) not null, -- this is the hashed opaque token
    identity_id varchar(50) not null,
    organization_id varchar(50) not null,
    roles jsonb not null default '[]',
    ip_address varchar(50) not null default '0.0.0.0',
    expires_at timestamp not null default now(),
    created_at timestamp not null default now(),
    primary key (id),
    foreign key (identity_id) references identities(id) on delete cascade,
    foreign key (organization_id) references organizations(id) on delete cascade
);

-- migrate:down
drop table if exists identity_sessions;
drop table if exists one_time_tokens;
drop table if exists identities;
drop table if exists organizations;

