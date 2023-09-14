package sources

type Filter struct {
	Limit uint
	Skip  uint
	Id    uint
}

type Source interface {
	Get(filter *Filter, callback func(rows any, err error)) (any, error)
	Count() uint
}
