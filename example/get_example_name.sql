create function get_example_name(p_id numeric) returns text as
$$
select name
from example
where id = p_id;
$$ 
language sql;