module github.com/admpub/webx

go 1.24.3

// replace github.com/admpub/nging/v5 => ../../../github.com/admpub/nging

// replace github.com/coscms/webcore => ../../coscms/webcore

// replace github.com/coscms/webfront => ../../coscms/webfront

exclude github.com/gomodule/redigo v2.0.0+incompatible

require (
	github.com/admpub/nging/v5 v5.3.4-0.20250731031910-eb90ae50e19c
	github.com/coscms/webcore v0.9.31
	github.com/coscms/webfront v0.0.0-20250731035826-ee75c9ee5c96
)

require (
	github.com/adamzy/cedar-go v0.0.0-20170805034717-80a9c64b256d // indirect
	github.com/admpub/bindata/v3 v3.2.1
	github.com/admpub/cache v0.8.0
	github.com/admpub/color v1.8.1 // indirect
	github.com/admpub/copier v0.1.1 // indirect
	github.com/admpub/decimal v1.3.2 // indirect
	github.com/admpub/errors v0.8.2
	github.com/admpub/fasttemplate v0.0.2 // indirect
	github.com/admpub/go-download/v2 v2.2.0 // indirect
	github.com/admpub/go-hashids v2.0.1+incompatible // indirect
	github.com/admpub/go-shortid v0.0.0-20140827050853-24d054c393fe // indirect
	github.com/admpub/godotenv v1.4.3
	github.com/admpub/imageproxy v0.10.1
	github.com/admpub/ipfilter v1.0.6 // indirect
	github.com/admpub/license_gen v0.1.1
	github.com/admpub/log v1.4.0
	github.com/admpub/marmot v0.0.0-20200702042226-2170d9ff59f5 // indirect
	github.com/admpub/null v8.0.5+incompatible
	github.com/admpub/once v0.0.1 // indirect
	github.com/admpub/pinyin-golang v1.0.1 // indirect
	github.com/admpub/pp v0.0.7
	github.com/admpub/redsync/v4 v4.0.3 // indirect
	github.com/admpub/resty/v2 v2.7.2 // indirect
	github.com/admpub/sensitive v0.0.1 // indirect
	github.com/admpub/useragent v0.0.2 // indirect
	github.com/caddy-plugins/ipfilter v1.1.8 // indirect
	github.com/coscms/forms v1.13.10
	github.com/coscms/oauth2s v0.4.3 // indirect
	github.com/coscms/sms v0.0.7
	github.com/gosimple/slug v1.15.0 // indirect
	github.com/huichen/sego v0.0.0-20210824061530-c87651ea5c76 // indirect
	github.com/prometheus/client_golang v1.22.0 // indirect
	github.com/spf13/cobra v1.9.1 // indirect
	github.com/stretchr/testify v1.10.0
	github.com/swaggo/swag v1.16.3
	github.com/webx-top/client v0.9.6
	github.com/webx-top/com v1.3.31
	github.com/webx-top/db v1.28.6
	github.com/webx-top/echo v1.16.12
	github.com/webx-top/echo-prometheus v1.1.2 // indirect
	github.com/webx-top/image v0.1.2
	github.com/webx-top/pagination v0.3.2
	github.com/webx-top/validation v0.0.3 // indirect
	golang.org/x/oauth2 v0.30.0 // indirect
	golang.org/x/sync v0.16.0
	gopkg.in/redis.v5 v5.2.9 // indirect
)

require (
	github.com/admpub/events v1.3.6
	github.com/admpub/goth v0.0.4
	github.com/admpub/sessions v0.3.0
	github.com/nging-plugins/dbmanager v1.8.7
	github.com/silenceper/wechat/v2 v2.1.9
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	gitee.com/admpub/certmagic v0.8.9 // indirect
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/abh/errorutil v1.0.0 // indirect
	github.com/admpub/9t v0.0.1 // indirect
	github.com/admpub/bart v0.0.1 // indirect
	github.com/admpub/boltstore v1.2.0 // indirect
	github.com/admpub/caddy v1.2.9 // indirect
	github.com/admpub/captcha-go v0.0.1 // indirect
	github.com/admpub/ccs-gm v0.0.5 // indirect
	github.com/admpub/checksum v1.1.0 // indirect
	github.com/admpub/collate v1.1.0 // indirect
	github.com/admpub/confl v0.2.4 // indirect
	github.com/admpub/cove v0.0.0-20241224063114-4fdd53c948a6 // indirect
	github.com/admpub/cron v0.1.1 // indirect
	github.com/admpub/dgoogauth v0.0.1 // indirect
	github.com/admpub/email v2.4.1+incompatible // indirect
	github.com/admpub/fasthttp v0.0.7 // indirect
	github.com/admpub/fsnotify v1.7.0 // indirect
	github.com/admpub/gifresize v1.0.2 // indirect
	github.com/admpub/go-bindata-assetfs v0.0.1 // indirect
	github.com/admpub/go-captcha-assets v0.0.0-20250122071745-baa7da4bda0d // indirect
	github.com/admpub/go-captcha/v2 v2.0.7 // indirect
	github.com/admpub/go-figure v0.0.2 // indirect
	github.com/admpub/go-isatty v0.0.11 // indirect
	github.com/admpub/go-lock v1.3.0 // indirect
	github.com/admpub/go-password v0.1.3 // indirect
	github.com/admpub/go-pretty/v6 v6.0.4 // indirect
	github.com/admpub/go-ps v0.0.1 // indirect
	github.com/admpub/go-reuseport v0.5.0 // indirect
	github.com/admpub/go-utility v0.0.1 // indirect
	github.com/admpub/godownloader v2.2.2+incompatible // indirect
	github.com/admpub/gohls v1.3.3 // indirect
	github.com/admpub/gohls-server v0.3.10 // indirect
	github.com/admpub/gotwilio v0.0.1 // indirect
	github.com/admpub/httpscerts v0.0.0-20180907121630-a2990e2af45c // indirect
	github.com/admpub/humanize v0.0.0-20190501023926-5f826e92c8ca // indirect
	github.com/admpub/i18n v0.4.6 // indirect
	github.com/admpub/identicon v1.0.2 // indirect
	github.com/admpub/imaging v1.6.3 // indirect
	github.com/admpub/ini v1.38.2 // indirect
	github.com/admpub/ip2region/v2 v2.0.1 // indirect
	github.com/admpub/json5 v0.0.1 // indirect
	github.com/admpub/mahonia v0.0.0-20151019004008-c528b747d92d // indirect
	github.com/admpub/mail v0.0.0-20170408110349-d63147b0317b // indirect
	github.com/admpub/map2struct v0.1.3 // indirect
	github.com/admpub/mysql-schema-sync v0.2.6 // indirect
	github.com/admpub/oauth2/v4 v4.0.3 // indirect
	github.com/admpub/pester v0.0.0-20200411024648-005672a2bd48 // indirect
	github.com/admpub/qrcode v0.0.3 // indirect
	github.com/admpub/realip v0.2.7 // indirect
	github.com/admpub/redistore v1.2.2 // indirect
	github.com/admpub/regexp2 v1.1.8 // indirect
	github.com/admpub/safesvg v0.0.8 // indirect
	github.com/admpub/securecookie v1.3.0 // indirect
	github.com/admpub/service v0.0.7 // indirect
	github.com/admpub/sockjs-go/v3 v3.0.1 // indirect
	github.com/admpub/sonyflake v0.0.1 // indirect
	github.com/admpub/tail v1.1.1 // indirect
	github.com/admpub/timeago v1.2.2 // indirect
	github.com/admpub/websocket v1.0.4 // indirect
	github.com/admpub/xencoding v0.0.3 // indirect
	github.com/aliyun/alibaba-cloud-sdk-go v1.63.100 // indirect
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/aws/aws-sdk-go v1.55.7 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/boombuler/barcode v1.1.0 // indirect
	github.com/bradfitz/gomemcache v0.0.0-20250403215159-8d39553ac7cf // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/coscms/captcha v0.2.3 // indirect
	github.com/coscms/go-imgparse v0.0.1 // indirect
	github.com/coscms/session-boltstore v0.0.0-20250617034717-a58d8848fa61 // indirect
	github.com/coscms/session-mysqlstore v0.0.0-20250617035706-a20b648443b1 // indirect
	github.com/coscms/session-redisstore v0.0.0-20250624032337-117cf04cfaf4 // indirect
	github.com/coscms/session-sqlitestore v0.0.4 // indirect
	github.com/coscms/session-sqlstore v0.0.1 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/davidbyttow/govips/v2 v2.15.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dsoprea/go-logging v0.0.0-20200710184922-b02d349568dd // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/ebitengine/purego v0.8.4 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/fcjr/aia-transport-go v1.3.0 // indirect
	github.com/flynn/go-shlex v0.0.0-20150515145356-3f9db97f8568 // indirect
	github.com/francoispqt/gojay v1.2.13 // indirect
	github.com/fynelabs/selfupdate v0.2.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.9 // indirect
	github.com/geoffgarside/ber v1.2.0 // indirect
	github.com/glebarez/go-sqlite v1.22.0 // indirect
	github.com/go-acme/lego/v4 v4.25.1 // indirect
	github.com/go-errors/errors v1.5.1 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/go-jose/go-jose/v4 v4.1.1 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/spec v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.27.0 // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/groupcache v0.0.0-20241129210726-2c02b8208cf8 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/gomodule/redigo v1.9.2 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/css v1.0.1 // indirect
	github.com/gorilla/schema v1.4.1 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/grafov/m3u8 v0.12.1 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/h2non/filetype v1.1.3 // indirect
	github.com/h2non/go-is-svg v0.0.0-20160927212452-35e8c4b0612c // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-syslog v1.0.0 // indirect
	github.com/hashicorp/go-version v1.7.0 // indirect
	github.com/hibiken/asynq v0.25.1 // indirect
	github.com/hirochachacha/go-smb2 v1.1.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jlaffaye/ftp v0.2.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/kisielk/errcheck v1.9.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/klauspost/cpuid v1.3.1 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20250317134145-8bc96cf8fc35 // indirect
	github.com/mailru/easyjson v0.9.0 // indirect
	github.com/markbates/going v1.0.3 // indirect
	github.com/martinlindhe/base36 v1.1.1 // indirect
	github.com/maruel/rs v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mattn/go-sqlite3 v1.14.30 // indirect
	github.com/microcosm-cc/bluemonday v1.0.27 // indirect
	github.com/miekg/dns v1.1.67 // indirect
	github.com/minio/crc64nvme v1.1.0 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/minio/minio-go/v7 v7.0.95 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/muesli/smartcrop v0.3.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/nging-plugins/dlmanager v1.8.2 // indirect
	github.com/opentracing/opentracing-go v1.2.1-0.20220228012449-10b1cf09e00b // indirect
	github.com/oschwald/maxminddb-golang v1.13.1 // indirect
	github.com/pcpl2/go-webp v0.0.1 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/phuslu/iploc v1.0.20250715 // indirect
	github.com/phuslu/lru v1.0.18 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/power-devops/perfstat v0.0.0-20240221224432-82ca36839d55 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.65.0 // indirect
	github.com/prometheus/procfs v0.17.0 // indirect
	github.com/quic-go/qpack v0.5.1 // indirect
	github.com/quic-go/quic-go v0.54.0 // indirect
	github.com/redis/go-redis/v9 v9.11.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/rs/xid v1.6.0 // indirect
	github.com/russross/blackfriday v1.6.0 // indirect
	github.com/rwcarlsen/goexif v0.0.0-20190401172101-9e8deecbddbd // indirect
	github.com/segmentio/fasthash v1.0.3 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/shirou/gopsutil/v4 v4.25.6 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/spf13/cast v1.9.2 // indirect
	github.com/spf13/pflag v1.0.7 // indirect
	github.com/studio-b12/gowebdav v0.10.0 // indirect
	github.com/syndtr/goleveldb v1.0.0 // indirect
	github.com/tdewolff/minify/v2 v2.23.10 // indirect
	github.com/tdewolff/parse/v2 v2.8.1 // indirect
	github.com/tidwall/btree v1.8.0 // indirect
	github.com/tidwall/buntdb v1.3.2 // indirect
	github.com/tidwall/gjson v1.18.0 // indirect
	github.com/tidwall/grect v0.1.4 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/rtred v0.1.2 // indirect
	github.com/tidwall/tinyqueue v0.1.1 // indirect
	github.com/tinylib/msgp v1.3.0 // indirect
	github.com/tklauser/go-sysconf v0.3.15 // indirect
	github.com/tklauser/numcpus v0.10.0 // indirect
	github.com/tuotoo/qrcode v0.0.0-20220425170535-52ccc2bebf5d // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/webx-top/captcha v0.1.0 // indirect
	github.com/webx-top/chardet v0.0.2 // indirect
	github.com/webx-top/codec v0.3.0 // indirect
	github.com/webx-top/poolx v0.0.0-20210912044716-5cfa2d58e380 // indirect
	github.com/webx-top/restyclient v0.0.6 // indirect
	github.com/webx-top/tagfast v0.0.1 // indirect
	github.com/webx-top/validator v0.3.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.etcd.io/bbolt v1.4.2 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/mock v0.5.2 // indirect
	golang.org/x/crypto v0.40.0 // indirect
	golang.org/x/exp v0.0.0-20250718183923-645b1fa84792 // indirect
	golang.org/x/image v0.29.0 // indirect
	golang.org/x/lint v0.0.0-20241112194109-818c5a804067 // indirect
	golang.org/x/mod v0.26.0 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	golang.org/x/time v0.12.0 // indirect
	golang.org/x/tools v0.35.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	modernc.org/libc v1.66.6 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.11.0 // indirect
	modernc.org/sqlite v1.38.2 // indirect
	rsc.io/qr v0.2.0 // indirect
)
