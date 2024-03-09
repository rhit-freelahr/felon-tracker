CREATE OR ALTER PROCEDURE insertCivilianPlaintiff(
	@Name nvarchar(30),
	@Address nvarchar(30) = NULL,
	@ID int OUT
)
AS
BEGIN
	SET NOCOUNT ON
	--Need to add data validation and error handling
	IF(@Name IS NULL OR @Name='')
	BEGIN
		RAISERROR('Name cannot be null or empty', 14, 6);
		RETURN 1;
	END
	--Address is allowed to be null I believe






	INSERT INTO Plaintiff DEFAULT VALUES;


	SET @ID = @@Identity;

	INSERT INTO CivilianPlaintiff(ID, [Name], [Address])
	VALUES(@ID, @Name, @Address);

	RETURN 0;
END