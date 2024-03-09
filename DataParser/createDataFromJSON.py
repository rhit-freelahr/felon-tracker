import pyodbc
import json
import os

current_directory = os.getcwd() + '\DataParser'

SERVER = "golem.csse.rose-hulman.edu"
DATABASE = "FelonTracker-Demo"
USERNAME = "FelonAdmin"
PASSWORD = "Password123"

connectionString = f'DRIVER={{ODBC Driver 18 for SQL Server}};SERVER={SERVER};DATABASE={DATABASE};UID={USERNAME};PWD={PASSWORD};TrustServerCertificate=yes'
conn = pyodbc.connect(connectionString) 
cursor = conn.cursor()


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

def main():
    path = os.path.join(current_directory, 'db.json')

    with open(path, 'r') as f:
        data = json.load(f)
    
    for d in data:
        court = d["Court"]
        case = d["Case"]
        events = d["Event"]
        defendant = d["Defendant"]
        plaintiff = d["Plaintiff"]
        charges = d["Charge"]
        courtID = writeCourt(court["Name"], court["Judge"])

        caseID = writeCase(case["CaseNumber"], case["Status"], case["Name"], case["Type"], case["DateFiled"], courtID)

        for event in events:
            writeEvent(event["Title"], event["DateFiled"], event["Description"], caseID)
        
        defendantID = writeDefendant(defendant["Address"], defendant["Name"], defendant["Description"])
        writeCaseDefendant(caseID, defendantID)
        
        if "Title" in plaintiff:
            plaintiffID = writeStatePlaintiff(plaintiff["Title"])
        else:
            plaintiffID = writeCivilianPlaintiff(plaintiff["Name"], plaintiff["Address"])

        writeCasePlaintiff(caseID, plaintiffID)

        for charge in charges:
            chargeID = writeCharge(charge["Name"], charge["Degree"], charge["Statute"])
            writeCaseCharge(caseID, chargeID, charge["Date"])
    
main()