
-- @CREATE-PRODUCT
CREATE TABLE Product (
	ID INT PRIMARY KEY AUTO_INCREMENT,
	Name VARCHAR(256) NOT NULL,
	Price DECIMAL(10, 2) NOT NULL
);

-- @CREATE-ALIAS
CREATE TABLE ProductAlias (
	ID INT PRIMARY KEY AUTO_INCREMENT,
	Alias VARCHAR(256) NOT NULL,
	ProductID INT NOT NULL,
	FOREIGN KEY (ProductID) REFERENCES Product(ID)
);

-- @CREATE-CATEGORY
CREATE TABLE Category (
	ID INT PRIMARY KEY AUTO_INCREMENT,
	Name VARCHAR(256) NOT NULL
);

-- @CREATE-CATEGORY-PRODUCT
CREATE TABLE Category_Product (
	ProductID INT NOT NULL,
	ProductCategoryID INT NOT NULL,
	PRIMARY KEY (ProductID, ProductCategoryID),
	FOREIGN KEY (ProductID) REFERENCES Product(ID),
	FOREIGN KEY (ProductCategoryID) REFERENCES Category(ID)
);

-- @CREATE-SALE
CREATE TABLE Sale(
	ID INT PRIMARY KEY AUTO_INCREMENT,
	ProductID INT NOT NULL,
	Quantity INT NOT NULL,
	Price INT NOT NULL,

	FOREIGN KEY (ProductID) REFERENCES Product(ID)
);

-- @CREATE-TICKET
CREATE TABLE Ticket(
	ID INT PRIMARY KEY AUTO_INCREMENT,
	Date TIMESTAMP NOT NULL
);

-- @CREATE-TICKET-SALE
CREATE TABLE Ticket_Sale(
	TicketID INT NOT NULL,
	SaleID INT NOT NULL,
	FOREIGN KEY (TicketID) REFERENCES Ticket(ID),
	FOREIGN KEY (SaleID) REFERENCES Sale(ID)
);


-- @DROP-TABLES
DROP TABLE IF EXISTS Product;
DROP TABLE IF EXISTS ProductAlias;
DROP TABLE IF EXISTS Category;
DROP TABLE IF EXISTS Category_Product;
DROP TABLE IF EXISTS Sale;
DROP TABLE IF EXISTS Ticket;
DROP TABLE IF EXISTS Ticket_Sale;
