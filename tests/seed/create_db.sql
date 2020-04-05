USE go_contacts_test;

CREATE TABLE `contact` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `first_name` varchar(50) NOT NULL,
  `last_name` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE `email` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `contact_id` INT NOT NULL, 
  `address` VARCHAR(320) NOT NULL,
  PRIMARY KEY(`id`),
  CONSTRAINT `fk_email_contact` FOREIGN KEY (`contact_id`)
    REFERENCES `contact`(`id`)
    ON DELETE CASCADE
);

CREATE TABLE `phone` (
  `id` int NOT NULL AUTO_INCREMENT,
  `contact_id` INT NOT NULL,
  `type` ENUM('mobile', 'home', 'work', 'fax') NOT NULL,
  `number` VARCHAR(30) NOT NULL,
  PRIMARY KEY(`id`),
  CONSTRAINT `fk_phone_contact` FOREIGN KEY (`contact_id`)
    REFERENCES `contact`(`id`)
    ON DELETE CASCADE
);