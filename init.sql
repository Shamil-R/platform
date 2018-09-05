CREATE DATABASE platform

GO

USE platform

CREATE TABLE users (
    id INT,
    name NVARCHAR(50)
)

INSERT INTO users (id, name) VALUES(1, 'Test')

GO