CREATE OR ALTER PROCEDURE updateUser
	--requires username and new password
	@Username nvarchar(30),
	@PassHash varchar(60)
AS
BEGIN
	IF NOT EXISTS (SELECT * FROM [User] WHERE Username = @Username)
	BEGIN
		RAISERROR('User does not exist', 14 , 6);
		RETURN 1;
	END
	IF @PassHash is null or @PassHash = ''
	BEGIN
		Print 'PasswordHash cannot be null or empty.';
		RETURN (2)
	END
	UPDATE [User]
	SET [PasswordHash] = @PassHash
	WHERE Username = @Username;
END