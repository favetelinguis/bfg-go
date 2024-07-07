package bfggo

type Executor interface {
	Place(bet int) (int, error)
}
