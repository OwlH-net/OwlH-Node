sudo yum -y install epel-release
sudo yum -y install wget 
sudo yum -y install gcc libpcap-devel pcre-devel libyaml-devel file-devel \
  zlib-devel jansson-devel nss-devel libcap-ng-devel libnet-devel tar make \
  libnetfilter_queue-devel lua-devel
sudo yum  -y install python-devel

pip install --upgrade pip
pip install pyyaml

if [ ! -f ./suricata-4.1.0-rc2.tar.gz ]; then
    wget https://openinfosecfoundation.org/download/suricata-4.1.0-rc2.tar.gz
    tar -xvzf suricata-4.1.0-rc2.tar.gz
fi

cd suricata-4.1.0-rc2
./configure --prefix=/usr --sysconfdir=/etc --localstatedir=/var --enable-nfqueue --enable-lua

make
sudo make install-full
sudo ldconfig