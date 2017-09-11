package main

import (
	"strings"
	"testing"
)

func TestPathByIndex(t *testing.T) {
	s := strings.TrimSpace(`
-- # start # example/example.sql
drop schema if exists public;

-- # start # example/schema.sql
drop schema if exists example cascade;
create schema example;

set search_path = example;

create table example(id serial primary key, name text);
-- # end # example/schema.sql
-- # start # example/get_example_name.sql
create function get_example_name(p_id numeric) returns text as
$$
select name
from example
where id = p_id;
$$ 
language sql;
-- # end # example/get_example_name.sql

truncate table example;
-- # end # example/example.sql
	`)

	path := pathByIndex(s, 55)
	if path != "example/example.sql" {
		t.Errorf("invalid start path %s", path)
	}

	path = pathByIndex(s, 150)
	if path != "example/schema.sql" {
		t.Errorf("invalid schema path '%s'", path)
	}

	path = pathByIndex(s, 390)
	if path != "example/get_example_name.sql" {
		t.Errorf("invalid get_example_name path '%s'", path)
	}

	path = pathByIndex(s, 500)
	if path != "example/example.sql" {
		t.Errorf("invalid end path '%s'", path)
	}

	text := textByIndex(s, 55)
	if text != "drop schema if exists "+colorFailChar.Sprint("p")+"ublic;" {
		t.Errorf("invalid root text '%s'", text)
	}
}
