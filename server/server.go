package server

type Server interface {
	Start() error
	Running() bool
}
