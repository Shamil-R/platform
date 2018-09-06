-- переменные
DECLARE @dbname VARCHAR(8) = 'platform'

-- удаление и создание БД
USE master
IF EXISTS ( SELECT * FROM sys.databases WHERE name = 'platform' )
BEGIN
    DROP DATABASE platform
END
CREATE DATABASE platform
GO

USE platform
-- создание таблицы users
IF OBJECT_ID('users', 'U') IS NOT NULL
BEGIN
    DROP TABLE users
END
CREATE TABLE users (
    id INT,
    uid NVARCHAR(50),
    name NVARCHAR(50)
)
GO

INSERT INTO users (id, uid, name) VALUES(1, '1', 'Test')
GO