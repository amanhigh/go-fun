
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
create table aman.MeraProduct
(
  id int(10) unsigned auto_increment
    primary key,
  created_at timestamp null,
  updated_at timestamp null,
  deleted_at timestamp null,
  code varchar(255) null,
  price int(10) unsigned null,
  vertical_id int(10) unsigned null
)
;

create index idx_MeraProduct_deleted_at
  on MeraProduct (deleted_at)
;

create table aman.verticals
(
  id int(10) unsigned auto_increment
    primary key,
  created_at timestamp null,
  updated_at timestamp null,
  deleted_at timestamp null,
  name varchar(255) default 'Shirts' null,
  constraint name
  unique (name)
)
;

create index idx_verticals_deleted_at
  on verticals (deleted_at)
;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE aman.MeraProduct;
DROP TABLE aman.verticals;