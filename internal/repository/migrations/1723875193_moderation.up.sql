IF NOT EXISTS create type status as enum('reported','resolved','deleted');

alter table resolutions
add status status
constraint default_value
default ('reported');