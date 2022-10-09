package top

import (
	"testing"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/testing/test"
)

func TestStars(t *testing.T) {
	test.Eq(t, STAR_HALF, StarsSlice(0.5)[0])
	test.Eq(t, STAR_HALF, StarsSlice(7.5)[7])
	echo.Dump(StarsSlicex(7.5, 10, 5))
	test.Eq(t, STAR_HALF, StarsSlicex(7.5, 10, 5)[3])
}

func TestParseDuration(t *testing.T) {
	dur, _ := ParseDuration(`2d`)
	test.Eq(t, DurationDay * 2, dur)
	dur, _ = ParseDuration(`2mo`)
	test.Eq(t, DurationMonth * 2, dur)
	dur, _ = ParseDuration(`2y`)
	test.Eq(t, DurationYear * 2, dur)
	dur, _ = ParseDuration(`2w`)
	test.Eq(t, DurationWeek * 2, dur)
}
