#! /bin/sh

echo "Updating system..."
apt update
#apt upgrade

echo "Installing python and environment..."

apt install -y python2.7 python-pip

python2 -m pip install --upgrade pip
python2 -m pip install nfcpy

echo "Installing other packages..."
apt install -y sqlite3 xscreensaver

echo "Installing service..."
echo "Copying files..."
cp elsystemd.service elsystemf.service /etc/systemd/system/

echo "Enabling Services..."
systemctl daemon-reload
systemctl enable elsystemd.service
systemctl enable elsystemf.service
systemctl start elsystemd.service
systemctl start elsystemf.service

echo "Setting config.txt for display config..."
echo "hdmi_blanking=1" >> /boot/config.txt

echo "Finished! Please restart the system."
