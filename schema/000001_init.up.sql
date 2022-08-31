create table users
(
    id  integer not null unique,
    name text not null,
    username text not null unique,
    password_hash text not null,
    email text
);

create table samples
(
    id integer,
    name text,
    path text
);

create table library
(
    id integer not null unique,
    sample_id int references samples (id) on delete cascade not null,
    user_id int references users (id) on delete cascade not null
);
