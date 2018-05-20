import sqlite3

dbname = 'database.db'

def getSIDByIDm(IDm):
    conn = sqlite3.connect(dbname)
    c = conn.cursor()
    c.execute('select sid from idcard where idm=?',(IDm,))
    ret = c.fetchone()
    conn.close()
    if ret is not None:
        return ret[0]
    return None

def getUserName(sid):
    conn = sqlite3.connect(dbname)
    c = conn.cursor()
    c.execute('select name from users where sid=?',(sid,))
    ret = c.fetchone()
    conn.close()
    if ret is not None:
        return ret[0]
    return None

def addLog(sid,isenter,time,ext):
    conn = sqlite3.connect(dbname)
    c = conn.cursor()
    c.execute('insert into log values(?,?,?,?)',(sid,isenter,time,ext))
    if isenter == 1:
        c.execute('update users set isenter=? where sid=?',(0,sid))
    else:
        c.execute('update users set isenter=? where sid=?',(1,sid))

    conn.commit()
    conn.close()

def addUser(idm,sid):
    conn = sqlite3.connect(dbname)
    c = conn.cursor()
    c.execute('insert into idcard values(?,?)',(idm,sid))
    conn.commit()
    conn.close()

def updateUser(idm,sid):
    conn = sqlite3.connect(dbname)
    c = conn.cursor()
    c.execute('update idcard set idm=?,sid=? where sid=?',(idm,sid,sid))
    conn.commit()
    conn.close()

def isUserInside(sid):
    conn = sqlite3.connect(dbname)
    c = conn.cursor()
    c.execute('select isenter from users where sid=?',(sid,))
    ret = c.fetchone()
    conn.close()
    return ret[0]==1

#sid = getSIDByIDm("0000000000000000")
#print(sid)
#print(getUserName(sid))
#addLog(sid,"In","20180519160755")
