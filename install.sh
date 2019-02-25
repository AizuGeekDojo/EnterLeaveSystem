#! /bin/sh

# if [ ! "$UID" = "0" ];then
#   echo "You are not root."
#   exit 1
# fi

echo "Updating system..."
apt update
#apt upgrade

echo "Installing python and environment..."

apt install -y python2.7 python-pip

sudo apt install -y build-essential tk-dev libncurses5-dev libncursesw5-dev libreadline6-dev libdb5.3-dev libgdbm-dev libsqlite3-dev libssl-dev libbz2-dev libexpat1-dev liblzma-dev zlib1g-dev

wget https://www.python.org/ftp/python/3.6.0/Python-3.6.0.tar.xz
tar xf Python-3.6.0.tar.xz
cd Python-3.6.0
./configure
make
sudo make altinstall
cd ../

python2 -m pip install --upgrade pip
python2 -m pip install nfcpy
python3 -m pip install --upgrade pip
python3 -m pip install -r back/requirements.txt

echo "Installing service..."
echo "Copying files..."
# cp enterleavesystemd.service /etc/systemd/system/
# cp enterleavesystemhttpd.service /etc/systemd/system/
cp elsystemd.service /etc/systemd/system/

echo "Enabling Services..."
systemctl daemon-reload
# systemctl enable enterleavesystemd.service
# systemctl start enterleavesystemd.service

# systemctl enable enterleavesystemhttpd.service
# systemctl start enterleavesystemhttpd.service

systemctl enable elsystemd.service
systemctl start elsystemd.service

echo "Configurating autostart"
mv /home/pi/.config/lxsession/LXDE-pi/autostart /home/pi/.config/lxsession/LXDE-pi/autostart_old
cp autostart /home/pi/.config/lxsession/LXDE-pi/autostart

echo "Finished! Please restart the system."
