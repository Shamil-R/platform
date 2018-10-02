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
