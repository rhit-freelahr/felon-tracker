CREATE PROCEDURE updateCharge(
	@ID int,
	@Name nvarchar(100) = NULL,
	@Degree varchar(8) = NULL,
	@Statute varchar(15) = NULL
)
AS
BEGIN
	SET NOCOUNT ON
    IF NOT EXISTS(SELECT * FROM Charge WHERE ID = @ID)
    BEGIN
        RAISERROR('Charge does not exist', 14 , 6);
		RETURN 1;
    END
	UPDATE Charge
	SET Name = ISNULL(@Name, Name),
		Degree = ISNULL(@Degree, Degree),
		Statute = ISNULL(@Statute, Statute)
	WHERE ID = @ID
END