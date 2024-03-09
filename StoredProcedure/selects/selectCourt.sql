CREATE PROCEDURE selectCourt(
	@ID int = NULL,
	@Name  nvarchar(30) = NULL,
	@Judge nvarchar(30) = NULL 
)
AS
BEGIN
	SET NOCOUNT ON
	IF(@ID IS NOT NULL)
	BEGIN
		SELECT * FROM Court WHERE ID = @ID
		RETURN
	END
	IF(@Name IS NULL AND @Judge IS NULL)
	BEGIN
		SELECT * FROM Court
		RETURN
	END
	IF(@Name IS NULL AND @Judge IS NOT NULL)
	BEGIN
		SELECT * FROM Court WHERE Judge = @Judge
		RETURN
	END
	IF(@Name IS NOT NULL AND @Judge IS NULL)
	BEGIN
		SELECT * FROM Court WHERE [Name] = @Name
		RETURN
	END
	SELECT * FROM Court WHERE Judge = @Judge AND [Name] = @Name
	RETURN
END