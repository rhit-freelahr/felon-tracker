CREATE PROCEDURE updateEvent(
	@ID int,
	@Title varchar(100) = NULL,
	@DateFiled date = NULL,
	@Description varchar(200) = NULL,
	@CaseID int = NULL
)
AS
BEGIN
	SET NOCOUNT ON
	IF NOT EXISTS(SELECT * FROM [Event] WHERE ID = @ID)
    BEGIN
        RAISERROR('Event does not exist', 14 , 6);
		RETURN 1;
    END
	IF NOT EXISTS (SELECT * FROM [Case] WHERE ID = @CaseID)
	BEGIN
		RAISERROR('Case ID does not exist', 14, 6)
		RETURN
	END
	UPDATE [Event]
	SET Title = ISNULL(@Title, Title),
		DateFiled = ISNULL(@DateFiled, DateFiled),
		Description = ISNULL(@Description, Description),
		CaseID = ISNULL(@CaseID, CaseID)
	WHERE ID = @ID
END