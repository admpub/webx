export PKGPATH=github.com/admpub/webx
export ENTRYDIR=${GOPATH}/src
export WORKDIR=${PWD}

export NGING_VERSION=`git describe --always --dirty | sed -e "s/^v//"`
export NGING_BUILD=`date +%Y%m%d%H%M%S`
export NGING_COMMIT=`git rev-parse --short HEAD`
export NGING_LABEL="stable"
export NGING_EXECUTOR="webx"