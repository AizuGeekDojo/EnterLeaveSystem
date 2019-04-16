import subprocess
def nfc_read():
    p = subprocess.Popen("python2 nfc_reader.py", shell=True,stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    stdout_data, stderr_data = p.communicate("".encode("utf-8"))
    rawdat = stdout_data.decode("utf-8").replace("\n", "")
    spldat = rawdat.split()
    if len(spldat) > 1:
      return spldat[0], spldat[1]
    else:
      return spldat[0]
