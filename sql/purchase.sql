
-- @CREATE-CLIENT
CREATE TABLE IF NOT EXISTS Client(
	ID       INTEGER PRIMARY KEY AUTOINCREMENT,
	Name     VARCHAR(40),
	Phone    VARCHAR(40)
);

-- @ADD-CLIENT-NAME
INSERT INTO Client (Name) VALUES (?);

-- @ADD-CLIENT-PHONE
INSERT INTO Client (Phone) VALUES (?);

-- @ADD-CLIENT-NAME-PHONE
INSERT INTO Client (Name, Phone) VALUES (?, ?);

-- @SEARCH-CLIENT-NAME
SELECT * FROM Client WHERE Name == ? ;

-- @SEARCH-CLIENT-PHONE
SELECT * FROM Client WHERE Phone== ?;

-- @SEARCH-CLIENT-NAME-PHONE
SELECT * FROM Client WHERE Phone== ? and Name == ?;

-- @CREATE-PURCHASE
CREATE TABLE IF NOT EXISTS Purchase(
	ID       INTEGER PRIMARY KEY AUTOINCREMENT,
    Desc     VARCHAR(512) NOT NULL,
    Cost     INTEGER NOT NULL,
    Done     BOOLEAN DEFAULT false,

	ClientID INTEGER,
    Prepay   INTEGER,

    FOREIGN KEY (ClientID) REFERENCES Client(ID)
);

-- @ADD-PURCHASE
INSERT INTO Purchase (Desc, ClientID, Cost) VALUES (?, ?, ?);

-- @ADD-PREPAYED-PURCHASE
INSERT INTO Purchase (Desc, ClientID, Cost, Prepay) VALUES (?, ?, ?, ?);

-- @GET-ALL-CLIENT-PURCHASES
SELECT Purchase.Id, Purchase.Desc, Client.Name, Client.Phone From
    Purchase JOIN Client ON
        Client.ID = Purchase.ClientID 
where
    Purchase.Done = false AND (
    Client.Name  = ? Or
    Client.Phone = ?);
