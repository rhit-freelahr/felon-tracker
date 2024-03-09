CREATE PROCEDURE selectCase (
	@CaseNumber char(20) = NULL
)
AS
BEGIN
	SET NOCOUNT ON
	IF(@CaseNumber IS NULL)
	BEGIN
		SELECT * FROM [Case]
		RETURN
	END
	SELECT * 
	FROM [Case]
	WHERE CaseNumber = @CaseNumber
END