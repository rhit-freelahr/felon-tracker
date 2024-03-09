CREATE TABLE [Court] (
	ID int IDENTITY
	PRIMARY KEY,
	Name nvarchar(30) NOT NULL
	UNIQUE,
	Judge nvarchar(30)--nullable
)
CREATE TABLE [Case] (
	ID int IDENTITY
	PRIMARY KEY,
	CaseNumber char(20) UNIQUE NOT NULL,
	Status varchar(10),--nullable
	Name nvarchar(30),--nullable
	Type varchar(35) NOT NULL,
	DateFiled date NOT NULL,
	CourtID int REFERENCES Court(ID),
	
)

ALTER TABLE [Case]
ADD CHECK (CaseNumber LIKE '_____-____-__-______');

CREATE TABLE [Defendant] (
	ID int IDENTITY
	PRIMARY KEY,
	Address varchar(50), --nullable
	Name nvarchar(30) NOT NULL,
	Description varchar(100)--nullable
)
CREATE TABLE [Event] (
	ID int IDENTITY
	PRIMARY KEY,
	Title varchar(100),--nullable
	DateFiled date NOT NULL,
	Description varchar(200),--nullable
	CaseID int REFERENCES [Case](ID) ON DELETE CASCADE
)
CREATE TABLE [Charge] (
	ID int IDENTITY
	PRIMARY KEY,
	Name nvarchar(100),--nullable
	Degree varchar(8), --nullable
	Statute varchar(15),--nullable
	UNIQUE(Degree, Statute) --This still feels wrong
)
CREATE TABLE [Plaintiff] (
	ID int IDENTITY
	PRIMARY KEY
)
CREATE TABLE [StatePlaintiff] ( --This table will probably just have one entry and it will be the state of indiana
	ID int NOT NULL REFERENCES Plaintiff(ID)
	PRIMARY KEY,
	Title varchar(30) NOT NULL UNIQUE
)
CREATE TABLE [CivilianPlaintiff] (
	ID int NOT NULL REFERENCES Plaintiff(ID)
	PRIMARY KEY,
	Name nvarchar(30) NOT NULL,
	Address nvarchar(30)--nullable
)
--CREATE TABLE [CaseEvents] (
--	CaseNumber char(17) REFERENCES [Case](CaseNumber),
--	EventID int REFERENCES Event(ID),
--	PRIMARY KEY(CaseNumber, EventID)
--)
CREATE TABLE [CaseDefendant] (
	--CaseNumber char(17) REFERENCES [Case](CaseNumber),
	CaseID int REFERENCES [Case](ID) ON DELETE CASCADE,
	DefendantID int REFERENCES Defendant(ID),
	PRIMARY KEY(CaseID, DefendantID)
)
CREATE TABLE [CasePlaintiff] (
	--CaseNumber char(17) REFERENCES [Case](CaseNumber),
	CaseID int REFERENCES [Case](ID) ON DELETE CASCADE,
	PlaintiffID int REFERENCES Plaintiff(ID),
	PRIMARY KEY(CaseID, PlaintiffID)
)
--CREATE TABLE [CaseCourt] (
--	--CaseNumber char(17) REFERENCES [Case](CaseNumber),
--	CaseID int REFERENCES [Case](ID),
--	CourtID int REFERENCES Court(ID),
--	PRIMARY KEY(CaseNumber, CourtID)
--)
CREATE TABLE [CaseCharge] (
	--CaseNumber char(17) REFERENCES [Case](CaseNumber),
	CaseID int REFERENCES [Case](ID) ON DELETE CASCADE,
	ChargeID int REFERENCES Charge(ID),
	[Date] date NOT NULL,
	PRIMARY KEY(CaseID, ChargeID)
)

CREATE TABLE [User] (
	Username nvarchar(30) PRIMARY KEY,
	PasswordHash varchar(60)
)