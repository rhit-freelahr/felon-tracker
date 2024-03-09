CREATE OR ALTER PROCEDURE insertCaseCharge(
	@CaseID int,
	@ChargeID int,
	@Date date
)
AS
BEGIN
	SET NOCOUNT ON
	IF(@CaseID IS NULL)
	BEGIN
		RAISERROR('CaseID cannot be null', 14, 2);
		RETURN 1;
	END
	IF(@ChargeID IS NULL)
	BEGIN
		RAISERROR('ChargeID cannot be null', 14, 2);
		RETURN 2;
	END
	IF(@Date IS NULL)
	BEGIN
		RAISERROR('Date cannot be null', 14, 2);
		RETURN 3;
	END
	IF(NOT EXISTS (SELECT * FROM Charge WHERE ID = @ChargeID))
	BEGIN
		RAISERROR('ChargeID must exist', 14, 2);
		RETURN 4;
	END
	IF(NOT EXISTS (SELECT * FROM [Case] WHERE ID = @CaseID))
	BEGIN
		RAISERROR('CaseID must exist', 14, 2);
		RETURN 5;
	END
	IF(EXISTS (SELECT * FROM [CaseCharge] WHERE CaseID=@CaseID AND ChargeID = @ChargeID))
	BEGIN
		RAISERROR('This CaseCharge Already Exists', 14, 2);
		RETURN 6;
	END

	INSERT INTO CaseCharge (CaseID, ChargeID, [Date])
	VALUES(@CaseID, @ChargeID, @Date);
END