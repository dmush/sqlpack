{{ include "./schema" "./get_example_name" }}

do $$
declare
  v_id numeric;
  v_name text := 'example';
begin
  insert into example(name)
  values (v_name)
  returning id into v_id;
  
  if v_name is distinct from get_example_name(v_id) then
    raise 'invalid name';
  end if;
end;
$$;
