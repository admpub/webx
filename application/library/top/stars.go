package top

import "html/template"

type VSTAR int

func (v VSTAR) IconStar() string {
	if v, ok := StarIcons[v]; ok {
		return v
	}
	return DefaultStarIcon
}

func (v VSTAR) String() string {
	if v, ok := StarNames[v]; ok {
		return v
	}
	return DefaultStarName
}

const (
	STAR_EMPTY VSTAR = 0
	STAR_FULL  VSTAR = 1
	STAR_HALF  VSTAR = 2
)

var (
	DefaultStarIcon = `fa-star-o`
	StarIcons       = map[VSTAR]string{
		STAR_EMPTY: `fa-star-o`,
		STAR_FULL:  `fa-star`,
		STAR_HALF:  `fa-star-half-o`,
	}
	DefaultStarName = `empty`
	StarNames       = map[VSTAR]string{
		STAR_EMPTY: `empty`,
		STAR_FULL:  `full`,
		STAR_HALF:  `half`,
	}
)

type Stars []VSTAR

func (s Stars) HTML() template.HTML {
	var v string
	for _, star := range s {
		v += `<i class="fa ` + star.IconStar() + `"></i>`
	}
	return template.HTML(v)
}

// StarsSlice 生成星标切片
func StarsSlice(cur float64, max ...int) Stars {
	m := 10
	if len(max) > 0 {
		m = max[0]
	}
	r := make([]VSTAR, m)
	n := cur
	for k := range r {
		if float64(k) < n {
			r[k] = STAR_FULL
			if float64(k+1) > n {
				r[k] = STAR_HALF
			}
			continue
		}
		break
	}
	return r
}

func StarsSlice5(cur float64) Stars {
	return StarsSlice(cur, 5)
}

func StarsSlicex5(curRating float64, maxRating float64) Stars {
	return StarsSlicex(curRating, maxRating, 5)
}

// StarsSlicex 将最大分数按照最大星标数量平分生成星标切片
func StarsSlicex(curRating float64, maxRating float64, maxStars ...int) Stars {
	m := 10
	if len(maxStars) > 0 {
		m = maxStars[0]
	}
	r := make([]VSTAR, m)
	n := curRating
	p := maxRating / float64(m)
	for k := range r {
		rating := float64(k) * p
		if rating < n {
			r[k] = STAR_FULL
			if rating+p > n {
				r[k] = STAR_HALF
			}
			continue
		}
		break
	}
	return r
}
