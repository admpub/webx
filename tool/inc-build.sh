cd ..
go generate
cd tool

export DISTPATH="../dist"
export OSVERSIONDIR=${NGING_EXECUTOR}_${GOOS}_${GOARCH}
export RELEASEDIR=${DISTPATH}/${OSVERSIONDIR}
mkdir -p ${RELEASEDIR}
go build -tags "bindata official${BUILDTAGS}" -ldflags="-X main.BUILD_TIME=${NGING_BUILD} -X main.COMMIT=${NGING_COMMIT} -X main.VERSION=${NGING_VERSION} -X main.LABEL=${NGING_LABEL} -X main.BUILD_OS=${GOOS} -X main.BUILD_ARCH=${GOARCH} -s -w -extldflags '-static'" -o ${RELEASEDIR}/${NGING_EXECUTOR}${NGINGEX} ..
mkdir ${RELEASEDIR}/data
mkdir ${RELEASEDIR}/data/logs
cp -R ../data/ip2region ${RELEASEDIR}/data/ip2region


mkdir ${RELEASEDIR}/config
#mkdir ${RELEASEDIR}/config/vhosts

#cp -R ../config/config.yaml ${RELEASEDIR}/config/config.yaml
cp -R ../config/config.yaml.sample ${RELEASEDIR}/config/config.yaml.sample
cp -R ../config/insert.* ${RELEASEDIR}/config/
# cp -R ../config/preupgrade.* ${RELEASEDIR}/config/
cp -R ../config/ua.txt ${RELEASEDIR}/config/ua.txt

export archiver_extension="tar.gz"

rm -rf ${RELEASEDIR}.${archiver_extension}

tar -zcvf ${RELEASEDIR}.${archiver_extension} -C ${DISTPATH} ${OSVERSIONDIR}

rm -rf ${RELEASEDIR}
