CREATE OR ALTER PROCEDURE updateCase(
	--Users/application should only have the CaseNumber
	@CaseNumber char(21),
	@Status varchar(10)=NULL,
	@Name nvarchar(30)=NULL,
	@Type varchar(35)=NULL
	--should a case be able to change courts or would that be filed under a different caseNumber
)
AS
BEGIN
	IF(NOT EXISTS (SELECT * FROM [Case] WHERE CaseNumber = @CaseNumber) OR @CaseNumber = NULL)
	BEGIN
		RAISERROR('Case does not exist', 14 , 6);
		RETURN 1;
	END

	UPDATE [Case]
	SET [Status] = ISNULL(@Status, [Status]),
		[Name] = ISNULL(@Name, [Name]),
		[Type] = ISNULL(@Type, [Type])
	WHERE CaseNumber = @CaseNumber;
END