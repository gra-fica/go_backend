
-- @CREATE-CLIENT
CREATE TABLE IF NOT EXISTS Client(
	ID       INTEGER PRIMARY KEY AUTOINCREMENT,
	Name     VARCHAR(40),
	Phone    VARCHAR(40)
);

-- @CREATE-ORDER
CREATE TABLE IF NOT EXISTS Order(
	ID       INTEGER PRIMARY KEY AUTOINCREMENT,
    Desc     VARCHAR(512),
	ClientID INTEGER,
    Prepay   INTEGER,
    Cost     INTEGER NOT NULL,
    Done     BOOLEAN DEFAULT false,

    FOREIGN KEY (ClientID) REFERENCES Client(ID)
);

-- @FIND-ALL-CLIENT-ORDERS
SELECT Order.Id, Client.Name, Client.Phone From Order
INNER JOIN Client ON Order.ID = Client.ClientID
where
    Order.Done = false AND(
    Client.Name  = ? Or
    Client.Phone = ?)
;
