CREATE PROCEDURE selectPassHash(
	@Username nvarchar(30)
)
AS
BEGIN
	SET NOCOUNT ON

	SELECT PasswordHash FROM [User] WHERE Username = @Username
	RETURN
END