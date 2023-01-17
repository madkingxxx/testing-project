package file_processing

import "file_processing_service/internal/entity"

type Repository interface {
	Create(file *entity.File) error
	Update(file *entity.File) error
	Get(name string) (*entity.File, error)
	List() ([]*entity.File, error)
}

type FileUsecase interface {
	UploadFile(file *entity.File) error
	DownloadFile(name string) (*entity.File, error)
	GetFileList() ([]*entity.File, error)
}
