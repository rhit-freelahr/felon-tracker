
CREATE NONCLUSTERED INDEX CaseNumber_IND
ON [Case](CaseNumber)
WITH(FILLFACTOR = 75);

CREATE NONCLUSTERED INDEX CaseName_IND
ON [Case]([Name]) 
WITH(FILLFACTOR = 75);

CREATE NONCLUSTERED INDEX DefendantName_IND
ON Defendant([Name])
WITH(FILLFACTOR = 75);

CREATE NONCLUSTERED INDEX Name_IND
ON Charge([Name])
WITH(FILLFACTOR = 75);