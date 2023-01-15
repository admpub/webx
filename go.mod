module github.com/admpub/webx

go 1.17

replace github.com/admpub/nging/v5 => ../../../github.com/admpub/nging

replace github.com/nging-plugins/dbmanager => ../../nging-plugins/dbmanager

replace github.com/nging-plugins/dlmanager => ../../nging-plugins/dlmanager

exclude github.com/gomodule/redigo v2.0.0+incompatible

require (
	github.com/RichardKnop/machinery v1.10.6
	github.com/adamzy/cedar-go v0.0.0-20170805034717-80a9c64b256d // indirect
	github.com/admpub/bindata/v3 v3.1.5
	github.com/admpub/cache v0.3.8
	github.com/admpub/color v1.8.1
	github.com/admpub/copier v0.1.1
	github.com/admpub/decimal v1.3.1
	github.com/admpub/errors v0.8.2
	github.com/admpub/fasttemplate v0.0.2
	github.com/admpub/go-download/v2 v2.1.12
	github.com/admpub/go-hashids v2.0.1+incompatible
	github.com/admpub/go-shortid v0.0.0-20140827050853-24d054c393fe
	github.com/admpub/godotenv v1.4.2
	github.com/admpub/imageproxy v0.9.3
	github.com/admpub/ipfilter v1.0.6
	github.com/admpub/license_gen v0.1.0
	github.com/admpub/log v1.3.3
	github.com/admpub/marmot v0.0.0-20200702042226-2170d9ff59f5
	github.com/admpub/null v8.0.4+incompatible
	github.com/admpub/once v0.0.1
	github.com/admpub/pinyin-golang v1.0.1
	github.com/admpub/pp v0.0.5
	github.com/admpub/redsync/v4 v4.0.3
	github.com/admpub/resty/v2 v2.7.0
	github.com/admpub/sensitive v0.0.0-20200106142752-42d1c505be7b
	github.com/admpub/useragent v0.0.1
	github.com/caddy-plugins/ipfilter v1.1.8
	github.com/chai2010/webp v1.1.1 // indirect
	github.com/coscms/forms v1.12.1
	github.com/coscms/oauth2s v0.2.1
	github.com/coscms/sms v0.0.4
	github.com/googollee/go-socket.io v1.6.2
	github.com/gosimple/slug v1.13.1
	github.com/hibiken/asynq v0.23.0
	github.com/huichen/sego v0.0.0-20210824061530-c87651ea5c76
	github.com/markbates/goth v1.76.0
	github.com/prometheus/client_golang v1.14.0
	github.com/spf13/cobra v1.6.1
	github.com/stretchr/testify v1.8.1
	github.com/swaggo/swag v1.8.8
	github.com/wangbin/jiebago v0.3.2
	github.com/webx-top/client v0.8.8
	github.com/webx-top/com v0.9.1
	github.com/webx-top/db v1.23.13
	github.com/webx-top/echo v2.33.1+incompatible
	github.com/webx-top/echo-prometheus v1.1.0
	github.com/webx-top/echo-socket.io v1.1.2
	github.com/webx-top/image v0.0.9
	github.com/webx-top/pagination v0.2.1
	github.com/webx-top/validation v0.0.3 // indirect
	github.com/yanyiwu/gojieba v1.2.0
	golang.org/x/oauth2 v0.4.0
	golang.org/x/sync v0.1.0
	gopkg.in/redis.v5 v5.2.9
)

require (
	github.com/admpub/confl v0.2.2
	github.com/admpub/events v1.3.5
	github.com/admpub/go-bindata-assetfs v0.0.0-20170428090253-36eaa4c19588
	github.com/admpub/go-lock v1.3.0
	github.com/admpub/go-zinc v0.0.8
	github.com/admpub/nging/v5 v5.0.0
	github.com/golang-jwt/jwt/v4 v4.4.3
	github.com/martinlindhe/base36 v1.1.1
	github.com/meilisearch/meilisearch-go v0.21.1
	github.com/nging-plugins/dbmanager v1.1.5
	github.com/webx-top/validator v0.2.0
)

require (
	cloud.google.com/go v0.107.0 // indirect
	cloud.google.com/go/compute v1.14.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/iam v0.8.0 // indirect
	cloud.google.com/go/pubsub v1.28.0 // indirect
	gitee.com/admpub/certmagic v0.8.8 // indirect
	github.com/KenmyZhang/aliyun-communicate v0.0.0-20180308134849-7997edc57454 // indirect
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/RichardKnop/logging v0.0.0-20190827224416-1a693bdd4fae // indirect
	github.com/abh/errorutil v0.0.0-20130729183701-f9bd360d00b9 // indirect
	github.com/admpub/9t v0.0.0-20190605154903-a68069ace5e1 // indirect
	github.com/admpub/archiver v1.1.4 // indirect
	github.com/admpub/caddy v1.1.13 // indirect
	github.com/admpub/ccs-gm v0.0.3 // indirect
	github.com/admpub/checksum v1.1.0 // indirect
	github.com/admpub/cron v0.0.1 // indirect
	github.com/admpub/dgoogauth v0.0.1 // indirect
	github.com/admpub/email v2.4.1+incompatible // indirect
	github.com/admpub/fasthttp v0.0.5 // indirect
	github.com/admpub/fsnotify v1.5.0 // indirect
	github.com/admpub/gifresize v1.0.2 // indirect
	github.com/admpub/go-figure v0.0.0-20180619031829-18b2b544842c // indirect
	github.com/admpub/go-isatty v0.0.10 // indirect
	github.com/admpub/go-password v0.1.3 // indirect
	github.com/admpub/go-pretty/v6 v6.0.3 // indirect
	github.com/admpub/go-reuseport v0.0.4 // indirect
	github.com/admpub/go-utility v0.0.1 // indirect
	github.com/admpub/godownloader v2.2.0+incompatible // indirect
	github.com/admpub/gohls v1.3.3 // indirect
	github.com/admpub/gohls-server v0.3.7 // indirect
	github.com/admpub/gotwilio v0.0.0-20210910030032-9a691aeea447 // indirect
	github.com/admpub/httpscerts v0.0.0-20180907121630-a2990e2af45c // indirect
	github.com/admpub/humanize v0.0.0-20190501023926-5f826e92c8ca // indirect
	github.com/admpub/i18n v0.1.3 // indirect
	github.com/admpub/identicon v1.0.2 // indirect
	github.com/admpub/imaging v1.5.0 // indirect
	github.com/admpub/ini v1.38.2 // indirect
	github.com/admpub/ip2region/v2 v2.0.1 // indirect
	github.com/admpub/json5 v0.0.1 // indirect
	github.com/admpub/mahonia v0.0.0-20151019004008-c528b747d92d // indirect
	github.com/admpub/mail v0.0.0-20170408110349-d63147b0317b // indirect
	github.com/admpub/mysql-schema-sync v0.2.1 // indirect
	github.com/admpub/pester v0.0.0-20200411024648-005672a2bd48 // indirect
	github.com/admpub/qrcode v0.0.3 // indirect
	github.com/admpub/randomize v0.0.2 // indirect
	github.com/admpub/realip v0.0.0-20210421084339-374cf5df122d // indirect
	github.com/admpub/redistore v1.2.1 // indirect
	github.com/admpub/securecookie v1.1.2 // indirect
	github.com/admpub/service v0.0.4 // indirect
	github.com/admpub/sessions v0.1.3 // indirect
	github.com/admpub/sockjs-go/v3 v3.0.1 // indirect
	github.com/admpub/sonyflake v0.0.1 // indirect
	github.com/admpub/tail v1.1.0 // indirect
	github.com/admpub/timeago v1.2.1 // indirect
	github.com/admpub/websocket v1.0.4 // indirect
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/aws/aws-sdk-go v1.44.180 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bits-and-blooms/bitset v1.2.2 // indirect
	github.com/boombuler/barcode v1.0.1 // indirect
	github.com/bradfitz/gomemcache v0.0.0-20221031212613-62deef7fc822 // indirect
	github.com/cenkalti/backoff/v4 v4.2.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/coscms/go-imgparse v0.0.0-20150925144422-3e3a099f7856 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/disintegration/imaging v1.6.2 // indirect
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/dsoprea/go-logging v0.0.0-20200710184922-b02d349568dd // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/fcjr/aia-transport-go v1.2.2 // indirect
	github.com/flynn/go-shlex v0.0.0-20150515145356-3f9db97f8568 // indirect
	github.com/francoispqt/gojay v1.2.13 // indirect
	github.com/friendsofgo/errors v0.9.2 // indirect
	github.com/garyburd/redigo v1.6.4 // indirect
	github.com/go-acme/lego/v4 v4.9.1 // indirect
	github.com/go-errors/errors v1.4.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/spec v0.20.7 // indirect
	github.com/go-openapi/swag v0.22.3 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.11.1 // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/go-redsync/redsync/v4 v4.7.1 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/go-task/slim-sprig v0.0.0-20210107165309-348f09dbbbc0 // indirect
	github.com/goccy/go-json v0.10.0 // indirect
	github.com/gofrs/uuid v4.3.1+incompatible // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gomodule/redigo v1.8.9 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/pprof v0.0.0-20230111200839-76d1ae5aea2b // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.1 // indirect
	github.com/googleapis/gax-go/v2 v2.7.0 // indirect
	github.com/gorilla/css v1.0.0 // indirect
	github.com/gorilla/schema v1.2.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/grafov/m3u8 v0.11.1 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/h2non/filetype v1.1.3 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-syslog v1.0.0 // indirect
	github.com/imdario/mergo v0.3.13 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/kisielk/errcheck v1.6.3 // indirect
	github.com/klauspost/compress v1.15.14 // indirect
	github.com/klauspost/cpuid v1.3.1 // indirect
	github.com/klauspost/cpuid/v2 v2.2.3 // indirect
	github.com/klauspost/pgzip v1.2.5 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/lucas-clemente/quic-go v0.31.1 // indirect
	github.com/lufia/plan9stats v0.0.0-20230110061619-bbe2e5e100de // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/markbates/going v1.0.3 // indirect
	github.com/marten-seemann/qpack v0.3.0 // indirect
	github.com/marten-seemann/qtls-go1-18 v0.1.4 // indirect
	github.com/marten-seemann/qtls-go1-19 v0.1.2 // indirect
	github.com/maruel/rs v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/mattn/go-sqlite3 v1.14.16 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/microcosm-cc/bluemonday v1.0.21 // indirect
	github.com/miekg/dns v1.1.50 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/minio/minio-go/v7 v7.0.47 // indirect
	github.com/minio/sha256-simd v1.0.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/montanaflynn/stats v0.6.6 // indirect
	github.com/mrjones/oauth v0.0.0-20190623134757-126b35219450 // indirect
	github.com/muesli/smartcrop v0.3.0 // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/nging-plugins/dlmanager v1.1.1 // indirect
	github.com/nwaples/rardecode v1.1.3 // indirect
	github.com/onsi/ginkgo/v2 v2.7.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/oschwald/maxminddb-golang v1.10.0 // indirect
	github.com/pborman/uuid v1.2.1 // indirect
	github.com/phuslu/iploc v1.0.20221130 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/power-devops/perfstat v0.0.0-20221212215047-62379fc7944b // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.39.0 // indirect
	github.com/prometheus/procfs v0.9.0 // indirect
	github.com/rivo/uniseg v0.4.3 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/rs/xid v1.4.0 // indirect
	github.com/russross/blackfriday v1.6.0 // indirect
	github.com/rwcarlsen/goexif v0.0.0-20190401172101-9e8deecbddbd // indirect
	github.com/segmentio/fasthash v1.0.3 // indirect
	github.com/shirou/gopsutil/v3 v3.22.12 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/smartwalle/crypto4go v1.0.3 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/streadway/amqp v1.0.0 // indirect
	github.com/syndtr/goleveldb v1.0.0 // indirect
	github.com/tklauser/go-sysconf v0.3.11 // indirect
	github.com/tklauser/numcpus v0.6.0 // indirect
	github.com/tuotoo/qrcode v0.0.0-20220425170535-52ccc2bebf5d // indirect
	github.com/ulikunitz/xz v0.5.11 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.43.0 // indirect
	github.com/volatiletech/inflect v0.0.1 // indirect
	github.com/volatiletech/strmangle v0.0.4 // indirect
	github.com/webx-top/captcha v0.1.0 // indirect
	github.com/webx-top/chardet v0.0.1 // indirect
	github.com/webx-top/codec v0.2.1 // indirect
	github.com/webx-top/poolx v0.0.0-20210912044716-5cfa2d58e380 // indirect
	github.com/webx-top/restyclient v0.0.3 // indirect
	github.com/webx-top/tagfast v0.0.0-20161020041435-9a2065ce3dd2 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	go.mongodb.org/mongo-driver v1.11.1 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/crypto v0.5.0 // indirect
	golang.org/x/exp v0.0.0-20230113213754-f9f960f08ad4 // indirect
	golang.org/x/image v0.3.0 // indirect
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/mod v0.7.0 // indirect
	golang.org/x/net v0.5.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	golang.org/x/tools v0.5.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/api v0.107.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230113154510-dbe35b8444a5 // indirect
	google.golang.org/grpc v1.52.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	rsc.io/qr v0.2.0 // indirect
)
