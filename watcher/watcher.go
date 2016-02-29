package watcher

type Watcher interface {
	Start() <-chan bool
	Stop()
}
