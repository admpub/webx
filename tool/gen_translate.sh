#go install github.com/admpub/i18n/cmd/fetchtext@latest
fetchtext --src=../ --dist=../config/i18n/messages --default=zh-cn --translate=true --onlyExport=false --clean=true --translator=tencent --translatorConfig="appid=&secret=" --vendorDirs="github.com/admpub/nging"
