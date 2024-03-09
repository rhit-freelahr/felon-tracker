CREATE OR ALTER PROCEDURE insertCaseDefendant(
	@CaseID int,
	@DefendantID int
)
AS
BEGIN
	SET NOCOUNT ON
	IF(@CaseID IS NULL)
	BEGIN
		RAISERROR('CaseID cannot be null', 14, 3);
		RETURN 1;
	END
	IF(@DefendantID IS NULL)
	BEGIN
		RAISERROR('DefendantID cannot be null', 14, 3);
		RETURN 2;
	END
	IF(NOT EXISTS (SELECT * FROM Defendant WHERE ID = @DefendantID))
	BEGIN
		RAISERROR('DefendantID must exist', 14, 3);
		RETURN 3;
	END
	IF(NOT EXISTS (SELECT * FROM [Case] WHERE ID = @CaseID))
	BEGIN
		RAISERROR('CaseID must exist', 14, 3);
		RETURN 4;
	END
	

	INSERT INTO CaseDefendant (CaseID, DefendantID)
	VALUES(@CaseID, @DefendantID);
END