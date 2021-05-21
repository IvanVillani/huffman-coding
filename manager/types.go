package manager

type ValueCodeLenght struct {
	Value           byte
	ValueCodeLenght int
}

type ValuesCodeLengthsArray struct {
	ValuesCodeLenghts []ValueCodeLenght
	Frequency         int
}

type Values []byte

func (v Values) Len() int {
	return len(v)
}

func (v Values) Less(i, j int) bool {
	return v[i] < v[j]
}

func (v Values) Swap(i, j int) {
	temp := v[i]
	v[i] = v[j]
	v[j] = temp
}
