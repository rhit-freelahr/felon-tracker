CREATE OR ALTER PROCEDURE selectCasesFromDefendantName(
	@Name nvarchar(30)
)
AS
BEGIN
	IF(@Name IS NOT NULL)
	BEGIN
		SELECT ca.ID, ca.CaseNumber, ca.[Name] AS [CaseName], d.[Name] AS [DefendantName]
		FROM Defendant d
		JOIN CaseDefendant cd ON cd.DefendantID = d.ID
		JOIN [Case] ca ON ca.ID = cd.CaseID
		WHERE d.[Name] LIKE '%' + @Name + '%'
		ORDER BY ca.CaseNumber;
		RETURN 0;
	END
	ELSE
	BEGIN
		SELECT ca.ID, ca.CaseNumber, ca.[Name] AS [CaseName], d.[Name] AS [DefendantName]
		FROM Defendant d
		JOIN CaseDefendant cd ON cd.DefendantID = d.ID
		JOIN [Case] ca ON ca.ID = cd.CaseID
		ORDER BY ca.CaseNumber;
		RETURN 0;
	END
	
END