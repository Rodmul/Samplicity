create table users
(
    id serial not null primary key,
    name text not null,
    username text not null unique,
    password_hash text not null,
    email text
);

create table samples
(
    id serial not null primary key,
    name text,
    path text not null,
    author text,
    type text
);

create table library
(
    id serial not null primary key,
    sample_id int references samples (id) on delete cascade not null,
    user_id int references users (id) on delete cascade not null
);
