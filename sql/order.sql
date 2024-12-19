-- @ADD-ORDER-PHONE
INSERT INTO Order (PhoneNum, Fee) VALUES (?, ?);

-- @ADD-ORDER-NAME
INSERT INTO Order (Name, Fee) VALUES (?, ?);

-- @ADD-ORDER-PHONE-PAYMENT
INSERT INTO Order (PhoneNum, Fee, Payment) VALUES (?, ?, ?);

-- @ADD-ORDER-Name-PAYMENT
INSERT INTO Order (Name, Fee, Payment) VALUES (?, ?, ?);

-- @GET-ORDER-UNFINISHED
SELECT * FROM Order WHERE Done == false and (PhoneNum == ? or Client == ?);

-- @GET-ORDERS
SELECT * FROM Order WHERE PhoneNum == ? or Client == ?;

-- @ORDER-DONE-ID
UPDATE TABLE Order
    SET Done = true
    WHERE ID == ?;

