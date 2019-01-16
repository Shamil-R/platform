USE master

-- создание логина
IF NOT EXISTS (SELECT * FROM master.sys.server_principals WHERE name = 'platform')
BEGIN
	CREATE LOGIN platform WITH PASSWORD = 'p@sSw0rd'
END
GO

-- удаление и создание БД
IF EXISTS (SELECT * FROM master.sys.databases WHERE name = 'platform')
BEGIN
    DROP DATABASE platform
END
CREATE DATABASE platform
GO

USE platform

-- создание пользователя бд
IF NOT EXISTS (SELECT * FROM master.sys.database_principals WHERE name = 'platform')
BEGIN
	CREATE USER platform FOR LOGIN platform WITH DEFAULT_SCHEMA = dbo
	EXEC sp_addrolemember 'db_owner', 'platform'
END
GO

-- создание таблицы user
IF OBJECT_ID('user', 'U') IS NOT NULL
BEGIN
    DROP TABLE [user]
END
CREATE TABLE [user] (
    id INT IDENTITY(1,1) NOT NULL PRIMARY KEY,
    name NVARCHAR(50),
    bio NVARCHAR(100)
)
GO

-- создание таблицы material
IF OBJECT_ID('material', 'U') IS NOT NULL
BEGIN
    DROP TABLE [material]
END
CREATE TABLE [material] (
    id INT IDENTITY(1,1) NOT NULL PRIMARY KEY,
    title NVARCHAR(50),
    user_id INT NOT NULL,
    CONSTRAINT FK_material_user FOREIGN KEY (user_id) REFERENCES [user]([id])
)
GO

-- создание таблицы order
IF OBJECT_ID('order', 'U') IS NOT NULL
BEGIN
    DROP TABLE [order]
END
CREATE TABLE [order] (
    id INT IDENTITY(1,1) NOT NULL PRIMARY KEY,
    date NVARCHAR(50),
    user_id INT NOT NULL,
    CONSTRAINT FK_order_user FOREIGN KEY (user_id) REFERENCES [user]([id])
)
GO

-- создание таблицы item
IF OBJECT_ID('item', 'U') IS NOT NULL
BEGIN
    DROP TABLE [item]
END
CREATE TABLE [item] (
    id INT IDENTITY(1,1) NOT NULL PRIMARY KEY,
    count INT,
    order_id INT NOT NULL,
    material_id INT NOT NULL,
    CONSTRAINT FK_item_order FOREIGN KEY (order_id) REFERENCES [order]([id]),
    CONSTRAINT FK_item_material FOREIGN KEY (material_id) REFERENCES [material]([id])
)
GO

-- создание таблицы trip
IF OBJECT_ID('trip', 'U') IS NOT NULL
BEGIN
    DROP TABLE [trip]
END
CREATE TABLE [trip] (
    id INT IDENTITY(1,1) NOT NULL PRIMARY KEY,
    city NVARCHAR(50),
    user_id INT NOT NULL,
    order_id INT NOT NULL,
    CONSTRAINT FK_trip_user FOREIGN KEY (user_id) REFERENCES [user]([id]),
    CONSTRAINT FK_trip_order FOREIGN KEY (order_id) REFERENCES [order]([id])
)
GO
