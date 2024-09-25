cd /home/webx
./webx service stop
cd /home

if [ "$1" = "" ]; then
  unzip ./webx_linux_amd64.zip -d ./webx
else
  tar -zxvf ./webx_linux_amd64.tar.gz -C ./webx
fi

cp -R ./webx/webx_linux_amd64/* ./webx
rm -rf ./webx/webx_linux_amd64
cd /home/webx
./webx service start
