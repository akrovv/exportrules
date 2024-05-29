package invalid

func NewM() Men {
	return Men{l: 10}
}

func New() Women {
	return Women{l: 10}
}

type Men struct {
	l int
}

type Women struct {
	l int
}

func (m Men) Men() int {
	w := Women{l: 10}
	w.l++
	w.l = 10
	return w.l
}
