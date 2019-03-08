#! /usr/local/bin/python3

# /usr/bin/python3
# /usr/local/bin/python3

import os
import sys
import json
import sqlite3
import datetime
import requests

# TOKEN = os.getenv("SLACK_FILEUP_TOKEN")
# CHANNEL = os.getenv("SLACK_FILEUP_CHANNEL")
TOKEN = ""
CHANNEL = ""



dbname = 'database.db'
csvname = 'out.csv'

def getAllLog():
    conn = sqlite3.connect(dbname)
    c = conn.cursor()
    c.execute('select name,log.isenter,time,log.sid,ext from log,users where log.sid=users.sid')
    ret = c.fetchall()
    conn.close()
    return ret

def removeAllLog():
    conn = sqlite3.connect(dbname)
    conn.execute('delete from log')
    conn.commit()
    conn.close()

ret = getAllLog()

csv = ["Date,StudentID,Name,Enter/Leave,Purpose,Comment\n"]
for i in ret:
    date = str(datetime.datetime.fromtimestamp(int(int(i[2])/1000)))
    ext = json.loads(i[4])
    comment = ""
    usepurpose = ""
    entstr = "N/A"

    if ('message' in ext):
        comment = ext['message']

    if ('Use' in ext):
        for u in ext['Use']:
            usepurpose += u + " "
    
    if int(i[1]) == 1:
        entstr = "Enter"
    elif int(i[1]) == 0:
        entstr = "Leave"

    csv.append(date + "," + i[3] + "," + i[0] + "," + entstr + "," + usepurpose + ",\"" + comment + "\"\n")

print

with open(csvname, mode='w') as f:
    f.writelines(csv)

files = {'file': open(csvname, 'rb')}
param = {
    'token':TOKEN,
    'channels':CHANNEL,
    'filename':"log.csv",
    'initial_comment': "Enter leave log (csv format) by csvexport",
    'title': "log.csv"
}
r = requests.post(url="https://slack.com/api/files.upload",params=param, files=files)

if r.json()["ok"] :
    removeAllLog()