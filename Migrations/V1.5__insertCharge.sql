CREATE OR ALTER PROCEDURE insertCharge(
@Name nvarchar(100),
@Degree varchar(8),
@Statute varchar(15),
@ID int OUT
)
AS
BEGIN 
	SET NOCOUNT ON
	IF(@Name IS NULL OR @Name = '')
	BEGIN
		RAISERROR('Name cannot be null or empty', 14, 5);
		RETURN 1;
	END
	IF(@Degree IS NULL OR @Degree = '')
	BEGIN
		RAISERROR('Degree cannot be null', 14, 5);
		RETURN 2;
	END
	IF(@Statute IS NULL OR @Statute = '')
	BEGIN
		RAISERROR('Statute cannot be null', 14, 5);
		RETURN 3;
	END

	DECLARE @Temp INT
	SELECT @Temp = ID 
	FROM [Charge] 
	WHERE Degree = @Degree AND Statute = @Statute
	IF (@Temp IS NOT NULL)
	BEGIN
		SET @ID = @Temp
		RETURN 0;
	END
	
	INSERT INTO Charge([Name], Degree, Statute)
	VALUES (@Name, @Degree, @Statute);

	SET @ID = @@Identity;

	RETURN 0;
END