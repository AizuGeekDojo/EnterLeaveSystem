import binascii
import nfc

def on_connect(tag):
    print binascii.hexlify(tag.idm)
    return False

clf = nfc.ContactlessFrontend('usb')
try:
    clf.connect(rdwr={'on-connect': on_connect})
finally:
    clf.close()
