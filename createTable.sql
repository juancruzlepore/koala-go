create table next_dates (
  date timestamptz
);

create table movies (
  creation_date timestamptz,
  added_by text,
  name text,
  seen bool
);

insert into movies values (current_timestamp, 'Juan', 'Ant-man', false);

alter table next_dates add column end_date timestamptz;

ALTER TABLE next_dates
  RENAME COLUMN date TO start_date;