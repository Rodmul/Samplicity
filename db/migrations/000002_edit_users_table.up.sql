alter table users
    drop column name;

alter table library
    rename to liked_samples;

create table created_samples
(
    id serial not null primary key,
    sample_id int references samples (id) on delete cascade not null,
    user_id int references users (id) on delete cascade not null
)