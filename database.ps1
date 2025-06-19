net start MySQL93
net start "Mongo DB Server (MongoDB)"

$mysql = @'
  CREATE TABLE IF NOT EXISTS `urmid`.`user` (
  `username` VARCHAR(32) NOT NULL,
  `email` VARCHAR(255) NULL,
  `password` VARCHAR(32) NOT NULL,
  `create_time` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `uuid` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `firstname` VARCHAR(255) NULL,
  `lastname` VARCHAR(255) NULL,
  `streetnumber` INT NULL,
  `street` VARCHAR(255) NULL,
  `city` VARCHAR(255) NULL,
  `state` VARCHAR(2) NULL,
  `country` VARCHAR(255) NULL,
  `zipcode` INT NULL,
  PRIMARY KEY (`uuid`),
  UNIQUE INDEX `uuid_UNIQUE` (`uuid` ASC) VISIBLE)

  CREATE TABLE IF NOT EXISTS `urmid`.`profile` (
  `firstname` VARCHAR(255) NOT NULL,
  `uuid` INT UNSIGNED ZEROFILL NOT NULL,
  `title` VARCHAR(255) NULL,
  `fuzzedlat` INT NULL,
  `fuzzedlong` INT NULL,
  `profile` LONGTEXT NULL,
  `job` VARCHAR(45) NULL,
  `energy` TINYINT(1) NULL,
  `dogs` TINYINT(1) NULL,
  `cats` TINYINT(1) NULL,
  `kids` TINYINT(1) NULL,
  `smoke` TINYINT(1) NULL,
  `drink` TINYINT(1) NULL,
  `age` INT NULL,
  `gender_identity` VARCHAR(45) NULL,
  `sexual_preference` VARCHAR(45) NULL,
  `relationship_desire` VARCHAR(45) NULL,
  `religion` VARCHAR(45) NULL,
  `religious_importance` TINYINT(1) NULL,
  `languages` VARCHAR(45) NULL,
  `dietary` VARCHAR(45) NULL,
  PRIMARY KEY (`uuid`),
  UNIQUE INDEX `uuid_UNIQUE` (`uuid` ASC) VISIBLE,
  CONSTRAINT `uuid`
    FOREIGN KEY (`uuid`)
    REFERENCES `urmid`.`user` (`uuid`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
'@