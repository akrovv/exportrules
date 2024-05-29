package valid

type bio struct {
	s string
}

type Fio struct {
	s string
}

func NewB() bio {
	return bio{}
}

type name struct {
	age int
	bio bio
	fio Fio
}

func NewN() name {
	n := name{
		age: 1,
		bio: NewB(),
	}

	return n
}

func NewNN() name {
	return name{
		age: 1,
		bio: NewB(),
	}
}

func NewNNN() name {
	n := name{
		age: 1,
	}

	n.bio = NewB()

	return n
}

func NewNNNN() name {
	return name{
		age: 1,
		fio: Fio{},
	}
}
