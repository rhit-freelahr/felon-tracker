CREATE OR ALTER PROCEDURE insertEvent(
@Title varchar(100),
@DateFiled date,
@Description varchar(200),
@CaseID int,
@ID int OUT
)
AS
BEGIN
	SET NOCOUNT ON
	IF(@Title IS NULL OR @Title='')
	BEGIN
		RAISERROR('Title cannot be null or empty', 14, 9);
		RETURN 1;
	END
	IF(@DateFiled IS NULL)
	BEGIN
		RAISERROR('DateFiled cannot be null', 14, 9);
		RETURN 2;
	END
	IF(@CaseID IS NULL)
	BEGIN
		RAISERROR('CaseID cannot be null', 14, 9);
		RETURN 3;
	END
	--Description can be null


	--Needs to create event
	INSERT INTO [Event] (Title, DateFiled, [Description], CaseID) --we need to add case number to EVENT Table
	VALUES(@Title, @DateFiled, @Description, @CaseID);



END
