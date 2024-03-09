USE FelonTracker
GO
DECLARE @Court_ID int;
DECLARE @Case_ID int;
DECLARE @Charge_ID int;
DECLARE @CivilianPlaintiff_ID int;
DECLARE @StatePlaintiff_ID int;
DECLARE @Defendant_ID int;
DECLARE @Event_ID int;

--From Court csv
--Allen Superior Court 6,"Bohdan, John C - MAG, II"
--@Name, @Judge
EXEC insertCourt @Name='Superior Court 6', @Judge='Bohdan, John C - MAG, II', @ID=@Court_ID OUTPUT;


--From Case csv
--45D09-2401-EV-000222,Pending,OFELIA FONSECA v. Brittany Smith,EV - Evictions (Small Claims Docket),01/18/2024 
--CaseNum, Status, Name of Case, Type of Case, DateFiled, CourtID
EXEC insertCase @CaseNumber='45D09-2401-EV-000222', @Status='Pending', @Name='OFELIA FONSECA v. Brittany Smith', @Type='EV - Evictions (Small Claims Docket)',  @Date='2024-01-18', @CourtID=@Court_ID, @ID=@Case_ID OUTPUT;


--From Charge csv
--35-43-4-2(a)/F6: Theft where value of property is between $750 &amp; $50k.,F6,35-43-4-2(a)
--Name, Degree, Statute
EXEC insertCharge @Name='Theft where value of property is between $750 &amp; $50k.', @Degree='F6', @Statute='35-43-4-2(a)', @ID=@Charge_ID OUTPUT ; 


--From CivillianPlaintiff csv
--"Avery, Tiffany","7899 Taft StreetMerrillville, IN 46410"
--Name, Address

EXEC insertCivilianPlaintiff @Name='Avery, Tiffany', @Address = '7899 Taft StreetMerrillville, IN 46410', @ID = @CivilianPlaintiff_ID OUTPUT;


--From StatePlaintiff csv
--State of Indiana
--Title

EXEC insertStatePlaintiff @Title='State of Indiana', @ID=@StatePlaintiff_ID OUTPUT;


--From Defendant csv
--"2191 S 900 WFARMLAND, IN 47340","DUNCAN, GABRIEL N","Male, White, 6' 2"", 220 lbs."
--Address, Name, Description

EXEC insertDefendant @Address='2191 S 900 WFARMLAND, IN 47340', @Name='DUNCAN, GABRIEL N', @Description='DUNCAN, GABRIEL N","Male, White, 6'' 2"", 220 lbs.', @ID=@Defendant_ID OUTPUT;


--From Event csv
--Service Issued,01/16/2024,Mailed summons SGray
--Title, DateFiled, Description
EXEC insertEvent @Title='Service Issued', @DateFiled='2024-01-16', @Description='Mailed summons SGray',@CaseID=@Case_ID, @ID=@Event_ID OUTPUT;


--USE data from previous creations
--

EXEC insertCasePlaintiff @CaseID = @Case_ID, @PlaintiffID = @CivilianPlaintiff_ID;


EXEC insertCaseDefendant @CaseID = @Case_ID, @DefendantID = @Defendant_ID;


EXEC insertCaseCharge @CaseID = @Case_ID, @ChargeID = @Charge_ID, @Date='2024-01-16';


