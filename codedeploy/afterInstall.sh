#!/bin/bash

sudo systemctl stop gowebapp
sudo chown ec2-user:ec2-user /home/ec2-user/webservice/webapp

# cleanup log files
sudo rm -rf /home/ec2-user/webservice/*.log
