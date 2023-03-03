package id

type Generater interface {
	Generate() ID
}

type ID interface {
	Int64() int64
	String() string
}


func New(nodeId int64) (Generater, error) {
	return NewTwitterSF(nodeId)
}

