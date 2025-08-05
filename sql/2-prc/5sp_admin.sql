-- pass table name as string, get count of table
create or replace function fn_cnt(_tbl regclass, out cnt bigint)
returns bigint
language plpgsql
as $$
begin
	execute format('select count(*) from %s', _tbl) into cnt;
end; $$;

select fn_cnt('stats.tbox');