CREATE PROCEDURE updateCivilianPlaintiff(
	@ID int,
	@Name nvarchar(30) = NULL,
	@Address nvarchar(30) = NULL
)
AS
BEGIN
	SET NOCOUNT ON
    IF NOT EXISTS(SELECT * FROM CivilianPlaintiff WHERE ID = @ID)
    BEGIN
        RAISERROR('CivilianPlaintiff does not exist', 14 , 6);
		RETURN 1;
    END
	UPDATE CivilianPlaintiff
	SET Name = ISNULL(@Name, Name),
		Address = ISNULL(@Address, Address)
	WHERE ID = @ID
END