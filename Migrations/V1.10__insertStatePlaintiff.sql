CREATE OR ALTER PROCEDURE insertStatePlaintiff(
	@Title varchar(30),
	@ID int OUT
)
AS
BEGIN
	SET NOCOUNT ON
	IF(@Title IS NULL OR @Title='')
	BEGIN
		RAISERROR('DateFiled cannot be null', 14, 10);
		RETURN 1;
	END

	DECLARE @Temp INT
	SELECT @Temp = ID 
	FROM StatePlaintiff 
	WHERE Title = @Title
	IF (@Temp IS NOT NULL)
	BEGIN
		SET @ID = @Temp
		RETURN 0;
	END

	INSERT INTO Plaintiff DEFAULT VALUES;

	SET @ID = @@Identity;

	INSERT INTO StatePlaintiff(ID, Title)
	VALUES(@ID, @Title);

	RETURN 0;


END