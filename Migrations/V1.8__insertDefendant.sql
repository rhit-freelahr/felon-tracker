CREATE OR ALTER PROCEDURE insertDefendant(
	@Address varchar(50) = NULL,
	@Name nvarchar(30),
	@Description varchar(100) = NULL,
	@ID int OUT
)
AS
BEGIN
	SET NOCOUNT ON
	IF(@Name IS NULL or @Name='')
	BEGIN
		RAISERROR('Name cannot be null or empty', 14, 8);
		RETURN 1;
	END
	--if a plaintiff address can be null I'm just gonna assume defendant address can be null
	--Description can be null
	
	
	
	INSERT INTO Defendant([Address], [Name], [Description])
	VALUES(@Address, @Name, @Description);

	SET @ID = @@Identity;
	RETURN 0;
END