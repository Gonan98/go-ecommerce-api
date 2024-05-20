CREATE TABLE IF NOT EXISTS products (
    id INT NOT NULL AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    brand VARCHAR(255) NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    stock INT NOT NULL,
    PRIMARY KEY(id),
    UNIQUE KEY(name)
);