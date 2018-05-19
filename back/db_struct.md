# Database仕様
## 学生証情報
```sql
  create table idcard (idm text,sid text);
```
## 学生情報
```sql
  create table users (sid text,name text);
```
## 入退室記録
```sql
  create table idcard (sid text,intime integer,outtime integer);
```
