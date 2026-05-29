package store

type Factory interface {
	Close() error
}

var client Factory

func Set(f Factory) {
	client = f
}

func Get() Factory {
	return client
}
