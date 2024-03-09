CREATE OR ALTER PROCEDURE [dbo].[Register]
	@Username nvarchar(30),
	@PasswordHash varchar(60)
AS
BEGIN
	if @Username is null or @Username = ''
	BEGIN
		Print 'Username cannot be null or empty.';
		RETURN (1)
	END
	if @PasswordHash is null or @PasswordHash = ''
	BEGIN
		Print 'PasswordHash cannot be null or empty.';
		RETURN (2)
	END
	IF (SELECT COUNT(*) FROM [User]
          WHERE Username = @Username) = 1
	BEGIN
      PRINT 'ERROR: Username already exists.';
	  RETURN(3)
	END
	INSERT INTO [User](Username, PasswordHash)
	VALUES (@username, @passwordHash)
END
GO
EXEC Register @Username='admin', @PasswordHash='$2a$10$oeRvKjF6ggv.tgLDvu/E3.G/wngMp3F4eacZ0P.vxTrDMPBvBO03O';