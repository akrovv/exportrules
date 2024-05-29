package valid

type People struct {
	Population float64
}

func (p People) Less(o float64) bool {
	return p.Population < o
}

var p = People{}

func setPopulation() {
	p.Population = 1000
}

func getPopulation() float64 {
	return p.Population
}

func getLess() {
	p.Less(10)
}
