#!/bin/sh

ssh -q -o "StrictHostKeyChecking no" $1 exit
while test $? -gt 0
do
   sleep 1 
   echo "Trying again..."
   ssh -q $1 -o "StrictHostKeyChecking no" exit
done
