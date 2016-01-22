package storage

type Store map[string][]byte

func MakeStore() *Store {
	s := make(Store, 1024)
	return &s
}

func (s *Store) Get(key string) []byte {
	return (*s)[key]
}

func (s *Store) Put(key string, val []byte) {
	(*s)[key] = val
}
