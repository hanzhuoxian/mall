package cache

type Factory interface {
	User() UserCache
}

var cache Factory

func Set(f Factory) {
	cache = f
}

func Get() Factory {
	return cache
}
