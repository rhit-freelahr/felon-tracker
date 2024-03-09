CREATE PROCEDURE selectChargeFromDefendant(
	@ID int
)
AS
BEGIN
	SET NOCOUNT ON
	
	SELECT c.ID, c.[Name], c.Degree, c.Statute
	FROM Defendant d
	JOIN CaseDefendant cd ON cd.DefendantID = d.ID
	JOIN CaseCharge cc ON cc.CaseID = cd.CaseID
	JOIN Charge c ON c.ID = cc.ChargeID
	WHERE d.ID = @ID
	RETURN
END
