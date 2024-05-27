package negative

type bio struct {
	s string
}

func NewB() bio {
	return bio{}
}

type name struct {
	age int
	bio bio
}

func NewN() name {
	n := name{
		age: 1,
		bio: bio{},
	}

	return n
}

func NewNN() name {
	return name{
		age: 1,
		bio: bio{},
	}
}

func NewNNN() name {
	n := name{
		age: 1,
	}

	n.bio = bio{}

	return n
}
