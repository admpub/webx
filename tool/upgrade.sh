cd /home/webx
./webx service stop
cd /home
unzip webx_linux_amd64.zip -d ./webx
cp -R ./webx/webx_linux_amd64/* ./webx
rm -rf ./webx/webx_linux_amd64
cd /home/webx
./webx service start
