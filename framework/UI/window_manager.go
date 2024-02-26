package ui

type WinManager struct {
	winStack []*Window
}

func NewWinManager() WinManager {
	ret := WinManager{
		winStack: []*Window{},
	}
	return ret
}
