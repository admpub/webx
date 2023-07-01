
source ${PWD}/inc-version.sh

go get github.com/admpub/xgo
source ${WORKDIR}/install-archiver.sh

# 回到入口
cd ${ENTRYDIR}

export NGINGEX=
export BUILDTAGS=

export GOOS=darwin
export GOARCH=amd64
source ${WORKDIR}/inc-build-x.sh

