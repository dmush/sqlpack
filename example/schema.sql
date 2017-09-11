drop schema if exists example cascade;
create schema example;

set search_path = example;

create table example(id serial primary key, name text);