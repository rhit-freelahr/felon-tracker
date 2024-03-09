CREATE OR ALTER PROCEDURE insertCase(
	@CaseNumber char(20),
	@Status varchar(10),
	@Name nvarchar(30),
	@Type varchar(35),
	@Date date,
	@CourtID int,
	@ID int OUT
)
AS
BEGIN
	SET NOCOUNT ON
	--needs error handling/Data validation
	IF (@CaseNumber IS NULL OR @CaseNumber = '')
	BEGIN
		RAISERROR('CaseNumber cannot be null or empty', 14, 1);
		RETURN 1;
	END

	IF (NOT (@CaseNumber LIKE '_____-____-__-______'))
	BEGIN
		RAISERROR('caseNumber is in the wrong format', 14, 1);
		RETURN 8;
	END

	IF (EXISTS(SELECT * FROM [Case] WHERE CaseNumber = @CaseNumber))
	BEGIN
		RAISERROR('caseNumber is already in the database', 14, 1);
		RETURN 9;
	END

	IF(@Date IS NULL)
	BEGIN
		RAISERROR('Date cannot be null', 14, 1);
		RETURN 2;
	END

	IF(@Status IS NULL OR @Status='')
	BEGIN
		RAISERROR('Status cannot be null or empty', 14, 1);
		RETURN 3;
	END

	IF(@Name IS NULL OR @Name='')
	BEGIN
		RAISERROR('Name cannot be null or empty', 14, 1);
		RETURN 4;
	END

	IF(@Type IS NULL OR @Type='')
	BEGIN
		RAISERROR('Type cannot be null or empty', 14, 1);
		RETURN 5;
	END

	IF(@CourtID IS NULL OR @CourtID='')
	BEGIN
		RAISERROR('CourtID cannot be null or empty', 14, 1);
		RETURN 6;
	END

	IF (NOT EXISTS(SELECT * FROM Court WHERE ID = @CourtID))
	BEGIN
		RAISERROR('CourtID does not exist',14,1);
		RETURN 7;
	END

	DECLARE @Temp INT
	SELECT @Temp = ID 
	FROM [Case] 
	WHERE CaseNumber = @CaseNumber
	IF (@Temp IS NOT NULL)
	BEGIN
		SET @ID = @Temp
		RETURN 0;
	END

	--basically the idea behind creating Cases
	INSERT INTO [Case]
	VALUES(@CaseNumber, @Status, @Name, @Type, @Date, @CourtID);

	SET @ID=@@Identity;
	
	RETURN 0;
	
END