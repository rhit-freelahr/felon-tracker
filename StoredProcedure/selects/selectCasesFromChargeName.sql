CREATE OR ALTER PROCEDURE selectCasesFromChargeName(
	@Name nvarchar(30)
)
AS
BEGIN
	IF(@Name IS NOT NULL)
	BEGIN
		SELECT ca.CaseNumber, ca.[Name] AS [CaseName], c.[Name] AS [ChargeName]
		FROM Charge c
		JOIN CaseCharge cc ON cc.ChargeID = c.ID
		JOIN [Case] ca ON ca.ID = cc.CaseID
		WHERE c.[Name] LIKE '%' + @Name + '%'
		ORDER BY ca.CaseNumber;
		RETURN 0;
	END
	ELSE
	BEGIN
		SELECT ca.CaseNumber, ca.[Name] AS [CaseName], c.[Name] AS [ChargeName]
		FROM Charge c
		JOIN CaseCharge cc ON cc.ChargeID = c.ID
		JOIN [Case] ca ON ca.ID = cc.CaseID
		ORDER BY ca.CaseNumber;
		RETURN 0;
	END

END 