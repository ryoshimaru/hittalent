-- +goose Up
CREATE TABLE employee(
    id SERIAL PRIMARY KEY, 
    department_id INT NOT NULL, 
    full_name VARCHAR(200) NOT NULL, 
    position VARCHAR(200) NOT NULL, 
    hired_at DATE NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_employees_department
        FOREIGN KEY (department_id)
        REFERENCES departments(id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS employee;
