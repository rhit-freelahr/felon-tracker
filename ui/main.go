package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"os"

	color "github.com/fatih/color"
	civil "github.com/golang-sql/civil"
	mssql "github.com/microsoft/go-mssqldb"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	sqlurl := fmt.Sprintf("sqlserver://FelonAdmin:Password123@golem.csse.rose-hulman.edu?database=FelonTracker-Demo")

	// flags
	user := flag.String("user", "", "username for update/delete")
	pass := flag.String("pass", "", "password for update/delete")
	reg := flag.Bool("r", false, "register username/password")
	ver := flag.Bool("v", false, "print version information")
	caseNum := flag.String("cn", "", "case number to look up")
	caseInsert := flag.Bool("ac", false, "add case")
	deleteCase := flag.String("dc", "", "delete case [Case Number]")
	addEvents := flag.Bool("ae", false, "add events to case")
	displayEvents := flag.Bool("de", false, "display case events")
	updateCase := flag.Bool("uc", false, "update case")
	deleteUser := flag.Bool("du", false, "delete user")
	updateUser := flag.Bool("uu", false, "update password for user")
	updateDefendant := flag.Bool("ud", false, "update defendant")
	updatePlaintiff := flag.Bool("up", false, "update civilian plaintiff")
	insertDefendant := flag.Bool("ad", false, "insert defendant")
	insertPlaintiff := flag.Bool("ap", false, "insert plaintiff")
	chargeInsert := flag.Bool("ach", false, "add charge")
	addChargeToCase := flag.Bool("ach2c", false, "add charge to a case")
	searchByCharge := flag.Bool("sch", false, "search cases by charge name")
	searchByCaseName := flag.Bool("scn", false, "search cases by case name")
	searchByDefendant := flag.Bool("scd", false, "search cases by defendant name")

	flag.Parse()

	// print version info
	if *ver {
		fmt.Println("Version 10.0.9")
		return
	}

	admin := false

	// connect to db
	db, err := connectToDB(sqlurl)
	defer db.Close()
	if err != nil {
		color.Red("Error connecting to database\n")
		return
	}

	// login
	if *user != "" {
		err = loginUser(*user, *pass, db)
		if err != nil {
			color.Red(err.Error())
			fmt.Println()
			return
		}
		admin = true
	}

	// inform user of admin permissions
	if admin {
		color.Green("Login successful")
		color.Red("You are operating as an admin")
		if *reg {
			uname, err := registerUser(db)
			if err != nil {
				if err == ErrUserExists {
					color.Red("User already Exists")
					return
				}
				if err == ErrEmptyUName {
					color.Red("Username cannot be empty")
					return
				}
				if err == ErrEmptyPWord {
					color.Red("Password cannot be empty")
					return
				}
				color.Red("registration error")
				fmt.Println()
				return
			} else {
				color.Green("User: " + uname + " successfully register\n")
			}
			return
		}
		// update/delete needs to happen before login as to not cause login error
		if *deleteUser {

			err := userDelete(db)
			if err != nil {
				if err == ErrNoSuchUser {
					color.Red("User doesn't exist")
					return
				}
				color.Red("Could not delete user")
				return
			}
			color.Green("User deleted successfully")
			return
		} else if *updateUser {
			err := userUpdate(db)
			if err != nil {
				if err == ErrNoSuchUser {
					color.Red("User does not exist")
					return
				}
				if err == ErrEmptyPWord {
					color.Red("Password cannot be empty")
					return
				}
				if err == ErrEmptyUName {
					color.Red("Username cannot be empty")
					return
				}
				color.Red("Could not update user")
				return
			}
			color.Green("User updated successfully")
			return
		} else if *caseInsert {
			err := addCase(db)
			if err != nil {
				color.Red(err.Error())
				return
			}
			color.Green("Case added Successfully")
			return
		} else if *deleteCase != "" {
			cases, err := getCase(*deleteCase, db)
			if err != nil {
				color.Red("error: cannot retrieve case for deletion")
				return
			}
			err = cases.deleteCase(db)
			if err != nil {
				color.Red(err.Error())
			} else {
				color.Green("Successfully deleted the case")
			}
			return
		} else if *updateCase {
			err := editCase(db)
			if err != nil {
				color.Red("error: could not update case")
			} else {
				color.Green("Case updated Successfully")
			}
			return
		} else if *addEvents {
			err := addEventToCase(db)
			if err != nil {
				color.Red(err.Error())
				return
			} else {
				color.Green("Event added successfully")
				return
			}
		} else if *updateDefendant {
			err := editDefendant(db)
			if err != nil {
				color.Red(err.Error())
			}
			return
		} else if *updatePlaintiff {
			err := editPlaintiff(db)
			if err != nil {
				color.Red(err.Error())
			}
			return
		} else if *insertDefendant {
			err := addDefendant(db)
			if err != nil {
				color.Red(err.Error())
			}
			return
		} else if *insertPlaintiff {
			err := addPlaintiff(db)
			if err != nil {
				color.Red(err.Error())
			}
			return
		} else if *chargeInsert {
			err := addCharge(db)
			if err != nil {
				color.Red(err.Error())
			} else {
				color.Green("Charge inserted Successfully")
			}
			return
		} else if *addChargeToCase {
			err := addCaseCharge(db)
			if err != nil {
				color.Red(err.Error())
			} else {
				color.Green("Charge added Successfully")
			}
			return
		}
	} else {
		color.Yellow("You are operating as a user\n")
		if *caseInsert || *deleteUser || *updateUser || *deleteCase != "" || *updateCase || *addEvents || *updateDefendant || *insertPlaintiff || *chargeInsert || *addChargeToCase {
			color.Red("You cannot add a case as a user. Please login using -user=username and -pass=password to add a case")
			return
		}
	}
	// get a case out
	if *caseNum != "" {
		courtCase, err := getCase(*caseNum, db)
		if err != nil {
			if err == sql.ErrNoRows {
				color.Red("No case with number " + *caseNum)
				return
			} else {
				color.Red("error getting case")
				return
			}
		}
		courtCase.displayCase(0, *displayEvents, db)
	} else if *searchByCharge {
		err := searchCaseByChargeName(db)
		if err != nil {
			if err == ErrOutBounds {
				color.Red("ERROR: index provided out of range")
			}
			color.Red(err.Error())
		}
		return
	} else if *searchByCaseName {
		err := searchCaseByCaseName(db)
		if err != nil {
			if err == ErrOutBounds {
				color.Red("ERROR: index provided out of range")
			}
			color.Red(err.Error())
		}
		return
	} else if *searchByDefendant {
		err := searchCaseByDefendantName(db)
		if err != nil {
			if err == ErrOutBounds {
				color.Red("ERROR: index provided out of range")
			}
			color.Red(err.Error())
		}
		return
	} else {
		cases, err := getCases(db)
		if err != nil {
			color.Red("error getting case")
		}
		for _, c := range cases {
			c.displayCase(0, *displayEvents, db)
			fmt.Println()
		}
	}
	return
}

func connectToDB(url string) (*sql.DB, error) {
	db, err := sql.Open("sqlserver", url)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func registerUser(db *sql.DB) (string, error) { //username string, password string,
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("UserName ?")
	scanner.Scan()
	username := scanner.Text()
	if scanner.Err() != nil {
		return username, scanner.Err()
	}
	if len(username) == 0 {
		return username, ErrEmptyUName
	}
	fmt.Println("Password?")
	scanner.Scan()
	password := scanner.Text()
	if scanner.Err() != nil {
		return username, scanner.Err()
	}
	if len(password) == 0 {
		return username, ErrEmptyPWord
	}

	passhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return username, err
	}
	stmt, err := db.Prepare("Register")
	if err != nil {
		return username, err
	}
	defer stmt.Close()
	var ret mssql.ReturnStatus
	_, err = stmt.Exec(username, passhash, &ret)
	if err != nil {
		return username, err
	}
	if ret == 3 {
		return username, ErrUserExists
	}
	if ret != 0 {
		return username, ErrRegistration
	}

	return username, nil
}

func loginUser(username string, password string, db *sql.DB) error {
	stmt, err := db.Prepare("selectPassHash")
	var passhash []byte
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows := stmt.QueryRow(username)
	if err != nil {
		return err
	}
	err = rows.Scan(&passhash)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword(passhash, []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func userDelete(db *sql.DB) error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("UserName ?")
	scanner.Scan()
	username := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}
	if len(username) == 0 {
		return ErrEmptyUName
	}

	stmt, err := db.Prepare("deleteUser")
	if err != nil {
		return err
	}
	defer stmt.Close()
	var ret mssql.ReturnStatus
	stmt.Exec(username, &ret)

	if ret == 1 {
		return ErrNoSuchUser
	}
	return nil
}

func userUpdate(db *sql.DB) error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("UserName ?")
	scanner.Scan()
	username := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}
	if len(username) == 0 {
		return ErrEmptyUName
	}

	fmt.Println("Password?")
	scanner.Scan()
	password := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}

	if len(password) == 0 {
		return ErrEmptyPWord
	}

	passhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	stmt, err := db.Prepare("updateUser")
	if err != nil {
		return err
	}
	defer stmt.Close()
	var ret mssql.ReturnStatus

	stmt.Exec(username, passhash, &ret)

	if ret == 1 {
		return ErrNoSuchUser
	}

	return nil
}

func addCase(db *sql.DB) error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Case Number?\nFormat: XXXXX-XXXX-XX-XXXXXX")
	scanner.Scan()
	cn := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}
	fmt.Println("Status?")
	scanner.Scan()
	status := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}
	fmt.Println("Name?")
	scanner.Scan()
	name := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}
	fmt.Println("Case Type?")
	scanner.Scan()
	casetype := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}
	fmt.Println("Date Filed? ([Year]-[Two digit month]-[Two digit day]")
	scanner.Scan()
	date := scanner.Text()
	datefiled, err := civil.ParseDate(date)
	if err != nil {
		return err
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	// courtid
	courts, err := getCourts(db)
	fmt.Println("Select Court (enter index to select)")
	for i, c := range courts {
		fmt.Printf("[%d]: %s\n", i, c.Name)
	}
	var index int
	_, err = fmt.Scanln(&index)
	if index >= len(courts) || index < 0 {
		return ErrOutBounds
	}
	court := courts[index]
	if err != nil {
		return err
	}
	err = insertCase(cn, status, name, casetype, datefiled, court, db)
	if err != nil {
		return err
	}
	color.Green("New Case:")
	c, err := getCase(cn, db)
	c.displayCase(0, false, db)
	if err != nil {
		return err
	}
	return nil
}

func addEventToCase(db *sql.DB) error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Case Number?\nFormat: XXXXX-XXXX-XX-XXXXXX")
	scanner.Scan()
	cn := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}
	c, err := getCase(cn, db)
	if err != nil {
		return err
	}

	fmt.Println("Title?")
	scanner.Scan()
	title := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}
	fmt.Println("Description")
	scanner.Scan()
	desc := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}
	fmt.Println("Date Filed? ([Year]-[Two digit month]-[Two digit day]")
	scanner.Scan()
	date := scanner.Text()
	datefiled, err := civil.ParseDate(date)
	if err != nil {
		return err
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	err = insertEvent(title, desc, datefiled, c, db)
	if err != nil {
		return err
	}
	return nil
}

func editCase(db *sql.DB) error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Case Number?\nFormat: XXXXX-XXXX-XX-XXXXXX")
	scanner.Scan()
	cn := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}

	fmt.Println("Current Data For Case: ")
	c, err := getCase(cn, db)
	if err != nil {
		return err
	}
	c.displayCase(0, false, db) //display current values in case

	fmt.Println("Status?")
	scanner.Scan()
	status := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}
	fmt.Println("Name?")
	scanner.Scan()
	name := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}
	fmt.Println("Case Type?")
	scanner.Scan()
	casetype := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}
	err = c.updateCase(cn, status, name, casetype, db) //is there a reason we aren't doing err :=
	if err != nil {
		return err
	}

	fmt.Println("Successfully Updated Record: ")
	c, err = getCase(cn, db)
	c.displayCase(0, false, db) //display new values in case

	return nil
}

func addCharge(db *sql.DB) error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Charge Name?")
	scanner.Scan()
	n := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}
	fmt.Println("Degree?")
	scanner.Scan()
	degree := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}
	fmt.Println("Statute?")
	scanner.Scan()
	statute := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}

	ID, err := insertCharge(n, degree, statute, db)
	if err != nil {
		return err
	}
	color.Green("New Charge:")
	c, err := getCharge(ID, db)
	c.displayCharge(0)
	if err != nil {
		return err
	}
	return nil
}

func addCaseCharge(db *sql.DB) error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Case Number?\nFormat: XXXXX-XXXX-XX-XXXXXX")
	scanner.Scan()
	cn := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}

	c, err := getCase(cn, db)
	if err != nil {
		return err
	}

	fmt.Println("Search Charge Name?")
	scanner.Scan()
	name := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}

	charges, err := getChargesFromChargeName(name, db)
	fmt.Println("Select Charge (enter index to select)")
	for i, c := range charges {
		fmt.Printf("[%d]: %s\n", i, c.Name)
	}
	var index int
	_, err = fmt.Scanln(&index)
	if len(charges) == 0 {
		return ErrNoResults
	}
	// if index == -1 {
	// 	return nil
	// }
	if index >= len(charges) || index < 0 {
		return ErrOutBounds
	}

	charge := charges[index]
	if err != nil {
		return err
	}

	fmt.Println("Date of Charge? ([Year]-[Two digit month]-[Two digit day]")
	scanner.Scan()
	date := scanner.Text()
	datefiled, err := civil.ParseDate(date)
	if err != nil {
		return err
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	err = insertCaseCharge(c.ID, charge.ID, datefiled, db) //might be wrong date
	if err != nil {
		return err
	}
	color.Green("Case With Charge Added:")
	Case, err := getCase(cn, db)
	Case.displayCase(0, false, db)
	if err != nil {
		return err
	}
	return nil
}

func searchCaseByChargeName(db *sql.DB) error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Search Charge Name?")
	scanner.Scan()
	chargeName := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}

	results, err := getCasesFromChargeName(chargeName, db)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNoResults
		}
		return err
	}
	fmt.Println("Select Case to view (enter index to select) (enter -1 to exit search)")
	for i, c := range results {
		fmt.Printf("[%d]: %s\t%s\t%s\n", i, c.CaseName, c.CaseNumber, c.ChargeName)
	}
	var index int
	_, err = fmt.Scanln(&index)
	if index == -1 {
		return nil
	}
	if index >= len(results) || index < 0 {
		return ErrOutBounds
	}
	res := results[index]
	if err != nil {
		return err
	}
	Case, err := getCasebyID(res.CaseID, db)
	if err != nil {
		return err
	}
	Case.displayCase(0, false, db)
	return nil
}

func searchCaseByDefendantName(db *sql.DB) error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Search Defendant Name?")
	scanner.Scan()
	chargeName := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}

	results, err := getCasesFromDefendantName(chargeName, db)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNoResults
		}
		return err
	}
	fmt.Println("Select Case to view (enter index to select) (-1 to exit search)")
	for i, c := range results {
		fmt.Printf("[%d]: %s\t%s\t%s\n", i, c.CaseName, c.CaseNumber, c.DefendantName)
	}
	var index int
	_, err = fmt.Scanln(&index)
	if index == -1 {
		return nil
	}
	if index >= len(results) || index < 0 {
		return ErrOutBounds
	}
	res := results[index]
	if err != nil {
		return err
	}
	Case, err := getCasebyID(res.CaseID, db)
	if err != nil {
		return err
	}
	Case.displayCase(0, false, db)
	return nil
}

func searchCaseByCaseName(db *sql.DB) error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Search Case Name?")
	scanner.Scan()
	chargeName := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}

	results, err := getCasesFromCaseName(chargeName, db)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNoResults
		}
		return err
	}
	fmt.Println("Select Case to view (enter index to select) (-1 to exit search)")
	for i, c := range results {
		fmt.Printf("[%d]: %s\t%s\t%s\t%s\n", i, c.CaseName, c.CaseNumber, c.CaseType, c.CaseDate)
	}
	var index int
	_, err = fmt.Scanln(&index)
	if index == -1 {
		return nil
	}
	if index >= len(results) || index < 0 {
		return ErrOutBounds
	}
	res := results[index]
	if err != nil {
		return err
	}
	Case, err := getCasebyID(res.CaseID, db)
	if err != nil {
		return err
	}
	Case.displayCase(0, false, db)
	return nil
}

func editDefendant(db *sql.DB) error {
	scanner := bufio.NewScanner(os.Stdin)
	defendants, err := getDefendants(db)
	if err != nil {
		return err
	}
	for i, d := range defendants {
		d.displayDefendantEdit(i)
		fmt.Println()
	}
	var index int
	color.Yellow("Select Defendant to edit (enter index to select) (-1 to exit search)")
	_, err = fmt.Scanln(&index)
	if index == -1 {
		return nil
	}
	if index >= len(defendants) || index < 0 {
		return ErrOutBounds
	}
	id := defendants[index].ID
	if err != nil {
		return err
	}
	fmt.Println("Current For Defendant: ")
	d, err := getDefendantFromID(id, db)
	if err != nil {
		color.Red("Invalid Defendant")
		return err
	}
	d.displayDefendant(0)
	fmt.Println("Address?")
	scanner.Scan()
	address := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}
	fmt.Println("Name?")
	scanner.Scan()
	name := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}
	fmt.Println("Description?")
	scanner.Scan()
	description := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}
	err = d.updateDefendant(int32(id), address, name, description, db) //is there a reason we aren't doing err :=
	if err != nil {
		return err
	}
	d, err = getDefendantFromID(id, db)

	fmt.Println("Successfully Updated Record: ")
	d.displayDefendant(0)
	return nil
}

func editPlaintiff(db *sql.DB) error {
	scanner := bufio.NewScanner(os.Stdin)
	plaintiffs, err := getCivilianPlaintiffs(db)
	if err != nil {
		return err
	}
	for i, p := range plaintiffs {
		p.displayCivilianPlaintiff(i)
		fmt.Println()
	}
	var index int
	color.Yellow("Select Plaintiff to edit (enter index to select) (-1 to exit search)")
	_, err = fmt.Scanln(&index)
	if index == -1 {
		return nil
	}
	if index >= len(plaintiffs) || index < 0 {
		return ErrOutBounds
	}
	id := plaintiffs[index].ID
	if err != nil {
		return err
	}
	fmt.Println("Current For Plaintiff: ")
	p, err := getPlaintiffFromID(id, db)
	if err != nil {
		color.Red("Invalid Plaintiff")
		return err
	}
	p.displayPlaintiff(0)

	fmt.Println("Name?")
	scanner.Scan()
	name := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}
	fmt.Println("Address?")
	scanner.Scan()
	address := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}
	err = p.updatePlaintiff(int32(id), name, address, db) //is there a reason we aren't doing err :=
	if err != nil {
		return err
	}
	np, err := getPlaintiffFromID(id, db)
	if err != nil {
		return err
	}

	fmt.Println("Successfully Updated Record: ")
	np.displayPlaintiff(0)

	return nil
}

func addDefendant(db *sql.DB) error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Case Number?\nFormat: XXXXX-XXXX-XX-XXXXXX")
	scanner.Scan()
	cn := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}
	c, err := getCase(cn, db)

	if c == nil {
		color.Red("Invalid Case Number\n")
		return err
	}

	fmt.Println("Address?")
	scanner.Scan()
	address := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}
	fmt.Println("Name?")
	scanner.Scan()
	name := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}
	fmt.Println("Description?")
	scanner.Scan()
	desc := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}
	id, err := insertDefendant(c.ID, address, name, desc, db)
	if err != nil {
		return err
	}
	color.Green("New Defendant:")
	d, err := getDefendantFromID(id, db)
	d.displayDefendant(0)
	if err != nil {
		return err
	}
	color.Green("Defendant inserted Successfully")
	return nil
}

func addPlaintiff(db *sql.DB) error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Case Number?\nFormat: XXXXX-XXXX-XX-XXXXXX")
	scanner.Scan()
	cn := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}
	c, err := getCase(cn, db)
	if c == nil {
		color.Red("Invalid Case Number\n")
		return err
	}
	fmt.Println("StatePlaintiff: 1\nCivilianPlaintiff: 2")
	scanner.Scan()
	state := scanner.Text()
	if scanner.Err() != nil {
		return scanner.Err()
	}

	if state == "2" {
		fmt.Println("Name?")
		scanner.Scan()
		name := scanner.Text()
		if scanner.Err() != nil {
			return scanner.Err()
		}
		fmt.Println("Address?")
		scanner.Scan()
		address := scanner.Text()
		if scanner.Err() != nil {
			return scanner.Err()
		}
		id, err := insertPlaintiff(c.ID, "", name, address, db)
		if err != nil {
			return err
		}
		color.Green("New Plaintiff:")
		p, err := getPlaintiffFromID(id, db)
		p.displayPlaintiff(0)
		if err != nil {
			return err
		}
	} else if state == "1" {
		fmt.Println("Title?")
		scanner.Scan()
		title := scanner.Text()
		if scanner.Err() != nil {
			return scanner.Err()
		}
		id, err := insertPlaintiff(c.ID, title, "", "", db)
		if err != nil {
			return err
		}
		color.Green("New Plaintiff:")
		p, err := getPlaintiffFromID(id, db)
		p.displayPlaintiff(0)
		if err != nil {
			return err
		}
	} else {
		return nil
	}

	color.Green("Plaintiff inserted Successfully")
	return nil
}
