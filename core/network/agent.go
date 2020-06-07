package ynet

type Agent interface {
	Run()
	OnClose()
}
