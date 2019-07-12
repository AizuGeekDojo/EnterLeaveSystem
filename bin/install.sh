#! /bin/sh

if [ ! "$UID" = "0" ];then
  echo "You are not root."
  exit 1
fi

echo "Updating system..."
apt update
#apt upgrade

echo "Installing python and environment..."

apt install -y python2.7 python-pip

python2 -m pip install --upgrade pip
python2 -m pip install nfcpy

echo "Installing service..."
echo "Copying files..."
cp elsystemd.service /etc/systemd/system/

echo "Enabling Services..."
systemctl daemon-reload
systemctl enable elsystemd.service
systemctl start elsystemd.service

echo "Configurating autostart"
mv /home/pi/.config/lxsession/LXDE-pi/autostart /home/pi/.config/lxsession/LXDE-pi/autostart_old
cp autostart /home/pi/.config/lxsession/LXDE-pi/autostart

echo "Finished! Please restart the system."
