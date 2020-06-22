package main

type DirectoryStore struct {
	root string
}

func (ds *DirectoryStore) GetListSubscribers(list List) []Subscriber {
	return []Subscriber{}
}
