# generate
cd ${PKGPATH}
go generate
cd ${ENTRYDIR}

export DISTPATH=${PKGPATH}/dist
export OSVERSIONDIR=${NGING_EXECUTOR}_${GOOS}_${GOARCH}
export RELEASEDIR=${DISTPATH}/${OSVERSIONDIR}
if [ "$GOARM" != "" ]; then
	export RELEASEDIR=${RELEASEDIR}v${GOARM}
fi
mkdir ${RELEASEDIR}

export LDFlagsExt=""
export LDFlagsExt=" -s -w -extldflags '-static'"

xgo -go=latest -goproxy=https://goproxy.cn,direct -image="admpub/xgo:1.17.4" -targets=${GOOS}/${GOARCH} -dest=${RELEASEDIR} -out=${NGING_EXECUTOR} -tags="bindata official sqlite zbar${BUILDTAGS}" -ldflags="-X main.BUILD_TIME=${NGING_BUILD} -X main.COMMIT=${NGING_COMMIT} -X main.VERSION=${NGING_VERSION} -X main.LABEL=${NGING_LABEL} -X main.BUILD_OS=${GOOS} -X main.BUILD_ARCH=${GOARCH} ${LDFlagsExt}" ./${PKGPATH}

mv ${RELEASEDIR}/${NGING_EXECUTOR}-${GOOS}-* ${RELEASEDIR}/${NGING_EXECUTOR}${NGINGEX}
mkdir ${RELEASEDIR}/data
mkdir ${RELEASEDIR}/data/logs
cp -R ${PKGPATH}/data/ip2region ${RELEASEDIR}/data/ip2region

mkdir ${RELEASEDIR}/config
mkdir ${RELEASEDIR}/config/vhosts

#cp -R ../config/config.yaml ${RELEASEDIR}/config/config.yaml
cp -R ${PKGPATH}/config/config.yaml.sample ${RELEASEDIR}/config/config.yaml.sample
cp -R ${PKGPATH}/config/insert.* ${RELEASEDIR}/config/
cp -R ${PKGPATH}/config/preupgrade.* ${RELEASEDIR}/config/
cp -R ${PKGPATH}/config/ua.txt ${RELEASEDIR}/config/ua.txt

export archiver_extension="tar.gz"

rm -rf ${RELEASEDIR}.${archiver_extension}

tar -zcvf ${RELEASEDIR}.${archiver_extension} -C ${DISTPATH} ${OSVERSIONDIR}

rm -rf ${RELEASEDIR}