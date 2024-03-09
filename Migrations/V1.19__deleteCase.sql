CREATE OR ALTER PROCEDURE deleteCase(
	@CaseNumber char(21)
)
AS
BEGIN 
	IF(NOT EXISTS(SELECT * FROM [Case] WHERE CaseNumber = @CaseNumber) OR @CaseNumber = NULL)
	BEGIN
		RAISERROR('CaseNumber does not exist', 14, 5);
		RETURN 1;
	END

	DELETE FROM [Case]
	WHERE CaseNumber = @CaseNumber;
	--Cascade should also nuke any connections this has with anything

END