{{ include "./schema" }}

do $$
declare
  v_id numeric;
  v_name text := 'example';
begin
  insert into example(name)
  values (v_name)
  returning id into v_id;
  
  if not exists (
    select *
    from example
    where id = v_id and name = v_name
  ) then
    raise 'record not found';
  end if;
end;
$$;