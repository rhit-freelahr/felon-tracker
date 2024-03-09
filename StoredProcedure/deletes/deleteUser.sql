CREATE OR ALTER PROCEDURE deleteUser(
	--requires username
	@Username nvarchar(30)
)
AS
BEGIN
	IF NOT EXISTS (SELECT * FROM [User] WHERE Username = @Username)
	BEGIN
		RAISERROR('User does not exist', 14 , 6);
		RETURN 1;
	END
	DELETE FROM [User]
	WHERE Username = @Username;
END