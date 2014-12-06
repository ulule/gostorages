package storages

type Storage interface {
	Save(filepath string, content []byte) (string, error)
	Path(filepath string) string
	Exists(filepath string) bool
	Delete(filepath string) error
	Open(filepath string) ([]byte, error)
	Params(params map[string]string) error
}
