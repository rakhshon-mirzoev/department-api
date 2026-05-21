-- +goose Up
CREATE TABLE departments(
    id serial primary key,
    name varchar(200) NOT NULL,
    parent_id int references departments(id),
    created_at timestamptz not null default now(),
    UNIQUE(parent_id, name)
);

CREATE TABLE employees(
    id serial primary key,
    department_id int references departments(id),
    full_name varchar(200) not null,
    position varchar(200) not null,
    hired_at timestamptz,
    created_at timestamptz not null default now()
);


-- +goose Down
DROP TABLE IF EXISTS employees;
DROP TABLE IF EXISTS departments;