package cache

type Backup interface {
	Write(content []byte) error
	Read() ([]byte, error)
	Close() error
	Clean() error
}