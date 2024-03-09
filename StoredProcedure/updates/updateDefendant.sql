CREATE PROCEDURE updateCharge(
	@ID int,
	@Address varchar(30) = NULL,
	@Name nvarchar(30) = NULL,
	@Description varchar(100) = NULL
)
AS
BEGIN
	SET NOCOUNT ON
    IF NOT EXISTS(SELECT * FROM Charge WHERE ID = @ID)
    BEGIN
        RAISERROR('Charge does not exist', 14 , 6);
		RETURN 1;
    END
	UPDATE Defendant
	SET Address = ISNULL(@Address, Address),
		Name = ISNULL(@Name, Name),
		Description = ISNULL(@Description, Description)
	WHERE ID = @ID
END