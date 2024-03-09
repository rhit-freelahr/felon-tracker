CREATE OR ALTER PROCEDURE deleteUser(
	--requires username
	@Username nvarchar(30)
)
AS
BEGIN
	IF ((@Username IS NULL) OR @Username='')
	BEGIN
		RAISERROR('User cannot be NULL or Empty', 14 , 6);
		RETURN 2;
	END

	IF NOT EXISTS (SELECT * FROM [User] WHERE Username = @Username)
	BEGIN
		RAISERROR('User does not exist', 14 , 6);
		RETURN 1;
	END
	DELETE FROM [User]
	WHERE Username = @Username;
END