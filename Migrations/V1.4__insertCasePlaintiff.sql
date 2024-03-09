CREATE OR ALTER PROCEDURE insertCasePlaintiff(
	@CaseID int,
	@PlaintiffID int
)
AS
BEGIN
	SET NOCOUNT ON
	IF(@CaseID IS NULL)
	BEGIN
		RAISERROR('CaseID cannot be null', 14, 4);
		RETURN 1;
	END
	IF(@PlaintiffID IS NULL)
	BEGIN
		RAISERROR('PlaintiffID cannot be null', 14, 4);
		RETURN 2;
	END
	IF(NOT EXISTS (SELECT * FROM Plaintiff WHERE ID = @PlaintiffID))
	BEGIN
		RAISERROR('PlaintiffID must exist', 14, 4);
		RETURN 3;
	END
	IF(NOT EXISTS (SELECT * FROM [Case] WHERE ID = @CaseID))
	BEGIN
		RAISERROR('CaseID must exist', 14, 4);
		RETURN 4;
	END



	INSERT INTO CasePlaintiff (CaseID, PlaintiffID)
	VALUES(@CaseID, @PlaintiffID);
END