CREATE PROCEDURE [dbo].[selectPlaintiffFromID](
	@ID INT
)
AS
BEGIN
	SET NOCOUNT ON

	IF EXISTS(SELECT * FROM StatePlaintiff WHERE ID = @ID)
	BEGIN
		SELECT * , '' AS [Name], '' AS [Address]
		FROM StatePlaintiff
		WHERE ID = @ID
	END
	IF EXISTS(SELECT * FROM CivilianPlaintiff WHERE ID = @ID)
	BEGIN
		SELECT ID, '' AS [Ttle], Name, Address
		FROM CivilianPlaintiff
		WHERE ID = @ID
	END
	RETURN
END