package version

import (
	"embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed testembed
var f embed.FS

//go:embed testembed/*
var f2 embed.FS

//go:embed testembed/.test.log
var f3 embed.FS

func TestEmbed(t *testing.T) {
	t.Logf(`------- embed folder -------`)
	dirs, err := f.ReadDir(`testembed`)
	assert.NoError(t, err)
	expected := []string{`test.md`}
	var actual []string
	for _, dir := range dirs {
		t.Logf(dir.Name())
		actual = append(actual, dir.Name())
	}
	assert.Equal(t, expected, actual)

	t.Logf(`------- embed folder/* -------`)
	dirs, err = f2.ReadDir(`testembed`)
	assert.NoError(t, err)
	expected = []string{".test.log", "_test.txt", "test.md"}
	actual = actual[0:0]
	for _, dir := range dirs {
		t.Logf(dir.Name())
		actual = append(actual, dir.Name())
	}
	assert.Equal(t, expected, actual)

	t.Logf(`------- embed folder/file -------`)
	dirs, err = f3.ReadDir(`testembed`)
	assert.NoError(t, err)
	expected = []string{`.test.log`}
	actual = actual[0:0]
	for _, dir := range dirs {
		t.Logf(dir.Name())
		actual = append(actual, dir.Name())
	}
	assert.Equal(t, expected, actual)
}

/*
* //go:embed images
* var images embed.FS // 不包含.b.jpg和_photo_metadata目录
*
* //go:embed images/*
* var images embed.FS // 注意！！！ 这里包含.b.jpg和_photo_metadata目录
*
* //go:embed images/.b.jog
* var bJPG []byte // 明确给出文件名也不会被忽略
 */
