import sys
import nfc

service_code = 0x300b

def on_connect(tag):
    # print tag
    if isinstance(tag, nfc.tag.tt3.Type3Tag):
        try:
            sc = nfc.tag.tt3.ServiceCode(number=service_code >> 6 ,attribute=service_code & 0x3f)
            bc = nfc.tag.tt3.BlockCode(0,access=0,service=0)
            data = tag.read_without_encryption([sc],[bc]).decode("utf-8")
            if data[0:4] == "1000":
                sid = data[4:11]
                if sid[0] == "1":
                    print "student s" + sid
                elif sid[0] == "5":
                    print "student m" + sid
                elif sid[0] == "8":
                    print "student d" + sid
            else:
                print "univ " + data
        except Exception as e:
            print "general " + "" . join(['%02x' % s for s in tag.idm])

clf = nfc.ContactlessFrontend('usb')
try:
    clf.connect(rdwr={'on-connect': on_connect})
finally:
    clf.close()
