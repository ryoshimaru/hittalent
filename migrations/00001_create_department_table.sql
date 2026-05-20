-- +goose Up
CREATE TABLE departments (
    id SERIAL PRIMARY KEY, 
    name VARCHAR(200) NOT NULL,
    parent_id INT NULL, 
    created_at TIMESTAMP NOT NULL DEFAULT NOW(), 

    CONSTRAINT fk_departments_parent
        FOREIGN KEY (parent_id)
        REFERENCES departments(id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS departments;
