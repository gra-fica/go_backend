-- @ADD-PRODUCT
INSERT INTO Product (Name, Price) VALUES ({}, {});

-- @ADD-ALIAS
INSERT INTO ALIAS (Alias, ProductID) VALUES ({}, {});

-- @ADD-CATEGORY
INSERT INTO Product (Name) VALUES ({}, {});

-- @ADD-CATEGORY-PRODUCT
INSERT INTO Category_Product (ProductID, ProductCategoryID) VALUES ({}, {});

-- @ADD-SALE
INSERT INTO Sale (ProductID, Quantity, Price) VALUES ({}, {}, {});

-- @ADD-TICKET
INSERT INTO Ticket (Date) VALUES ({});

-- @ADD-TICKET-SALE
INSERT INTO Ticket_Sale (TicketID, SaleID) VALUES ({}, {});

