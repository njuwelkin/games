package mkf

type DataMkf struct {
	Mkf
}

func NewDataMkf(path string) (DataMkf, error) {
	ret := DataMkf{Mkf{}}
	return ret, ret.Open(path)
}
