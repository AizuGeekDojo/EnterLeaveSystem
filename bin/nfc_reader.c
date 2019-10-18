// Reference: https://qiita.com/ysomei/items/32f366b61a7b631c4750
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>

#include "libusb.h"

#define PRODUCT_ID 0x06C3  // RC-S380
#define VENDOR_ID 0x054C   // SONY
#define TIMEOUT 5000

struct device_info {
  libusb_device *dev;
  libusb_device_handle *dh;
  unsigned char ep_in;
  unsigned char ep_out;
  int interface_num;
};
typedef struct device_info usb_device_info;
void show_data(unsigned char *buf, int size) {
  for(int i = 0; i != size; i++){
    printf("%02x", buf[i]);
  }
  printf("\n");  
}


usb_device_info *get_usb_information(libusb_device_handle *dh) {
  usb_device_info *devinfo;
  devinfo = (usb_device_info *)malloc(sizeof(usb_device_info));
  if(devinfo == NULL) return NULL;
  memset(devinfo, 0, sizeof(usb_device_info));

  libusb_device *dev;
  struct libusb_config_descriptor *conf;
  const struct libusb_endpoint_descriptor *endp;
  const struct libusb_interface *intf;
  const struct libusb_interface_descriptor *intdesc;
  int ret;

  dev = libusb_get_device(dh);
  if(dev == NULL){
    return NULL;
  }
  devinfo->dh = dh;
  devinfo->dev = dev;

  libusb_get_config_descriptor(dev, 0, &conf);  
  for(int i = 0; i < (int)conf->bNumInterfaces; i++){
    intf = &conf->interface[i];
    for(int j = 0; j < intf->num_altsetting; j++){
      intdesc = &intf->altsetting[j];
      for(int k = 0; k < (int)intdesc->bNumEndpoints; k++){
        endp = &intdesc->endpoint[k];

        switch(endp->bmAttributes & LIBUSB_TRANSFER_TYPE_MASK) {
        case LIBUSB_TRANSFER_TYPE_BULK:
          if((endp->bEndpointAddress & 0x80) == LIBUSB_ENDPOINT_IN){
            devinfo->ep_in = endp->bEndpointAddress;
          }
          if((endp->bEndpointAddress & 0x80) == LIBUSB_ENDPOINT_OUT){
            devinfo->ep_out = endp->bEndpointAddress;
          }
          break;
        case LIBUSB_TRANSFER_TYPE_INTERRUPT:
          break;
        }
      }
    }    
  }
  libusb_free_config_descriptor(conf);

  return devinfo;
}
int nfc_init(usb_device_info *di){
  unsigned char cmd[6]={0x00,0x00,0xff,0x00,0xff,0x00};

  return libusb_bulk_transfer(di->dh, di->ep_out, cmd, sizeof(cmd), NULL, TIMEOUT);
}

int nfc_send(usb_device_info *devinfo, unsigned char *buf, int size,unsigned char *res,int reslen) {
  unsigned char cmd[255]={0x00,0x00,0xff,0xff,0xff};
  int len,i;
  short csum;
  int ret;

  // 00 00 ff ff ff len(L) len(H) checksum(len) data checksum(data) 00

  cmd[5] = ((size + 1) & 0xff);
  cmd[6] = ((size + 1) & 0xff00) >> 8;
  csum = (0x100 - (cmd[5] + cmd[6])) % 0x100;
  cmd[7] = csum;

  cmd[8] = 0xd6;
  for (i=0;i<size;i++)
    cmd[i+9]=buf[i];
  

  csum = cmd[8];
  for(int i = 0; i < size; i++)
    csum += buf[i];
  
  csum = (0x100 - csum) % 0x100; 
  cmd[9 + size] = csum;

  cmd[10 + size] = 0x00;


  ret = libusb_bulk_transfer(devinfo->dh, devinfo->ep_out, cmd, 11 + size, NULL, TIMEOUT);
  if(ret != 0){
    return -1;
  }

  // receive ack/nck
  ret = libusb_bulk_transfer(devinfo->dh, devinfo->ep_in, cmd, sizeof(cmd), NULL, TIMEOUT);
  if(ret != 0){
    return -2;
  }

  // receive response
  if (res != NULL) {
    ret = libusb_bulk_transfer(devinfo->dh, devinfo->ep_in, res, reslen, &len, TIMEOUT);
  }else{
    ret = libusb_bulk_transfer(devinfo->dh, devinfo->ep_in, cmd, sizeof(cmd), &len, TIMEOUT);
  }
  if(ret != 0){
    return -3;
  }

  return len;
}
int packet_setcommandtype(usb_device_info *devinfo) {
  unsigned char cmd[2]={0x2a,0x01};
  return nfc_send(devinfo, cmd, sizeof(cmd),NULL,0);
}
int packet_switch_rf(usb_device_info *devinfo) {
  unsigned char cmd[2]={0x06,0x00};
  return nfc_send(devinfo, cmd, sizeof(cmd),NULL,0);
}
int packet_inset_rf(usb_device_info *devinfo) {
  unsigned char cmd[5]={0x00,0x01,0x01,0x0f,0x01};
  return nfc_send(devinfo, cmd, sizeof(cmd),NULL,0);  
}
int packet_inset_protocol_1(usb_device_info *devinfo) {
  unsigned char cmd[39]=
  {0x02,0x00,0x18,0x01,0x01,0x02,0x01,0x03,0x00,0x04,
  0x00,0x05,0x00,0x06,0x00,0x07,0x08,0x08,0x00,0x09,
  0x00,0x0a,0x00,0x0b,0x00,0x0c,0x00,0x0e,0x04,0x0f,
  0x00,0x10,0x00,0x11,0x00,0x12,0x00,0x13,0x06};
  return nfc_send(devinfo, cmd, sizeof(cmd),NULL,0);  
}
int packet_inset_protocol_2(usb_device_info *devinfo) {
  unsigned char cmd[3]={0x02,0x00,0x18};
  return nfc_send(devinfo, cmd, sizeof(cmd),NULL,0);  
}
int packet_sens_req(usb_device_info *devinfo, unsigned char *buf,int size) {
  // 04 (InCommRF) 6e 00 (timeout) 06 (size including size) 00 (felica polling)
  // ff ff (system code :wild) 01 (request code) 00 (timeslots)
  unsigned char cmd[9]={0x04,0x6e,0x00,0x06,0x00,0xff,0xff,0x01,0x00};


  return nfc_send(devinfo, cmd, sizeof(cmd),buf,size);  
}
int packet_felicaread(usb_device_info *devinfo, unsigned char *buf,int size,unsigned char *idm) {
  unsigned char cmd[19]={0x04,0x43,0x60,0x10,0x06,
                  0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,
                  0x01,0x0b,0x30,0x01,0x80,0x00};
  int i;

  for (i = 0;i< 8; i++)
    cmd[i+5]=idm[i];

  //04 (InCommRF) 6e 00 (timeout) 10 (size including size)
  //06 (Read W/O Enc) 01 16 06 00 3e 16 dd 01 (idm)
  //01 (num of sev code) 8b 1a (sev code) 01 (num of blk code) 80 00 (blk code)

  return nfc_send(devinfo, cmd, sizeof(cmd), buf, size);
}

usb_device_info *opennfc(){
  libusb_device_handle *dh = NULL;
  usb_device_info *devinfo;


  if (libusb_init(NULL)!=0) {
    return NULL;
  }
  dh = libusb_open_device_with_vid_pid(NULL, VENDOR_ID, PRODUCT_ID);

  if (dh == NULL){
    libusb_exit(NULL);
    return NULL;
  }

  // usb interface setting
  libusb_set_auto_detach_kernel_driver(dh, 1);
  libusb_set_configuration(dh, 1);

  libusb_claim_interface(dh, 0);
  libusb_set_interface_alt_setting(dh, 0, 0);

  // get usb information
  devinfo = get_usb_information(dh);
  nfc_init(devinfo);
  packet_setcommandtype(devinfo);
  packet_switch_rf(devinfo);

  packet_inset_rf(devinfo);
  packet_inset_protocol_1(devinfo);    
  packet_inset_protocol_2(devinfo);

  return devinfo;
}




int main(){
  int i,rlen;
  usb_device_info *devinfo = opennfc();
  unsigned char rbuf[256];


  while (1){      
  i = packet_sens_req(devinfo, rbuf,256);

  if(rbuf[9] == 0x05 && rbuf[10] == 0x00){

    rlen = ((rbuf[6] << 8) + rbuf[5]);
    for (i=0;i<rlen - 10;i++){
      rbuf[i]=rbuf[i+9];
    }
    //      // Type-F
    if(rbuf[6] == 0x14 && rbuf[7] == 0x01){

      printf(" IDm: "); show_data(rbuf + 8, 8);
      printf(" PMm: "); show_data(rbuf + 16, 8);

      packet_felicaread(devinfo,rbuf,255,rbuf + 8);

      show_data(rbuf,46);
      printf("%s\n",&rbuf[28]);

    }
        usleep(2500 * 1000);

  }
    usleep(250 * 1000);
  }

  //     // close
  libusb_release_interface(devinfo->dh, 0);
  libusb_close(devinfo->dh);

  libusb_exit(NULL);

}