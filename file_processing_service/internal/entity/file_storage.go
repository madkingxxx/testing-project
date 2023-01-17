package entity

import (
	"sync"
)

type FileStorage struct {
	sync.Map
}

func NewFileStorage() *FileStorage {
	return &FileStorage{}
}

func (s *FileStorage) Create(file *File) error {
	if _, ok := s.Load(file.Name); ok {
		return ErrAlreadyExists
	}
	s.Store(file.Name, file)
	return nil
}

func (s *FileStorage) Get(name string) (*File, error) {
	if file, ok := s.Load(name); ok {
		return file.(*File), nil
	}
	return nil, ErrNotFound
}

func (s *FileStorage) List() ([]*File, error) {
	var files []*File
	s.Range(func(key, value interface{}) bool {
		files = append(files, value.(*File))
		return true
	})
	return files, nil
}

func (s *FileStorage) Update(file *File) error {
	if _, ok := s.Load(file.Name); !ok {
		return ErrNotFound
	}
	s.Store(file.Name, file)
	return nil
}
