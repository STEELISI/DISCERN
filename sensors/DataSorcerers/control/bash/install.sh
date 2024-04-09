#!/bin/bash

# This needs to be run as su

apt update
apt -y install perl
apt -y install git
apt -y install strace

pushd .

cd ~

git clone https://github.com/STEELISI/ACSLE.git

mkdir -p /usr/local/src/ttylog
mkdir -p /var/log/ttylog

cp ACSLE/monitor/analyze.py /usr/local/src/ttylog/
# cp ACSLE/monitor/pre_process.py /usr/local/src/
cp ACSLE/monitor/script.sh /usr/local/src/
cp ACSLE/monitor/ttylog /usr/local/src/ttylog/

rm -rf ACSLE

popd

cp ./start_ttylog.sh /usr/local/src/start_ttylog.sh

echo 'ForceCommand /usr/local/src/script.sh' >> /etc/ssh/sshd_config # (Start 'script.sh' as soon as a user SSH's in)
# echo 'python3 /usr/local/src/pre_process.py &' >> /usr/local/etc/emulab/rc/rc.testbed # (Launch 'pre_process.py' at system startup)
systemctl restart sshd

echo "Restart system to finish bash installation"

