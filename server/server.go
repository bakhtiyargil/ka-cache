package server

type Server interface {
	Start()
	Stop()
	Running() bool
}
