#!/bin/bash
echo "Hello world" > index.html
sudo nohup busybox httpd -f -p 80 &
