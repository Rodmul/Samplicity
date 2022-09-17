alter table users
    add column name text;

alter table liked_samples
    rename to library;

drop table created_samples;