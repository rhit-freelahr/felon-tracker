CREATE PROCEDURE updateCourt(
	@ID int,
	@Name  nvarchar(30) = NULL,
	@Judge nvarchar(30) = NULL
)
AS
BEGIN
	SET NOCOUNT ON
    IF NOT EXISTS(SELECT * FROM Court WHERE ID = @ID)
    BEGIN
        RAISERROR('Court does not exist', 14 , 6);
		RETURN 1;
    END
	UPDATE Court
	SET Name = ISNULL(@Name, Name),
		Judge = ISNULL(@Judge, Judge)
	WHERE ID = @ID
END