-- создание логина
IF NOT EXISTS (SELECT * FROM master.sys.server_principals WHERE name = 'platform')
BEGIN
	CREATE LOGIN platform WITH PASSWORD = 'p@sSw0rd'
END

-- удаление и создание БД
IF EXISTS (SELECT * FROM master.sys.databases WHERE name = 'platform')
BEGIN
    DROP DATABASE platform
END
CREATE DATABASE platform

USE platform

-- создание пользователя бд
IF NOT EXISTS (SELECT * FROM master.sys.database_principals WHERE name = 'platform')
BEGIN
	-- CREATE USER platform FOR LOGIN platform WITH DEFAULT_SCHEMA = dbo
	EXEC sp_addrolemember 'db_owner', 'platform'
END

-- создание таблицы user
IF OBJECT_ID('user', 'U') IS NOT NULL
BEGIN
    DROP TABLE [user]
END
CREATE TABLE [user] (
    id INT IDENTITY(1,1) NOT NULL PRIMARY KEY,
    name NVARCHAR(50)
)

-- INSERT INTO user (id, name) VALUES(1, 'Test')
-- GO

IF OBJECT_ID('object_access', 'U') IS NOT NULL
BEGIN
    DROP TABLE [object_access]
END
CREATE TABLE [object_access] (
    id INT IDENTITY(1,1) NOT NULL PRIMARY KEY,
    user_id INT NOT NULL,
    object NVARCHAR(50) NOT NULL
)

--INSERT INTO object_access (user_id, object) VALUES(1, 'User')

IF OBJECT_ID('action_object_access', 'U') IS NOT NULL
BEGIN
    DROP TABLE [action_object_access]
END
CREATE TABLE [action_object_access] (
    id INT IDENTITY(1,1) NOT NULL PRIMARY KEY,
    object_id INT NOT NULL,
    action NVARCHAR(8),
    allow BIT NOT NULL
)

--INSERT INTO action_object_access (object_id, allow) VALUES(1, 1)

IF OBJECT_ID('field_access', 'U') IS NOT NULL
BEGIN
    DROP TABLE [field_access]
END
CREATE TABLE [field_access] (
    id INT IDENTITY(1,1) NOT NULL PRIMARY KEY,
    object_id INT NOT NULL,
    field NVARCHAR(50) NOT NULL
)

IF OBJECT_ID('action_field_access', 'U') IS NOT NULL
BEGIN
    DROP TABLE [action_field_access]
END
CREATE TABLE [action_field_access] (
    id INT IDENTITY(1,1) NOT NULL PRIMARY KEY,
    field_id INT NOT NULL,
    action NVARCHAR(8),
    allow BIT NOT NULL
)
