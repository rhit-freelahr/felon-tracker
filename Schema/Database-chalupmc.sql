CREATE DATABASE [FelonTracker-chalupmc]
ON
PRIMARY(
	NAME = [FelonTracker-chalupmc],
	FILENAME = '/var/opt/mssql/data/FelonTracker-chalupmc.mdf'
)
LOG ON(
	NAME = [FelonTracker-chalupmc_log],
	FILENAME = '/var/opt/mssql/data/FelonTracker-chalupmc_log.ldf'
)
GO

USE [FelonTracker-chalupmc]
GO

IF (NOT EXISTS (select name, suser_sname(sid)
	from master.dbo.sysdatabases
	where 'freelahr' is null or name = DB_NAME()))
BEGIN
	CREATE USER freelahr FROM LOGIN freelahr
	exec sp_addrolemember 'db_owner', freelahr;
END
GO

IF (NOT EXISTS (select name, suser_sname(sid)
	from master.dbo.sysdatabases
	where 'blomecj' is null or name = DB_NAME()))
BEGIN
	CREATE USER blomecj FROM LOGIN blomecj
	exec sp_addrolemember 'db_owner', blomecj;
END
GO
IF (NOT EXISTS (select name, suser_sname(sid)
	from master.dbo.sysdatabases
	where 'chalupmc' is null or name = DB_NAME()))
BEGIN
	CREATE USER chalupmc FROM LOGIN chalupmc
	exec sp_addrolemember 'db_owner', chalupmc;
END
GO
IF (NOT EXISTS (select name, suser_sname(sid)
	from master.dbo.sysdatabases
	where 'FelonAdmin' is null or name = DB_NAME()))
BEGIN
	CREATE USER FelonAdmin FROM LOGIN FelonAdmin
	exec sp_addrolemember 'db_datawriter', FelonAdmin;
	exec sp_addrolemember 'db_datareader', FelonAdmin;
	GRANT EXECUTE TO FelonAdmin;
END
GO