CREATE OR ALTER PROCEDURE selectCaseCharges(
	@CaseID int
)
AS
BEGIN
	SELECT *
	FROM CaseCharge cc
	WHERE cc.CaseID = @CaseID;
	RETURN 0;
	
END