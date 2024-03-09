package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	color "github.com/fatih/color"
	civil "github.com/golang-sql/civil"
	mssql "github.com/microsoft/go-mssqldb"
)

var (
	ErrOutBounds    = errors.New("Index Out of Bounds")
	ErrEmptyUName   = errors.New("UserName field is empty")
	ErrEmptyPWord   = errors.New("Password field is empty")
	ErrNoSuchUser   = errors.New("User doesn't exist")
	ErrRegistration = errors.New("Error Registering")
	ErrUserExists   = errors.New("User already exists")
	//insertCaseErrors
	ErrCaseNumberEmpty       = errors.New("CaseNumber Cannot be NULL or Empty")
	ErrCaseDateNull          = errors.New("Date cannot be null")
	ErrCaseStatusEmpty       = errors.New("Status cannot be null or empty")
	ErrCaseNameEmpty         = errors.New("Name cannot be null or empty")
	ErrCaseTypeEmpty         = errors.New("Type cannot be null or empty")
	ErrCaseCourtIDEmpty      = errors.New("CourtID cannot be null or empty")
	ErrCaseCourtIDNotExist   = errors.New("CourtID does not exist")
	ErrCaseNumberWrongFormat = errors.New("caseNumber is in the wrong format")
	ErrCaseNumberExists      = errors.New("Case already exists")
	//deleteCaseErrors/updateCaseErrors
	ErrCaseNumberNotExist = errors.New("CaseNumber does not exist")
	//insertEventErrors
	ErrEventTitleNull  = errors.New("Title cannot be null or empty")
	ErrEventDateNull   = errors.New("DateFiled cannot be null")
	ErrEventCaseIDNull = errors.New("CaseID cannot be null")
	//editDefendant
	ErrDefendantNotExist = errors.New("Defendant does not exist")
	//editPlaintiff
	ErrPlaintiffNotExist = errors.New("CivilianPlaintiff does not exist")
	//insertDefendant
	ErrDefendantNameNull = errors.New("Name cannot be null or empty")
	//insertCaseDefendant
	ErrCaseNull      = errors.New("CaseID cannot be null")
	ErrDefendantNull = errors.New("DefendantID cannot be null")
	// ErrDefendantNotExist = errors.New("DefendantID must exist") //already implemented
	ErrCaseIDNotExist = errors.New("CaseID must exist")

	//insertPlaintiff
	ErrPlaintiffNameNull = errors.New("Name cannot be null or empty")
	//insertCasePlaintiff
	// ErrCaseNull      = errors.New("CaseID cannot be null")
	ErrPlaintiffNull = errors.New("PlaintiffID cannot be null")
	// ErrPlaintiffNotExist = errors.New("PlaintiffID must exist") already implemented
	// ErrCaseIDNotExist = errors.New("CaseID must exist")

	//insertCharge
	ErrEmptyChargeName   = errors.New("Charge Name is empty or null")
	ErrNullChargeDegree  = errors.New("Charge Degree is null")
	ErrNullChargeStatute = errors.New("Charge Statute is null")
	//insertCaseCharge
	//ErrCaseNull
	ErrChargeNull     = errors.New("ChargeID cannot be null")
	ErrDateNull       = errors.New("Date cannot be null")
	ErrChargeNotExist = errors.New("ChargeID must exist")
	// ErrCaseIDNotExist
	ErrCaseChargeAlreadyExists = errors.New("This CaseCharge Already Exists")

	//TO REPLACE NO ROWS
	ErrNoResults = errors.New("No results found")
)

type Case struct {
	ID         int
	CaseNumber string
	Status     string
	Name       string
	Type       string
	DateFiled  string
	CourtID    *Court
}

func (c *Case) displayCase(indent int, displayEvents bool, db *sql.DB) {
	d, err := getDefendant(c.ID, db)
	if err != nil && err != sql.ErrNoRows {
		color.Red("error getting defendant: ", err)
		return
	}
	p, err := getPlaintiff(c.ID, db)
	if err != nil && err != sql.ErrNoRows {
		color.Red("error getting plaintiff")
		return
	}
	cs, err := getCaseCharges(c.ID, db)
	if err != nil && err != sql.ErrNoRows {
		color.Red("error getting charges")
		return
	}
	es, err := getEvents(c, db)
	if err != nil && err != sql.ErrNoRows {
		color.Red("error getting events")
		return
	}
	printIndents(indent)
	color.Cyan(c.Name + ":")
	printIndents(indent)
	color.Green("\tCase Number: " + c.CaseNumber)
	printIndents(indent)
	color.Yellow("\tStatus: " + c.Status)
	printIndents(indent)
	color.Red("\tType: " + c.Type)
	printIndents(indent)
	color.Magenta("\tDate Filed: " + c.DateFiled)
	printIndents(indent + 1)
	color.White("Court: ")
	c.CourtID.displayCourt(indent + 2)
	if len(d) != 0 {
		printIndents(indent + 1)
		color.White("Defendant: ")
		for _, de := range d {
			de.displayDefendant(indent + 2)
			fmt.Println()
		}
	}
	if len(cs) != 0 {
		printIndents(indent + 1)
		color.White("Charges: ")
		for _, ch := range cs {
			ch.displayCaseCharge(indent + 2)
			fmt.Println()
		}
	}
	if len(p) != 0 {
		printIndents(indent + 1)
		color.White("Plaintiff: ")
		for _, pl := range p {
			pl.displayPlaintiff(indent + 2)
			fmt.Println()
		}
	}
	if displayEvents {
		printIndents(indent + 1)
		color.White("Events: ")
		for _, e := range es {
			e.displayEvent(indent + 2)
		}
	}
}

func (c *Case) deleteCaseID(db *sql.DB) error {
	stmt, err := db.Prepare("deleteCase")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.ID)
	if err != nil {
		return err
	}
	return nil
}

func (c *Case) deleteCase(db *sql.DB) error {
	stmt, err := db.Prepare("deleteCase")
	if err != nil {
		return err
	}
	defer stmt.Close()
	var ret mssql.ReturnStatus
	_, err = stmt.Exec(c.CaseNumber, &ret)
	if ret == 1 {
		return ErrCaseNumberNotExist
	}
	if err != nil {
		return err
	}
	return nil
}

func (c *Case) updateCase(CaseNumber string, Status string, Name string, Type string, db *sql.DB) error {
	stmt, err := db.Prepare("updateCase")
	if err != nil {
		return err
	}
	defer stmt.Close()

	var SQLCaseNumber sql.NullString
	var SQLStatus sql.NullString
	var SQLName sql.NullString
	var SQLType sql.NullString

	SQLCaseNumber.String = CaseNumber
	SQLStatus.String = Status
	SQLName.String = Name
	SQLType.String = Type

	SQLCaseNumber.Valid, SQLStatus.Valid, SQLName.Valid, SQLType.Valid = true, true, true, true

	if len(SQLCaseNumber.String) == 0 {
		SQLCaseNumber.Valid = false
	}
	if len(SQLStatus.String) == 0 {
		SQLStatus.Valid = false
	}
	if len(SQLName.String) == 0 {
		SQLName.Valid = false
	}
	if len(SQLType.String) == 0 {
		SQLType.Valid = false
	}
	var ret mssql.ReturnStatus

	_, err = stmt.Exec(SQLCaseNumber, SQLStatus, SQLName, SQLType, &ret)
	if ret == 1 {
		return ErrCaseNumberNotExist
	}
	if err != nil {
		return err
	}

	return nil
}

func insertCase(CaseNumber string, Status string, Name string, Type string, DateFiled civil.Date, CourtID *Court, db *sql.DB) error {
	stmt, err := db.Prepare("insertCase")
	if err != nil {
		return err
	}
	defer stmt.Close()
	var ID int
	var ret mssql.ReturnStatus
	_, err = stmt.Exec(CaseNumber, Status, Name, Type, DateFiled, CourtID.ID, sql.Out{Dest: &ID}, &ret)
	if ret == 1 {
		return ErrCaseNumberEmpty
	}
	if ret == 2 {
		return ErrCaseDateNull
	}
	if ret == 3 {
		return ErrCaseStatusEmpty
	}
	if ret == 4 {
		return ErrCaseNameEmpty
	}
	if ret == 5 {
		return ErrCaseTypeEmpty
	}
	if ret == 6 {
		return ErrCaseCourtIDEmpty
	}
	if ret == 7 {
		return ErrCaseCourtIDNotExist
	}
	if ret == 8 {
		return ErrCaseNumberWrongFormat
	}
	if ret == 9 {
		return ErrCaseNumberExists
	}
	if err != nil {
		return err
	}
	return nil
}

func getCase(CaseNumber string, db *sql.DB) (*Case, error) {
	stmt, err := db.Prepare("selectCase")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row, err := stmt.Query(CaseNumber)
	if err != nil {
		return nil, err
	}
	if !row.Next() {
		return nil, ErrNoResults
	}
	courtCase, err := parseToCase(row, db)
	if err != nil {
		return nil, err
	}
	return courtCase, nil
}

func getCasebyID(CaseID int, db *sql.DB) (*Case, error) {
	stmt, err := db.Prepare("selectCaseFromID")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row, err := stmt.Query(CaseID)
	if err != nil {
		return nil, err
	}
	if !row.Next() {
		return nil, ErrNoResults
	}
	courtCase, err := parseToCase(row, db)
	if err != nil {
		return nil, err
	}
	return courtCase, nil
}

func getCases(db *sql.DB) ([]*Case, error) {
	stmt, err := db.Prepare("selectCase")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	cases := make([]*Case, 0)
	for rows.Next() {
		courtCase, err := parseToCase(rows, db)
		if err != nil {
			return nil, err
		}
		cases = append(cases, courtCase)
	}
	return cases, nil
}

func parseToCase(row *sql.Rows, db *sql.DB) (*Case, error) {
	courtCase := Case{}
	var courtID int
	row.Scan(&courtCase.ID, &courtCase.CaseNumber, &courtCase.Status, &courtCase.Name,
		&courtCase.Type, &courtCase.DateFiled, &courtID)
	court, err := getCourt(courtID, db)
	if err != nil {
		return nil, err
	}
	courtCase.CourtID = court
	return &courtCase, nil
}

type Court struct {
	ID    int
	Name  string
	Judge string
}

func (c *Court) displayCourt(indent int) {
	printIndents(indent)
	color.Cyan("Name: " + c.Name)
	if c.Judge != "" {
		printIndents(indent)
		color.Green("Judge: " + c.Judge)
	}
}

func getCourt(id int, db *sql.DB) (*Court, error) {
	stmt, err := db.Prepare("selectCourt")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row, err := stmt.Query(id, nil, nil)
	if err != nil {
		return nil, err
	}
	row.Next()
	return parseToCourt(row), nil
}

func getCourts(db *sql.DB) ([]*Court, error) {
	stmt, err := db.Prepare("selectCourt")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	courts := make([]*Court, 0)
	for rows.Next() {
		court := parseToCourt(rows)
		if err != nil {
			return nil, err
		}
		courts = append(courts, court)
	}
	return courts, nil
}

func parseToCourt(row *sql.Rows) *Court {
	court := Court{}
	row.Scan(&court.ID, &court.Name, &court.Judge)
	return &court
}

type Defendant struct {
	ID          int
	Address     string
	Name        string
	Description string
}

func (d *Defendant) displayDefendant(indent int) {
	printIndents(indent)
	color.Cyan("Name: " + d.Name)
	if d.Address != "" {
		printIndents(indent)
		color.Magenta("Address: " + d.Address)
	}
	if d.Description != "" {
		printIndents(indent)
		color.Yellow("Description: " + d.Description)
	}
}

func (d *Defendant) displayDefendantEdit(index int) {
	color.Cyan("[" + strconv.Itoa(index) + "]" + " Name: " + d.Name)
	if d.Address != "" {
		color.Magenta("Address: " + d.Address)
	}
	if d.Description != "" {
		color.Yellow("Description: " + d.Description)
	}
}

func getDefendant(caseID int, db *sql.DB) ([]*Defendant, error) {
	stmt, err := db.Prepare("selectDefendantFromCase")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(caseID)
	if err != nil {
		return nil, err
	}
	defendants := make([]*Defendant, 0)
	for rows.Next() {
		defendant := Defendant{}
		rows.Scan(&defendant.ID, &defendant.Address, &defendant.Name, &defendant.Description)
		defendants = append(defendants, &defendant)
	}
	return defendants, nil
}

func insertDefendant(caseID int, address string, name string, desc string, db *sql.DB) (int, error) {
	stmt, err := db.Prepare("insertDefendant")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	var ID int

	var ret mssql.ReturnStatus

	_, err = stmt.Exec(address, name, desc, sql.Named("ID", sql.Out{Dest: &ID}), &ret)
	if ret == 1 {
		return -1, ErrDefendantNameNull
	}
	if err != nil {
		return -1, err
	}
	stmt, err = db.Prepare("insertCaseDefendant")

	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	_, err = stmt.Exec(caseID, ID)
	if err != nil {
		return -1, err
	}
	return ID, nil
}

func getDefendantFromID(ID int, db *sql.DB) (*Defendant, error) {
	stmt, err := db.Prepare("selectDefendantFromID")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row, err := stmt.Query(ID)
	if err != nil {
		return nil, err
	}
	if !row.Next() {
		return nil, nil
	}

	return parseToDefendant(row), nil
}

func getDefendants(db *sql.DB) ([]*Defendant, error) {
	stmt, err := db.Prepare("selectDefendant")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defendants := make([]*Defendant, 0)
	for rows.Next() {
		defendant := Defendant{}
		rows.Scan(&defendant.ID, &defendant.Address, &defendant.Name, &defendant.Description)
		defendants = append(defendants, &defendant)
	}
	return defendants, nil
}

func parseToDefendant(row *sql.Rows) *Defendant {
	defendant := Defendant{}
	row.Scan(&defendant.ID, &defendant.Address, &defendant.Name, &defendant.Description)
	return &defendant
}

func (c *Defendant) updateDefendant(ID int32, Address string, Name string, Description string, db *sql.DB) error {
	stmt, err := db.Prepare("updateDefendant")
	if err != nil {
		return err
	}
	defer stmt.Close()

	var SQLAddress sql.NullString
	var SQLName sql.NullString
	var SQLDescription sql.NullString
	var SQLID sql.NullInt32

	SQLAddress.String = Address
	SQLName.String = Name
	SQLDescription.String = Description
	SQLID.Int32 = ID

	SQLAddress.Valid, SQLName.Valid, SQLDescription.Valid, SQLID.Valid = true, true, true, true

	if len(SQLAddress.String) == 0 {
		SQLAddress.Valid = false
	}
	if len(SQLName.String) == 0 {
		SQLName.Valid = false
	}
	if len(SQLDescription.String) == 0 {
		SQLDescription.Valid = false
	}
	if SQLID.Int32 < 0 {
		SQLID.Valid = false
	}

	var ret mssql.ReturnStatus
	_, err = stmt.Exec(SQLID, SQLAddress, SQLName, SQLDescription, &ret)
	if ret == 1 {
		return ErrDefendantNotExist
	}
	if err != nil {
		return err
	}

	return nil
}

type Event struct {
	ID          int
	Title       string
	DateFiled   string
	Description string
	CaseID      *Case
}

func (e *Event) displayEvent(indent int) {
	printIndents(indent)
	color.Cyan("Title: " + e.Title)
	if e.DateFiled != "" {
		printIndents(indent)
		color.Magenta("Date Filed: " + e.DateFiled)
	}
	if e.Description != "" {
		printIndents(indent)
		color.Yellow("Description: " + e.Description)
	}
}

func insertEvent(Title string, Desc string, DateFiled civil.Date, CaseID *Case, db *sql.DB) error {
	stmt, err := db.Prepare("insertEvent")
	if err != nil {
		return err
	}
	defer stmt.Close()
	var ID int
	var ret mssql.ReturnStatus
	_, err = stmt.Exec(Title, DateFiled, Desc, CaseID.ID, sql.Out{Dest: &ID}, &ret)
	if ret == 1 {
		return ErrEventTitleNull
	}
	if ret == 2 {
		return ErrEventDateNull
	}
	if ret == 3 {
		return ErrEventCaseIDNull
	}
	if err != nil {
		return err
	}
	return nil
}

func getEvents(caseID *Case, db *sql.DB) ([]*Event, error) {
	stmt, err := db.Prepare("selectEventsFromCase")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(caseID.ID)
	if err != nil {
		return nil, err
	}
	events := make([]*Event, 0)
	for rows.Next() {
		events = append(events, parseToEvent(rows, caseID))
	}
	return events, nil
}

func parseToEvent(row *sql.Rows, caseID *Case) *Event {
	event := Event{}
	var ID int
	row.Scan(&event.ID, &event.Title, &event.DateFiled, &event.Description, &ID)
	event.CaseID = caseID
	return &event
}

type Charge struct {
	ID      int
	Name    string
	Degree  string
	Statute string
}

func (c *Charge) displayCharge(indent int) {
	printIndents(indent)
	color.Cyan("Name: " + c.Name)
	if c.Degree != "" {
		printIndents(indent)
		color.Magenta("Degree: " + c.Degree)
	}
	if c.Statute != "" {
		printIndents(indent)
		color.Yellow("Statute: " + c.Statute)
	}
}

func parseToCharge(row *sql.Rows) *Charge {
	charge := Charge{}
	row.Scan(&charge.ID, &charge.Name, &charge.Degree, &charge.Statute)
	return &charge
}

func getCharge(id int, db *sql.DB) (*Charge, error) {
	stmt, err := db.Prepare("selectCharge")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	row.Next()
	return parseToCharge(row), nil
}

type CaseCharge struct {
	Case   *Case
	Charge *Charge
	Date   string
}

func getCaseCharges(CaseID int, db *sql.DB) ([]*CaseCharge, error) {
	stmt, err := db.Prepare("selectCaseCharges")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(CaseID)
	if err != nil {
		return nil, err
	}
	charges := make([]*CaseCharge, 0)
	for rows.Next() {
		charge, err := parseToCaseCharge(rows, db)
		if err != nil {
			return nil, err
		}
		charges = append(charges, charge)
	}
	return charges, nil
}

func parseToCaseCharge(row *sql.Rows, db *sql.DB) (*CaseCharge, error) { //TODO
	caseCharge := CaseCharge{}
	var caseID int
	var chargeID int
	row.Scan(&caseID, &chargeID, &caseCharge.Date)
	Case, err := getCasebyID(caseID, db)
	charge, err := getCharge(chargeID, db)
	if err != nil {
		return nil, err
	}
	caseCharge.Case = Case
	caseCharge.Charge = charge
	return &caseCharge, nil
}

func (c *CaseCharge) displayCaseCharge(indent int) {
	c.Charge.displayCharge(indent) //might indent weird we'll see
	printIndents(indent)
	color.Blue("Date: " + c.Date)
}

func getChargesFromChargeName(Name string, db *sql.DB) ([]*Charge, error) {
	stmt, err := db.Prepare("selectChargeFromChargeName")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var SQLName sql.NullString
	SQLName.String = Name
	SQLName.Valid = true
	if len(SQLName.String) == 0 {
		SQLName.Valid = false
	}

	rows, err := stmt.Query(SQLName)
	if err != nil {
		return nil, err
	}
	charges := make([]*Charge, 0)
	for rows.Next() {
		charge := Charge{} //this is parsing the rows into Charge structs
		rows.Scan(&charge.ID, &charge.Name, &charge.Degree, &charge.Statute)
		charges = append(charges, &charge)
	}
	return charges, nil
}

func getCasesFromChargeName(Name string, db *sql.DB) ([]*ChargeNameSearchResult, error) {
	stmt, err := db.Prepare("selectCasesFromChargeName")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var SQLChargeName sql.NullString
	SQLChargeName.String = Name
	SQLChargeName.Valid = true
	if len(SQLChargeName.String) == 0 {
		SQLChargeName.Valid = false
	}

	rows, err := stmt.Query(SQLChargeName)
	if err != nil {
		return nil, err
	}
	Results := make([]*ChargeNameSearchResult, 0)
	for rows.Next() {
		res := ChargeNameSearchResult{} //this is parsing the rows into Charge structs
		rows.Scan(&res.CaseID, &res.CaseNumber, &res.CaseName, &res.ChargeName)
		Results = append(Results, &res)
	}
	if len(Results) == 0 {
		return nil, ErrNoResults
	}
	return Results, nil
}

func getCasesFromDefendantName(Name string, db *sql.DB) ([]*DefendantNameSearchResult, error) {
	stmt, err := db.Prepare("selectCasesFromDefendantName")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var SQLChargeName sql.NullString
	SQLChargeName.String = Name
	SQLChargeName.Valid = true
	if len(SQLChargeName.String) == 0 {
		SQLChargeName.Valid = false
	}

	rows, err := stmt.Query(SQLChargeName)
	if err != nil {
		return nil, err
	}
	Results := make([]*DefendantNameSearchResult, 0)
	for rows.Next() {
		res := DefendantNameSearchResult{} //this is parsing the rows into Charge structs
		rows.Scan(&res.CaseID, &res.CaseNumber, &res.CaseName, &res.DefendantName)
		Results = append(Results, &res)
	}
	if len(Results) == 0 {
		return nil, ErrNoResults
	}

	return Results, nil
}

func getCasesFromCaseName(Name string, db *sql.DB) ([]*CaseNameSearchResult, error) {
	stmt, err := db.Prepare("selectCasesFromCaseName")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var SQLChargeName sql.NullString
	SQLChargeName.String = Name
	SQLChargeName.Valid = true
	if len(SQLChargeName.String) == 0 {
		SQLChargeName.Valid = false
	}

	rows, err := stmt.Query(SQLChargeName)
	if err != nil {
		return nil, err
	}
	Results := make([]*CaseNameSearchResult, 0)
	for rows.Next() {
		res := CaseNameSearchResult{} //this is parsing the rows into Charge structs
		rows.Scan(&res.CaseID, &res.CaseNumber, &res.CaseName, &res.CaseType, &res.CaseDate)
		Results = append(Results, &res)
	}
	if len(Results) == 0 {
		return nil, ErrNoResults
	}

	return Results, nil
}

func insertCharge(Name string, Degree string, Statute string, db *sql.DB) (int, error) {
	stmt, err := db.Prepare("insertCharge")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	var ID int
	var ret mssql.ReturnStatus
	_, err = stmt.Exec(Name, Degree, Statute, sql.Named("ID", sql.Out{Dest: &ID}), &ret)
	if ret == 1 {
		return -1, ErrEmptyChargeName
	}
	if ret == 2 {
		return -1, ErrNullChargeDegree
	}
	if ret == 3 {
		return -1, ErrNullChargeStatute
	}
	if err != nil {
		return -1, err
	}
	return ID, nil
}

func insertCaseCharge(CaseID int, ChargeID int, Date civil.Date, db *sql.DB) error { //CaseID and ChargeID (use Casenumber to get CaseID use charge name and select charge from a list)
	stmt, err := db.Prepare("insertCaseCharge")
	if err != nil {
		return err
	}
	defer stmt.Close()
	var ret mssql.ReturnStatus
	_, err = stmt.Exec(CaseID, ChargeID, Date, &ret)
	if ret == 1 {
		return ErrCaseNull
	}
	if ret == 2 {
		return ErrChargeNull
	}
	if ret == 3 {
		return ErrDateNull
	}
	if ret == 4 {
		return ErrChargeNotExist
	}
	if ret == 5 {
		return ErrCaseIDNotExist
	}
	if ret == 6 {
		return ErrCaseChargeAlreadyExists
	}

	if err != nil {
		return err
	}
	return nil
}

type Plaintiff struct {
	ID      int
	Title   string
	Name    string
	Address string
}

func (p *Plaintiff) displayPlaintiff(indent int) {
	printIndents(indent)
	if p.Name != "" {
		color.Cyan("Name: " + p.Name)
	}
	if p.Address != "" {
		printIndents(indent)
		color.Magenta("Address: " + p.Address)
	}
	if p.Title != "" {
		color.Yellow("Title: " + p.Title)
	}
}
func (p *Plaintiff) displayCivilianPlaintiff(index int) {

	if p.Name != "" {
		color.Cyan("[" + strconv.Itoa(index) + "]" + " Name: " + p.Name)
	}
	if p.Address != "" {
		color.Magenta("Address: " + p.Address)
	}
	if p.Title != "" {
		color.Yellow("Title: " + p.Title)
	}
}

func getPlaintiffFromID(ID int, db *sql.DB) (*Plaintiff, error) {
	stmt, err := db.Prepare("selectPlaintiffFromID")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row, err := stmt.Query(ID)
	if err != nil {
		return nil, err
	}
	if !row.Next() {
		return nil, ErrNoResults
	}
	plaintiff := Plaintiff{}
	row.Scan(&plaintiff.ID, &plaintiff.Title, &plaintiff.Name, &plaintiff.Address)
	return &plaintiff, nil
}

func getCivilianPlaintiffs(db *sql.DB) ([]*Plaintiff, error) {
	stmt, err := db.Prepare("selectCivilianPlaintiff")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	plaintiffs := make([]*Plaintiff, 0)
	for rows.Next() {
		plaintiff := Plaintiff{}
		rows.Scan(&plaintiff.ID, &plaintiff.Name, &plaintiff.Address)
		plaintiffs = append(plaintiffs, &plaintiff)
	}
	return plaintiffs, nil
}

func (c *Plaintiff) updatePlaintiff(ID int32, Name string, Address string, db *sql.DB) error {
	stmt, err := db.Prepare("updateCivilianPlaintiff")
	if err != nil {
		return err
	}
	defer stmt.Close()

	var SQLAddress sql.NullString
	var SQLName sql.NullString
	var SQLID sql.NullInt32

	SQLAddress.String = Address
	SQLName.String = Name
	SQLID.Int32 = ID

	SQLAddress.Valid, SQLName.Valid, SQLID.Valid = true, true, true

	if len(SQLAddress.String) == 0 {
		SQLAddress.Valid = false
	}
	if len(SQLName.String) == 0 {
		SQLName.Valid = false
	}
	if SQLID.Int32 < 0 {
		SQLID.Valid = false
	}

	var ret mssql.ReturnStatus
	_, err = stmt.Exec(SQLID, SQLName, SQLAddress, &ret)
	if ret == 1 {
		return ErrPlaintiffNotExist
	}
	if err != nil {
		return err
	}

	return nil
}

func insertPlaintiff(CaseID int, Title string, Name string, Address string, db *sql.DB) (int, error) {
	var ID int
	var ret mssql.ReturnStatus
	if Title != "" {
		stmt, err := db.Prepare("insertStatePlaintiff")
		if err != nil {
			return -1, err
		}
		defer stmt.Close()
		_, err = stmt.Exec(Title, sql.Named("ID", sql.Out{Dest: &ID}), &ret)
		if ret == 1 {
			return -1, ErrPlaintiffNameNull
		}
		if err != nil {
			return -1, err
		}
		stmt, err = db.Prepare("insertCasePlaintiff")
		if err != nil {
			return -1, err
		}
		defer stmt.Close()

		_, err = stmt.Exec(CaseID, ID, &ret)
		if ret == 1 {
			return -1, ErrCaseNull
		}
		if ret == 2 {
			return -1, ErrPlaintiffNull
		}
		if ret == 3 {
			return -1, ErrPlaintiffNotExist
		}
		if ret == 4 {
			return -1, ErrCaseIDNotExist
		}
		if err != nil {
			return -1, err
		}
	} else {
		stmt, err := db.Prepare("insertCivilianPlaintiff")
		if err != nil {
			return -1, err
		}
		defer stmt.Close()
		_, err = stmt.Exec(Name, Address, sql.Named("ID", sql.Out{Dest: &ID}), &ret)
		if ret == 1 {
			return -1, ErrPlaintiffNameNull
		}
		if err != nil {
			return -1, err
		}
		stmt, err = db.Prepare("insertCasePlaintiff")
		if err != nil {
			return -1, err
		}
		defer stmt.Close()

		_, err = stmt.Exec(CaseID, ID, &ret)
		if ret == 1 {
			return -1, ErrCaseNull
		}
		if ret == 2 {
			return -1, ErrPlaintiffNull
		}
		if ret == 3 {
			return -1, ErrPlaintiffNotExist
		}
		if ret == 4 {
			return -1, ErrCaseIDNotExist
		}
		if err != nil {
			return -1, err
		}
	}

	return ID, nil
}

func getPlaintiff(caseID int, db *sql.DB) ([]*Plaintiff, error) {
	stmt, err := db.Prepare("selectPlaintiffFromCase")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(caseID)
	if err != nil {
		return nil, err
	}
	plaintiffs := make([]*Plaintiff, 0)
	for rows.Next() {
		plaintiff := Plaintiff{}
		rows.Scan(&plaintiff.ID, &plaintiff.Title, &plaintiff.Name, &plaintiff.Address)
		plaintiffs = append(plaintiffs, &plaintiff)
	}
	return plaintiffs, nil
}

func parseToPlaintiff(row *sql.Row) *Plaintiff {
	plaintiff := Plaintiff{}
	row.Scan(&plaintiff.ID, &plaintiff.Title, &plaintiff.Name, &plaintiff.Address)
	return &plaintiff
}

type ChargeNameSearchResult struct {
	CaseID     int
	CaseNumber string
	CaseName   string
	ChargeName string
}

func parseToChargeNameSearchResult(row *sql.Row) *ChargeNameSearchResult {
	chargeSearchRes := ChargeNameSearchResult{}
	row.Scan(&chargeSearchRes.CaseID, &chargeSearchRes.CaseNumber, &chargeSearchRes.CaseName, &chargeSearchRes.ChargeName)
	return &chargeSearchRes
}

type CaseNameSearchResult struct { //could have just tried to make a case with using select star but Im here now
	CaseID     int
	CaseNumber string
	CaseName   string
	CaseType   string
	CaseDate   string //case used string so I'll just copy whatever it is using
}

func parseToCaseNameSearchResult(row *sql.Row) *CaseNameSearchResult {
	caseSearchRes := CaseNameSearchResult{}
	row.Scan(&caseSearchRes.CaseID, &caseSearchRes.CaseNumber, &caseSearchRes.CaseName, &caseSearchRes.CaseDate)
	return &caseSearchRes
}

type DefendantNameSearchResult struct {
	CaseID        int
	CaseNumber    string
	CaseName      string
	DefendantName string
}

func parseToDefendantNameSearchResult(row *sql.Row) *DefendantNameSearchResult {
	defendantSearchRes := DefendantNameSearchResult{}
	row.Scan(&defendantSearchRes.CaseID, &defendantSearchRes.CaseNumber, &defendantSearchRes.CaseName, &defendantSearchRes.DefendantName)
	return &defendantSearchRes
}

func printIndents(indent int) {
	for i := 0; i < indent; i++ {
		fmt.Print("\t")
	}
}
