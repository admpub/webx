
source ${PWD}/inc-version.sh

# 回到入口
cd ${ENTRYDIR}

export NGINGEX=
export BUILDTAGS=

export GOOS=linux
export GOARCH=amd64
source ${WORKDIR}/inc-build-x.sh

