def nfc_read():
    p = subprocess.Popen("python nfc.py", shell=True,stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    stdout_data, stderr_data = p.communicate("".encode("utf-8"),timeout=20)
    return stdout_data.decode("utf-8")
