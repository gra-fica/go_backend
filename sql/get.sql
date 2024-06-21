-- everything query here should return only one result

-- @GET-PRODUCT-ID
SELECT * FROM Product WHERE ID == ?;

-- @GET-PRODUCT-NAME
SELECT * FROM Product WHERE name == ?;

-- @GET-ALIAS-ID
SELECT * FROM Alias WHERE ID == ?;
-- @GET-CATEGORY-ID
SELECT * FROM Category WHERE ID == ?;
-- @GET-CATEGORY-PRODUCT
SELECT * FROM Category_Product WHERE ProductID == ? AND ProductCategoryID == ?;
-- @GET-SALE-ID
SELECT * FROM Sale WHERE ID == ?;
-- @GET-TICKET-ID
SELECT * FROM Ticket WHERE ID == ?;
-- @GET-TICKET-SALE-ID
SELECT * FROM Ticket_Sale WHERE TicketID == ? AND SaleID == ?;
