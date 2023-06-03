#!/bin/bash
sudo apt update
sudo apt install apache2

INTERFACE=$(curl -s http://169.254.169.254/latest/meta-data/network/interfaces/macs/)
SUBNETID=$(curl -s http://169.254.169.254/latest/meta-data/network/interfaces/macs/${INTERFACE}/subnet-id)
echo '<center><h1>subnet: SUBNETID </h1></center>' > /var/www/html/index.txt
sed "s/SUBNETID/$SUBNETID/" /var/www/html/index.txt > /var/www/html/index.html

sudo ufw app list
sudo ufw allow 'Apache'
