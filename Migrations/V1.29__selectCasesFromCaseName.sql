CREATE OR ALTER PROCEDURE selectCasesFromCaseName (
	@Name nvarchar(30)
)
AS
BEGIN
	IF(@Name IS NOT NULL)
	BEGIN
		SELECT ID, CaseNumber, [Name], [Type], DateFiled
		FROM [Case]
		WHERE [Name] LIKE '%' + @Name +'%'
		ORDER BY CaseNumber;
		RETURN 0;
	END
	ELSE
	BEGIN
		SELECT ID, CaseNumber, [Name], [Type], DateFiled
		FROM [Case]
		ORDER BY CaseNumber;
		RETURN 0;
	END
END