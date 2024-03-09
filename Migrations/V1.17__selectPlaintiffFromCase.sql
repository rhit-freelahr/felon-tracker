CREATE PROCEDURE selectPlaintiffFromCase(
	@ID INT
)
AS
BEGIN
	SET NOCOUNT ON

	IF EXISTS (SELECT * 
	FROM StatePlaintiff sp
	JOIN CasePlaintiff cp ON cp.PlaintiffID = sp.ID
	WHERE cp.CaseID = @ID)
	BEGIN
		SELECT sp.ID, sp.Title, '' AS [Name], '' AS [Address]
		FROM StatePlaintiff sp
		JOIN CasePlaintiff cp ON cp.PlaintiffID = sp.ID
		WHERE cp.CaseID = @ID
		RETURN
	END
	SELECT sp.ID, '' AS Title, [Name], [Address]
	FROM CivilianPlaintiff sp
	JOIN CasePlaintiff cp ON cp.PlaintiffID = sp.ID
	WHERE cp.CaseID = @ID
	RETURN
END