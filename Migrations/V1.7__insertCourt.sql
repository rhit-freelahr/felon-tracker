CREATE OR ALTER PROCEDURE insertCourt(
	@Name  nvarchar(30),
	@Judge nvarchar(30) = NULL, 
	@ID int OUT
)
AS
BEGIN
	SET NOCOUNT ON
	--Need to add error handling/data validation
	IF(@Name IS NULL OR @Name='')
	BEGIN
		RAISERROR('Name cannot be null or empty', 14, 7);
		RETURN 1;
	END
	--Judge can be null

	DECLARE @Temp INT
	SELECT @Temp = ID 
	FROM [Court] 
	WHERE [Name] = @Name
	IF (@Temp IS NOT NULL)
	BEGIN
		SET @ID = @Temp
		RETURN 0;
	END

	INSERT INTO Court([Name], Judge)
	VALUES(@Name, @Judge);

	SET @ID = @@Identity; 
	RETURN 0;
END