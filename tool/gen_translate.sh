#go install github.com/admpub/i18n/cmd/fetchtext@latest
if [ "$1" = "" ]; then
fetchtext --src=../ --dist=../config/i18n/messages --default=zh-CN --translate=true --onlyExport=false --clean=true\
 --translator=tencent --translatorConfig="appid=&secret="\
 --vendorDirs="github.com/admpub/nging"\
 --envFile="$HOME/go/src/github.com/admpub/nging/tool/.translator_tencent.env" --onlyTranslateIncr=true
else
translator --src=../ --dist=../config/i18n/messages --default=zh-CN --translate=true --onlyExport=false --clean=true\
 --translator=ollama --translatorConfig="appid=&secret="\
 --vendorDirs="github.com/admpub/nging"\
 --onlyTranslateIncr=true
fi