package model

import (
	"net/url"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID        string // UUID
	Title     string
	FileName  string
	Size      int
	MimeType  string
	CreatedAt time.Time

	User *User

	baseURL string
}

func NewFile() (*File, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	return &File{ID: id.String()}, nil
}

func (f *File) SetBaseURL(base string) {
	f.baseURL = base
}

func (f *File) SetUser(user *User) {
	f.User = user
}

func (f *File) GetDownloadURL() (string, error) {
	ext := filepath.Ext(f.FileName)
	endpoint, err := url.JoinPath(f.baseURL, f.ID+ext)
	if err != nil {
		return "", err
	}
	return endpoint, nil
}
