package sources

type Config interface {
}

type Source interface {
	Size() int64
	GetRandomImage() (string, string, error)
	GetConfig() Config
}
