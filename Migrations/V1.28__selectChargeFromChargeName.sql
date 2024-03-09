CREATE OR ALTER PROCEDURE selectChargeFromChargeName(
	@Name nvarchar(30)
)
AS
BEGIN
	IF(@Name is NOT NULL)
	BEGIN
		SELECT * 
		FROM Charge
		WHERE [Name] LIKE '%' + @Name + '%'
		RETURN 0;
	END
	ELSE
	BEGIN
		SELECT *
		FROM Charge
		RETURN 0;
	END

END