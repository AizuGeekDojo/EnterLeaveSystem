#!/usr/bin/env python3
"""NFC card reader for student ID and university cards."""

import sys
import nfc

SERVICE_CODE = 0x300b


def on_connect(tag):
    """Handle NFC tag connection and read card data."""
    if isinstance(tag, nfc.tag.tt3.Type3Tag):
        try:
            sc = nfc.tag.tt3.ServiceCode(
                number=SERVICE_CODE >> 6, attribute=SERVICE_CODE & 0x3f
            )
            bc = nfc.tag.tt3.BlockCode(0, access=0, service=0)
            data = tag.read_without_encryption([sc], [bc]).decode("utf-8")

            if data[0:4] == "1000":
                sid = data[4:11]
                if sid[0] == "1":
                    print(f"student s{sid}")
                elif sid[0] == "5":
                    print(f"student m{sid}")
                elif sid[0] == "8":
                    print(f"student d{sid}")
            else:
                print(f"univ {data}")
        except Exception as e:
            # General card - use IDm as identifier
            idm_hex = "".join([f"{s:02x}" for s in tag.idm])
            print(f"general {idm_hex}")


def main():
    """Main function to initialize NFC reader and wait for cards."""
    clf = nfc.ContactlessFrontend('usb')
    try:
        clf.connect(rdwr={'on-connect': on_connect})
    finally:
        clf.close()


if __name__ == "__main__":
    main()
