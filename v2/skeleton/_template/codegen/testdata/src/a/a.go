package a

type DB interface {
	Get(id string) int
	Set(id string, v int)
}

type db struct{}

func (db) Get(id string) int    { return 0 }
func (db) Set(id string, v int) {}

type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}
