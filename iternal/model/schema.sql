-- Создание базы данных
CREATE DATABASE IF NOT EXISTS password_manager;

-- Использование базы данных
USE password_manager;

CREATE TABLE Passwords (
    id INT AUTO_INCREMENT PRIMARY KEY,
    resource VARCHAR(100) NOT NULL,
    login VARCHAR(100) NOT NULL,
    password VARCHAR(100),
    description VARCHAR(100),
);

CREATE TABLE Reminders (
    FOREIGN KEY (id) REFERENCES Passwords(id),
    remind DATETIME
)