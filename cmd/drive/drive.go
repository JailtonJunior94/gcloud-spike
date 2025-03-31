package drive

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type upload struct {
	service *drive.Service
}

const prefixLog = "upload:"

func NewUpload(credentials []byte) (*upload, error) {
	service, err := drive.NewService(context.Background(), option.WithCredentialsJSON(credentials))
	if err != nil {
		return nil, fmt.Errorf("%s%v", prefixLog, err)
	}
	return &upload{service: service}, nil
}

func (u *upload) Upload(folder *drive.File, content bytes.Buffer) error {
	gDriveFile := &drive.File{
		Name:     fmt.Sprintf("customers-%v.csv", time.Now().Format("2006-01-02")),
		MimeType: "text/csv",
		Parents:  []string{folder.Id},
	}

	fileCreated, err := u.service.Files.Create(gDriveFile).Media(&content).Do()
	if err != nil {
		return fmt.Errorf("%s%v", prefixLog, err)
	}

	permission := &drive.Permission{Role: "writer", Type: "anyone"}
	_, err = u.service.Permissions.Create(fileCreated.Id, permission).Do()
	if err != nil {
		return fmt.Errorf("%s%v", prefixLog, err)
	}

	log.Printf("%s file %s created in folder %s\n", prefixLog, fileCreated.Name, folder.Name)
	return nil
}

func (u *upload) GetFolder(name string) (*drive.File, error) {
	query := fmt.Sprintf("name='%s'", name)
	fileList, err := u.service.Files.List().Q(query).Do()
	if err != nil {
		return nil, fmt.Errorf("%s%v", prefixLog, err)
	}

	if len(fileList.Files) == 0 {
		return nil, fmt.Errorf("%s file not found", prefixLog)
	}

	return fileList.Files[0], nil
}

func (u *upload) GetFolderByID(id string) (*drive.File, error) {
	file, err := u.service.Files.Get(id).Do()
	if err != nil {
		return nil, fmt.Errorf("%s%v", prefixLog, err)
	}
	return file, nil
}
