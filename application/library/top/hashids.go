package top

import (
	"regexp"
	"sync"

	"github.com/webx-top/echo"

	hashids "github.com/admpub/go-hashids"
	"github.com/coscms/webcore/library/common"
	syncOnce "github.com/admpub/once"
	"github.com/martinlindhe/base36"
)

var (
	hashidsCodec         *hashids.HashID
	mutex                = &sync.Mutex{}
	hashidsSalt          string
	hashidsOnce          syncOnce.Once
	snCodec              *hashids.HashID
	hashidsDefaultPrefix = `A`
	snRegexp             = regexp.MustCompile(`^[A-Z0-9]{14,}$`)
)

func init() {
	snCodec = NewSNCodec()
}

func NewSNCodec() *hashids.HashID {
	data := hashids.NewUpperCaseData()
	data.Uint64Mode = true
	codec, err := hashids.NewWithData(data)
	if err != nil {
		panic(err)
	}
	return codec
}

func IsSN(v string) bool {
	return snRegexp.MatchString(v)
}

// GenSN 生成唯一序列号 (经密码加密)
// 生成结果长度一般为 15 个字符
func GenSN(prefixes ...string) (string, error) {
	id, err := common.NextID()
	if err != nil {
		return ``, err
	}
	sn, err := snCodec.EncodeUint64([]uint64{id})
	if err != nil {
		return ``, err
	}
	var prefix string
	if len(prefixes) > 0 {
		prefix = prefixes[0]
	} else {
		prefix = hashidsDefaultPrefix
	}
	return prefix + sn, nil
}

// RawSN 生成唯一序列号 (无加密)
// 生成结果长度一般为 12 个字符
func RawSN() (string, error) {
	id, err := common.NextID()
	if err != nil {
		return ``, err
	}
	return Base36Encode(id), nil
}

func Base36Encode(i uint64) string {
	return base36.Encode(i)
}

func Base36Decode(v string) uint64 {
	return base36.Decode(v)
}

func HashID() *hashids.HashID {
	if cfg, ok := echo.Get(common.ConfigName).(common.ConfigFromDB); ok {
		salt := cfg.ConfigFromDB().GetStore(`base`).String(`hashidSalt`)
		if salt != hashidsSalt {
			mutex.Lock()
			hashidsSalt = salt
			hashidsCodec = NewHashID(hashidsSalt)
			mutex.Unlock()
		}
	}
	if hashidsCodec == nil {
		hashidsCodec = NewHashID(hashidsSalt)
	}
	return hashidsCodec
}

func HashIDOnce() *hashids.HashID {
	hashidsOnce.Do(func() {
		HashID()
	})
	return hashidsCodec
}

func NewHashID(salt string) *hashids.HashID {
	data := hashids.NewData()
	data.Salt = salt
	data.Uint64Mode = true
	codec, err := hashids.NewWithData(data)
	if err != nil {
		panic(err)
	}
	return codec
}

func HashIDEncode(id uint64) (string, error) {
	encoded, err := HashIDOnce().EncodeUint64([]uint64{id})
	return encoded, err
}

func HashIDDecode(encoded string) (uint64, error) {
	ids, err := HashIDOnce().DecodeUint64WithError(encoded)
	if err != nil {
		return 0, err
	}
	return ids[0], nil
}
