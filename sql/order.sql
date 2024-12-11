
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

-- @CREATE-ORDER
CREATE TABLE IF NOT EXISTS Order(
	ID       INTEGER PRIMARY KEY AUTOINCREMENT,
    Desc     VARCHAR(512),
	ClientID INTEGER,
    Cost     INTEGER NOT NULL,
    Prepay   INTEGER,
    Done     BOOLEAN DEFAULT false NOT NULl,

    FOREIGN KEY (ClientID) REFERENCES Client(ID)
);

-- @ADD-ORDER
INSERT INTO Order (Desc, ClientID, Cost) VALUES (?, ?, ?);

-- @ADD-PREPAYED-ORDER
INSERT INTO Order (Desc, ClientID, Cost, Prepay) VALUES (?, ?, ?, ?);

-- @FIND-ALL-CLIENT-ORDERS
SELECT Order.Id, Client.Name, Client.Phone From Order
INNER JOIN Client ON Order.ID = Client.ClientID
where
    Order.Done = false AND(
    Client.Name  = ? Or
    Client.Phone = ?);
