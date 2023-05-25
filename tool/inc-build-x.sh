export DISTPATH=${PKGPATH}/dist
export RELEASEDIR=${DISTPATH}/${NGING_EXECUTOR}_${GOOS}_${GOARCH}
if [ "$GOARM" != "" ]; then
	export RELEASEDIR=${RELEASEDIR}v${GOARM}
fi
mkdir ${RELEASEDIR}

export LDFlagsExt=""
#export LDFlagsExt="-linkmode external -extldflags '-static'"

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

if [ "$GOOS" = "windows" ]; then
    cp -R ${PKGPATH}/support/sqlite3_${GOARCH}.dll ${RELEASEDIR}/
	export archiver_extension=zip
else
	export archiver_extension=zip
fi

rm -rf ${RELEASEDIR}.${archiver_extension}

#${NGING_VERSION}${NGING_LABEL}
arc archive ${RELEASEDIR}.${archiver_extension} ${RELEASEDIR}

rm -rf ${RELEASEDIR}