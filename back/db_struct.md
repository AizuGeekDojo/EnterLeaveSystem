# Database仕様
## 学生証情報
```sql
  create table idcard (idm text,sid text);
```
## 学生情報
```sql
  create table users (sid text,name text,isenter integer);
```
## 入退室記録
```sql
  create table log (sid text,isenter integer,time integer,ext text);
```
