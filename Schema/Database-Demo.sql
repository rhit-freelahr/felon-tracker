CREATE DATABASE [FelonTracker-Demo]
ON
PRIMARY(
	NAME = [FelonTracker-Demo],
	FILENAME = '/var/opt/mssql/data/FelonTracker-Demo.mdf'
)
LOG ON(
	NAME = [FelonTracker-Demo_log],
	FILENAME = '/var/opt/mssql/data/FelonTracker-Demo_log.ldf'
)
GO

USE [FelonTracker-Demo]
GO


CREATE USER freelahr FROM LOGIN freelahr
	exec sp_addrolemember 'db_owner', freelahr;
GO
	CREATE USER chalupmc FROM LOGIN chalupmc
	exec sp_addrolemember 'db_owner', chalupmc;
GO
	CREATE USER FelonAdmin FROM LOGIN FelonAdmin
	exec sp_addrolemember 'db_datawriter', FelonAdmin;
	exec sp_addrolemember 'db_datareader', FelonAdmin;
	GRANT EXECUTE TO FelonAdmin;
GO