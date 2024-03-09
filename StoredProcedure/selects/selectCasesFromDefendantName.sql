CREATE OR ALTER PROCEDURE selectCasesFromDefendantName(
	@Name nvarchar(30)
)
AS
BEGIN
	SELECT ca.CaseNumber, ca.[Name] AS [CaseName], d.[Name] AS [DefendantName], ca.DateFiled
	FROM Defendant d
	JOIN CaseDefendant cd ON cd.DefendantID = d.ID
	JOIN [Case] ca ON ca.ID = cd.CaseID
	WHERE d.[Name] LIKE '%' + @Name + '%'
	
END