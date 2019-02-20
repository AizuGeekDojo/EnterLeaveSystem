import binascii
# import os
# import struct
# import sys
import nfc

# service_code = 0x1a8b
service_code = 0x300b
num_blocks = 4



def on_connect(tag):
    # print binascii.hexlify(tag.idm)
    # return False

    print tag
    # for i in xrange(0x10000):
    #     area_or_service = tag.search_service_code(i)
    #     if area_or_service is None:
    #     break
    #     elif len(area_or_service) == 1:
    #     sc = area_or_service[0]
    #     print(nfc.tag.tt3.ServiceCode(sc >> 6, sc & 0x3f))
    #     elif len(area_or_service) == 2:
    #     area_code, area_last = area_or_service
    #     print("Area {0:04x}--{0:04x}".format(area_code, area_last))
    if isinstance(tag, nfc.tag.tt3.Type3Tag):
        try:
        for i in range(num_blocks):
            sc = nfc.tag.tt3.ServiceCode(number=service_code >> 6 ,attribute=service_code & 0x3f)
            bc = nfc.tag.tt3.BlockCode(i,access=0,service=0)
            data = tag.read_without_encryption([sc],[bc])
            print "" . join(['%2x ' % (s - 0x30) for s in data])
        except Exception as e:
        # print "error: %s" % e
    # else:
        # print "error: tag isn't Type3Tag"


clf = nfc.ContactlessFrontend('usb')
try:
    clf.connect(rdwr={'on-connect': on_connect})
finally:
    clf.close()
