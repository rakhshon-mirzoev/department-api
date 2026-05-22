-- +goose Up
ALTER TABLE departments DROP CONSTRAINT departments_parent_id_fkey;
ALTER TABLE departments ADD CONSTRAINT departments_parent_id_fkey
    FOREIGN KEY (parent_id) REFERENCES departments(id) ON DELETE CASCADE;

ALTER TABLE employees DROP CONSTRAINT employees_department_id_fkey;
ALTER TABLE employees ADD CONSTRAINT employees_department_id_fkey
    FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE departments DROP CONSTRAINT departments_parent_id_fkey;
ALTER TABLE departments ADD CONSTRAINT departments_parent_id_fkey
    FOREIGN KEY (parent_id) REFERENCES departments(id);

ALTER TABLE employees DROP CONSTRAINT employees_department_id_fkey;
ALTER TABLE employees ADD CONSTRAINT employees_department_id_fkey
    FOREIGN KEY (department_id) REFERENCES departments(id);
