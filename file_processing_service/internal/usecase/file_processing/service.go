package file_processing

import (
	"file_processing_service/internal/entity"
	"fmt"
	"io"
	"os"
	"time"
)

const fileStorage = "file_storage"

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) UploadFile(file *entity.File) error {
	f, err := os.Create(fmt.Sprintf("./%s/%s", fileStorage, file.Name))
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, file.Content)
	if err != nil {
		return err
	}
	err = s.repo.Create(file)
	if err == entity.ErrAlreadyExists {
		return s.repo.Update(&entity.File{
			Name:      file.Name,
			CreatedAt: file.CreatedAt,
			UpdatedAt: time.Now().UTC(),
		})
	}
	return err
}

func (s *Service) DownloadFile(name string) (*entity.File, error) {
	file, err := s.repo.Get(name)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(fmt.Sprintf("./%s/%s", fileStorage, file.Name))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	_, err = io.Copy(file.Content, f)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (s *Service) GetFileList() ([]*entity.File, error) {
	return s.repo.List()
}
