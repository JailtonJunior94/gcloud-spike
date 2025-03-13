package drive

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type upload struct {
	service *drive.Service
}

const prefixLog = "upload: "

func NewUpload(credentials []byte) (*upload, error) {
	service, err := drive.NewService(context.Background(), option.WithCredentialsJSON(credentials))
	if err != nil {
		return nil, fmt.Errorf("%s%v", prefixLog, err)
	}
	return &upload{service: service}, nil
}

func (u *upload) Upload() error {
	files, err := u.service.Files.List().Do()
	if err != nil {
		return fmt.Errorf("%s%v", prefixLog, err)
	}

	for _, file := range files.Files {
		log.Printf("%sfile: %s, id: %s\n", prefixLog, file.Name, file.Id)
	}

	file := &drive.File{
		Name:     "file.json",
		MimeType: "application/json",
		Parents:  []string{"18Btpom3U6GJj3ZAzxgPaW3Fm16ITNxTd"},
	}

	driveFile, err := u.service.Files.Create(file).Media(nil).Do()
	if err != nil {
		return fmt.Errorf("%s%v", prefixLog, err)
	}

	permission := &drive.Permission{
		Role: "writer",
		Type: "anyone",
	}

	permissions, err := u.service.Permissions.Create(driveFile.Id, permission).Do()
	if err != nil {
		return fmt.Errorf("%s%v", prefixLog, err)
	}

	log.Println(permissions)
	return nil
}
