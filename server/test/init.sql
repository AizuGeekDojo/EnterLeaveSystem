PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS "users" (sid TEXT,name TEXT,isenter INTEGER);
INSERT INTO users VALUES('s1240095','Yusuke Namiki',0);
CREATE TABLE IF NOT EXISTS "idcard" (idm TEXT,sid TEXT);
CREATE TABLE IF NOT EXISTS "products" (id TEXT,name TEXT,barcode TEXT,borrowersid TEXT);
INSERT INTO products VALUES('1','テスト商品A','','s1240095');
INSERT INTO products VALUES('2','テスト商品B','12345678','s1240095');
CREATE TABLE IF NOT EXISTS "log" (sid TEXT,isenter INTEGER,time INTEGER,ext TEXT);
COMMIT;

