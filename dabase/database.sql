-- Create the Payments database
CREATE DATABASE Payments;
GO

USE Payments;
GO

-- Create the table to store payment records
CREATE TABLE Payments (
    PaymentID INT IDENTITY(1,1) PRIMARY KEY,
    PaymentDate DATE NOT NULL,
    Amount DECIMAL(18, 2) NOT NULL,
    PaymentMethod VARCHAR(50) NOT NULL,
    Status VARCHAR(50) NOT NULL,
    CustomerID INT NOT NULL,
    TransactionID VARCHAR(50) NOT NULL
);

-- Create the Customers table
CREATE TABLE Customers (
    CustomerID INT IDENTITY(1,1) PRIMARY KEY,
    CustomerName VARCHAR(100) NOT NULL,
    Email VARCHAR(100) NOT NULL,
    Phone VARCHAR(15) NOT NULL,
    Address VARCHAR(255) NOT NULL
);

-- Insert sample customer data
INSERT INTO Customers (CustomerName, Email, Phone, Address)
VALUES ('John Doe', 'john.doe@example.com', '555-1234', '123 Elm Street'),
       ('Jane Smith', 'jane.smith@example.com', '555-5678', '456 Oak Street');

-- Generate 237 payment records
DECLARE @i INT = 1;
WHILE @i <= 237
BEGIN
    INSERT INTO Payments (PaymentDate, Amount, PaymentMethod, Status, CustomerID, TransactionID)
    VALUES (GETDATE(), ROUND(RAND() * 100 + 1, 2), 'Credit Card', 'Completed', (SELECT TOP 1 CustomerID FROM Customers ORDER BY NEWID()), NEWID());
    
    SET @i = @i + 1;
END;
GO
