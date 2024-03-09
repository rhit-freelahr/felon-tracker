from selenium import webdriver
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.by import By
from selenium.common.exceptions import *

from selenium.webdriver.chrome.service import Service
import pyodbc
import os
import json

current_directory = os.getcwd() + '\DataParser'

SERVER = "golem.csse.rose-hulman.edu"
DATABASE = "FelonTracker-Demo"
USERNAME = "FelonAdmin"
PASSWORD = "Password123"

connectionString = f'DRIVER={{ODBC Driver 18 for SQL Server}};SERVER={SERVER};DATABASE={DATABASE};UID={USERNAME};PWD={PASSWORD};TrustServerCertificate=yes'
conn = pyodbc.connect(connectionString) 
cursor = conn.cursor()

# website definitions
url = 'https://public.courts.in.gov/mycase/#/vw/CaseSummary/eyJ2Ijp7IkNhc2VUb2tlbiI6IlhZUzJOSENfTVVPcHNnLXVudlFCWHpEMEswLTVNa193dFZDdXQyUGJkT1kxIn19'
path = "C:/Users/freelahr/.cache/selenium/chromedriver/win64/120.0.6099.109/chromedriver.exe" # path to your driver

service = Service(executable_path = path)
options = webdriver.ChromeOptions()
options.add_argument('--headless')
driver = webdriver.Chrome(service = service)
driver.get(url)
delay = 10

def writeCase(Number, Status, Name, Type, DateFiled, CourtID):
    try:
        sql = """
            DECLARE @out INT

            EXEC insertCase ?, ?, ?, ?, ?, ?, @out OUTPUT

            SELECT @out
        """
        cursor.execute(sql, (Number, Status, Name, Type, DateFiled, CourtID))
        id = cursor.fetchval()
        conn.commit()
        return id
    except:
        sql = """
            EXEC selectCase ?
        """
        cursor.execute(sql, (Number))
        id = cursor.fetchval()
        conn.commit()
        return id

def writeCharge(Name, Degree, Statute):
    sql = """
        DECLARE @out INT

        EXEC insertCharge ?, ?, ?, @out OUTPUT

        SELECT @out
    """
    cursor.execute(sql, (Name, Degree, Statute))
    id = cursor.fetchval()
    conn.commit()
    return id

def writeDefendant(Address, Name, Description):
    sql = """
        DECLARE @out INT

        EXEC insertDefendant ?, ?, ?, @out OUTPUT

        SELECT @out
    """
    cursor.execute(sql, (Address, Name, Description))
    id = cursor.fetchval()
    conn.commit()
    return id

def writeEvent(Title, DateFiled, Description, CaseID):
    sql = """
        DECLARE @out INT

        EXEC insertEvent ?, ?, ?, ?, @out OUTPUT

        SELECT @out
    """
    cursor.execute(sql, (Title, DateFiled, Description, CaseID))
    id = cursor.fetchval()
    conn.commit()
    return id

def writeCourt(Name, Judge):
    sql = """
        DECLARE @out INT

        EXEC insertCourt ?, ?, @out OUTPUT

        SELECT @out
    """
    cursor.execute(sql, (Name, Judge))
    id = cursor.fetchval()
    conn.commit()
    return id

def writeStatePlaintiff(Title): 
    sql = """
        DECLARE @out INT

        EXEC insertStatePlaintiff ?, @out OUTPUT

        SELECT @out
    """
    cursor.execute(sql, (Title))
    id = cursor.fetchval()
    conn.commit()
    return id

def writeCivilianPlaintiff(Name, Address):
    sql = """
        DECLARE @out INT

        EXEC insertCivilianPlaintiff ?, ?, @out OUTPUT

        SELECT @out
    """
    cursor.execute(sql, (Name, Address))
    id = cursor.fetchval()
    conn.commit()
    return id

def writeCaseDefendant(cID, dID):
    sql = """
            EXEC insertCaseDefendant ?, ?
        """
    cursor.execute(sql, (cID, dID))
    conn.commit()

def writeCasePlaintiff(cID, pID):
    sql = """
            EXEC insertCasePlaintiff ?, ?
        """
    cursor.execute(sql, (cID, pID))
    conn.commit()

def writeCaseCharge(cID, chID, Date):
    sql = """
        EXEC insertCaseCharge ?, ?, ?
    """
    cursor.execute(sql, (cID, chID, Date))
    conn.commit()

def makeEvent():
    Events = driver.find_element(By.CLASS_NAME, "event-list")
    EventTable = Events.find_element(By.TAG_NAME, "tbody")
    EventRows = EventTable.find_elements(By.XPATH, "./*")
    arr = []
    for row in EventRows:
        EventTitle = row.find_element(By.CLASS_NAME, "bold").get_attribute("innerHTML")
        EventDate = row.find_element(By.CSS_SELECTOR, "span[data-bind='html: event.EventDate']").get_attribute("innerHTML")
        EventDescription = None
        try:
            EventDescription = row.find_element(By.CSS_SELECTOR, "span[data-bind='html: event.CaseEvent.Comment']").get_attribute("innerHTML")
        except NoSuchElementException:
            pass
        if EventDescription == "":
            EventDescription = None
        dict = {
            "Title": EventTitle,
            "DateFiled": EventDate,
            "Description": EventDescription
        }
        arr.append(dict)
    return arr

def makeCase():
    CaseName = driver.find_element(By.XPATH, '/html/body/div[2]/div[8]/div/div[2]/div/h4').text
    CaseStatus = driver.find_element(By.CSS_SELECTOR, "span[data-bind='html: CaseStatus']").text
    CaseNumber = driver.find_element(By.CSS_SELECTOR, "span[data-bind='html: CaseNumber']").text
    CaseType = driver.find_element(By.CSS_SELECTOR, "span[data-bind='html: CaseType']").text
    DateFiled = driver.find_element(By.CSS_SELECTOR, "span[data-bind='html: FileDate']").text

    dict = {
        "CaseNumber": CaseNumber,
        "Status": CaseStatus,
        "Name": CaseName,
        "Type": CaseType,
        "DateFiled": DateFiled
    }
    return dict

def makeCivilianPlaintiff():
    PlaintiffName = driver.find_elements(By.CLASS_NAME,"ccs-party-c3")[1]
    Plaintiff = PlaintiffName.find_element(By.TAG_NAME,"span").text
    CivilianDetail = driver.find_elements(By.CLASS_NAME,"ccs-party-detail-row")[1]
    try:
        CivilianPlaintiffAddress = CivilianDetail.find_element(By.CSS_SELECTOR, "span[aria-labelledby='labelPartyAttyAddr']").get_attribute("innerText")
    except NoSuchElementException:
        pass
    try:
        CivilianPlaintiffAddress = CivilianDetail.find_element(By.CSS_SELECTOR, "span[aria-labelledby='labelPartyAddr']").get_attribute("innerText")
    except NoSuchElementException:
        CivilianPlaintiffAddress = None
        pass
    dict = {
        "Name": Plaintiff,
        "Address": CivilianPlaintiffAddress
    }
    return dict

def makeStatePlaintiff():
    PlaintiffName = driver.find_elements(By.CLASS_NAME,"ccs-party-c3")[1]
    Plaintiff = PlaintiffName.find_element(By.TAG_NAME,"span").text
    dict = {
        "Title": Plaintiff
    }
    return dict

def makeCharge():
    try:
        ChargeList = driver.find_elements(By.CLASS_NAME, "ccs-charge-row")
        ChargeDetail = driver.find_elements(By.CLASS_NAME, "ccs-charge-detail-row")
        arr = []
        for i in range(len(ChargeList)):
            ChargeName = ChargeList[i].find_element(By.CSS_SELECTOR, "span[data-bind='html: charge.OffenseDescription']").get_attribute('innerHTML')
            ChargeDegree = ChargeDetail[i].find_element(By.CSS_SELECTOR, "span[data-bind='html: charge.OffenseDegree']").get_attribute('innerHTML')
            ChargeStatute = ChargeDetail[i].find_element(By.CSS_SELECTOR, "span[data-bind='html: charge.OffenseStatute']").get_attribute('innerHTML')
            ChargeDate = ChargeList[i].find_element(By.CSS_SELECTOR, "span[data-bind='html: charge.OffenseDate']").get_attribute('innerHTML')
            dict = {
                "Name": ChargeName,
                "Degree": ChargeDegree,
                "Statute": ChargeStatute,
                "Date": ChargeDate      # for CaseCharge table
            }
            arr.append(dict)
    except NoSuchElementException:
        pass
    return arr

def makeDefendant():
    DefendantClass = driver.find_elements(By.CLASS_NAME, "ccs-party-c3")[0]
    DefendantDetail = driver.find_elements(By.CLASS_NAME,"ccs-party-detail-row")[0]
    DefendantName = DefendantClass.find_element(By.CSS_SELECTOR, "span").text
    DefendantAddress = DefendantDetail.find_element(By.CSS_SELECTOR, "span[aria-labelledby='labelPartyAddr']").get_attribute("innerText")

    DefendantDescription = None
    try:
        DefendantDescription = driver.find_element(By.CSS_SELECTOR, "span[aria-labelledby='labelPartyDesc']").get_attribute("innerText")
    except NoSuchElementException:
        pass
    dict = {
        "Address": DefendantAddress,
        "Name": DefendantName,
        "Description": DefendantDescription
    }
    return dict

def makeCourt():
    CourtName = driver.find_element(By.CSS_SELECTOR, "span[aria-labelledby='lblCourt']").get_attribute('innerHTML')

    try:
        JudgeName = driver.find_element(By.CSS_SELECTOR, "span[data-bind='html: $parent.Judge']").get_attribute("innerHTML")
    except NoSuchElementException:
        pass
    try:
        JudgeName = driver.find_element(By.CSS_SELECTOR, "span[data-bind='html: event.Judge']").get_attribute("innerHTML")
    except NoSuchElementException:
        JudgeName = None
        pass
    
    dict = {
        "Name": CourtName,
        "Judge": JudgeName
    }
    return dict

def main():
    WebDriverWait(driver, delay).until(EC.presence_of_element_located((By.XPATH, '/html/body/div[2]/div[8]/div/div[2]/div/h4')))

    PlaintiffClass = driver.find_elements(By.CLASS_NAME,"ccs-party-c2")[1]
    PlaintiffType = PlaintiffClass.find_element(By.TAG_NAME,"span").text

    if(PlaintiffType == "State Plaintiff"):
        plaintiff = makeStatePlaintiff()
    else:
        plaintiff = makeCivilianPlaintiff()
    events = makeEvent()
    case = makeCase()
    charges = makeCharge()
    defendant = makeDefendant()
    court = makeCourt()

    courtID = writeCourt(court["Name"], court["Judge"])

    caseID = writeCase(case["CaseNumber"], case["Status"], case["Name"], case["Type"], case["DateFiled"], courtID)

    for charge in charges:
        chargeID = writeCharge(charge["Name"], charge["Degree"], charge["Statute"])
        writeCaseCharge(caseID, chargeID, charge["Date"])
        conn.commit()

    defendantID = writeDefendant(defendant["Address"], defendant["Name"], defendant["Description"])
    writeCaseDefendant(caseID, defendantID)

    for event in events:
        writeEvent(event["Title"], event["DateFiled"], event["Description"], caseID)

    if(PlaintiffType == "State Plaintiff"):
        plaintiffID = writeStatePlaintiff(plaintiff["Title"])
    else:
        plaintiffID = writeCivilianPlaintiff(plaintiff["Name"], plaintiff["Address"])
    
    writeCasePlaintiff(caseID, plaintiffID)

    cursor.close()
    conn.close()

    dict = [
        {
            "Case": case,
            "Defendant": defendant,
            "Plaintiff": plaintiff,
            "Charge": charges,
            "Event": events,
            "Court": court
        }
    ]
    path = os.path.join(current_directory, 'db.json')
    with open(path, 'r') as f:
        prev = json.load(f)
        dict = prev + dict

    with open(path, 'w') as f:
        json.dump(dict, f, indent=4)

main()