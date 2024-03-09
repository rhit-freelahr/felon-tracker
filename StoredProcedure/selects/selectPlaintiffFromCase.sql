CREATE PROCEDURE selectPlaintiffFromCase(
	@ID INT
)
AS
BEGIN
	SET NOCOUNT ON

	IF EXISTS (
		SELECT * 
		FROM StatePlaintiff s
		JOIN CasePlaintiff cp ON cp.PlaintiffID = s.ID
		WHERE caseID = @ID)
	BEGIN
		SELECT *, '' AS [Name], '' AS [Address]
		FROM StatePlaintiff s
		JOIN CasePlaintiff cp ON cp.PlaintiffID = s.ID
		WHERE caseID = @ID
		RETURN
	END
	SELECT ID, '' AS Title, [Name], [Address]
	FROM CivilianPlaintiff c
	JOIN CasePlaintiff cp ON cp.PlaintiffID = c.ID
	WHERE cp.CaseID = @ID
	RETURN
	
END