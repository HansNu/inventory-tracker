select 'CREATE TABLE ' || c.relname || E' (\n' ||
  string_agg(
    '  ' || quote_ident(a.attname) || ' ' || format_type(a.atttypid, a.atttypmod)
      || coalesce(' DEFAULT ' || pg_get_expr(d.adbin, d.adrelid), '')
      || case when a.attnotnull then ' NOT NULL' else '' end,
    E',\n' order by a.attnum
  ) || E'\n);' as ddl
from pg_class c
join pg_namespace n on n.oid = c.relnamespace
join pg_attribute a on a.attrelid = c.oid and a.attnum > 0 and not a.attisdropped
left join pg_attrdef d on d.adrelid = c.oid and d.adnum = a.attnum
where n.nspname = 'public' and c.relkind = 'r'
group by c.relname
order by c.relname;

select indexdef || ';' as ddl
from pg_indexes
where schemaname = 'public'
  and indexname not in (
    select conname from pg_constraint where contype in ('p','u')
  )
order by tablename, indexname;

select indexdef || ';' as ddl
from pg_indexes
where schemaname = 'public'
  and indexname not in (
    select conname from pg_constraint where contype in ('p','u')
  )
order by tablename, indexname;