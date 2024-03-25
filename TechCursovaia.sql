CREATE TABLE Category (
    ID serial PRIMARY KEY,
    Name text NOT NULL UNIQUE,
    Relation int NULL default(0),
    Is_Exist boolean NULL default(true)
);

INSERT INTO Category (Name, Relation)
VALUES ('Стиральные машины', 0),
       ('Стиральные машины с фронтальной загрузкой', 1),
       ('Компактные стиральные машины', 1),
       ('Пылесосы', 0),
       ('Пылесосы с контейнером', 4);

CREATE TABLE Product (
    ID serial PRIMARY KEY,
    Name text NOT NULL UNIQUE,
    Price numeric NOT NULL CONSTRAINT Positive_Price CHECK (Price::numeric >= 0.00),
    Amount int NOT NULL CONSTRAINT Positive_Amount CHECK (Amount >= 0),
    Discount int NULL default(0) CONSTRAINT Correct_Discount CHECK (Discount >= 0 and Discount <= 100),
    Image_Link text NULL default(''),
    Category_ID integer REFERENCES Category (ID),
    Is_Exist boolean NULL default(true)
);

INSERT INTO Product (Name, Price, Amount, Category_ID)
VALUES ('Фронтальная стиральная машина Korting KWM 42D1460', 44000.00, 30, 2),
       ('Пылесос с контейнером Scarlett SC-VC80C86', 3000.00, 100, 5),
       ('Стиральная машина Electrolux EWC 1350', 70000.00, 23, 3);

CREATE TABLE Characteristic (
    ID serial PRIMARY KEY,
    Name text UNIQUE NOT NULL,
    Type text NULL default('END'),
    Relation int NULL default(0),
    Is_Exist boolean NULL default(true)
);

INSERT INTO Characteristic (Name, Type, Relation)
VALUES ('Скорость отжима, об/мин', 'INT', 0),
       ('Защита от протечек', 'SET', 0),
       ('Полная', 'END', 2),
       ('Частичная', 'END', 2),
       ('Отсутствует', 'END', 2),
       ('Тип уборки', 'SET', 0),
       ('Сухая', 'END', 6),
       ('Влажная', 'END', 6);

INSERT INTO Characteristic (Name, Type, Relation)
VALUES ('Test', 'END', 2);

CREATE TABLE Set (
    ID serial PRIMARY KEY,
    Value numeric NULL default(0),
    Product_ID integer REFERENCES Product (ID),
    Characteristic_ID integer REFERENCES Characteristic (ID)
);

INSERT INTO Set (Product_ID, Characteristic_ID)
VALUES (1, 3),
       (2, 7),
       (3, 4);

INSERT INTO Set (Product_ID, Characteristic_ID, Value)
VALUES (1, 1, 1400),
       (3, 1, 1300);

CREATE TABLE Role (
    ID serial PRIMARY KEY,
    Name text NOT NULL UNIQUE,
    Is_Exist boolean NULL default(true)
);

INSERT INTO Role (Name)
VALUES ('Администратор'),
       ('Клиент'),
       ('Товарный менеджер');

CREATE TABLE "User" (
    ID serial PRIMARY KEY,
    Email text NOT NULL UNIQUE,
    Password text NOT NULL,
    Phone text NULL default(''),
    Last_Name text NULL default(''),
    Name text NULL default(''),
    Middle_Name text NULL default(''),
    Gender text NULL default(''),
    Role_ID integer REFERENCES Role (ID),
    Is_Exist boolean NULL default(true)
);

INSERT INTO "User" (Email, Password, Role_ID)
VALUES ('qwerty@mail.com', '$2a$10$A2GmgcX735aXaoPklA4PxuoBpr7Ld.jMLq9HlybvzuqxSOQCieRDy', 1),
       ('Almostreal@mail.com', '$2a$10$4JR2gIRGWrsY7XlzC7SFFuUR5BxhOTKevssDjVYvbubizRi5gPGJe', 3),
       ('User@mail.com', '$2a$10$IFzoe6wTzO08706XsHrum.juEY1FL1DYH6V2q5HkAtAz9GBsIU8/K', 2);

CREATE TABLE Discount_Card (
    ID integer PRIMARY KEY REFERENCES "User" (ID),
    Date date NULL default(now()),
    Discount integer NOT NULL CONSTRAINT Correct_Discount CHECK (Discount >= 0 and Discount <= 100)
);

INSERT INTO Discount_Card (ID, Discount)
VALUES (3, 0);

CREATE TABLE Cart_Position (
    ID serial PRIMARY KEY,
    Amount integer NULL default(1),
    Visibility boolean NULL default(True),
    Product_ID integer REFERENCES Product (ID),
    User_ID integer REFERENCES "User" (ID)
);

CREATE TABLE Status (
    ID serial PRIMARY KEY,
    Status text NOT NULL UNIQUE,
    Is_Exist boolean NULL default(true)
);

INSERT INTO Status (Status)
VALUES ('Составлен'), ('Согласован'), ('Скомплектован'), ('Доставляется'),
       ('Готов к выдаче'), ('Выдан'), ('Возврат'), ('Отменён');

CREATE TABLE "Order" (
    ID serial PRIMARY KEY,
    Date timestamp NULL default(now()),
    Address text NOT NULL,
    Status_ID integer REFERENCES Status (ID),
    User_ID integer REFERENCES "User" (ID)
);

CREATE TABLE Order_Position (
    ID serial PRIMARY KEY,
    Checkout_Price numeric NOT NULL,
    Amount integer NULL default(1),
    Order_ID integer REFERENCES "Order" (ID),
    Product_ID integer REFERENCES Product (ID)
);

drop table Order_Position, "Order", Status, Cart_Position, Discount_Card, "User", Role, Set, Characteristic, Product, Category;

select * from "User";

select * from Characteristic;

