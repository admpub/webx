go install github.com/admpub/nging-builder@latest
platform="linux_amd64,linux_386,freebsd_amd64,freebsd_386,linux_arm,linux_arm64,windows_amd64,windows_386,darwin_amd64,linux_mips,linux_mips64,linux_mipsle,linux_mips64le"
if [ "$1" != "" ]; then
    platform="$1"
fi
nging-builder --conf ./builder.conf $platform min