
source ${PWD}/inc-version.sh

#go get github.com/admpub/xgo
#source ${WORKDIR}/install-archiver.sh

cd ..
go generate

cd tool

export TMPDIR=

export NGINGEX=
export BUILDTAGS=

export GOOS=linux
export GOARCH=amd64
source ${PWD}/inc-build.sh

