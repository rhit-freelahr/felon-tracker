CREATE PROCEDURE selectDefendantFromCase(
	@ID INT
)
AS
BEGIN
	SET NOCOUNT ON
	SELECT d.ID, d.[Address], d.[Name], d.[Description]
	FROM [Case] c
	JOIN CaseDefendant cd ON cd.CaseID = c.ID
	JOIN Defendant d on d.ID = cd.DefendantID
	WHERE c.ID = @ID
	RETURN
END