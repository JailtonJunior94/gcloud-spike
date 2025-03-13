package drive

import (
	"bytes"
	"context"
	"encoding/csv"
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

func (u *upload) Upload(fileID string) error {
	file, err := u.service.Files.Get(fileID).Do()
	if err != nil {
		return fmt.Errorf("%s%v", prefixLog, err)
	}

	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)

	records := [][]string{
		{"Name", "Age"},
		{"Jailton", "31"},
		{"Stefany", "30"},
		{"Antony", "04"},
		{"Noah", "00"},
	}

	if err := writer.WriteAll(records); err != nil {
		return fmt.Errorf("%s%v", prefixLog, err)
	}

	gDriveFile := &drive.File{
		Name:     "users.csv",
		MimeType: "text/csv",
		Parents:  []string{file.Id},
	}

	fileCreated, err := u.service.Files.Create(gDriveFile).Media(&buffer).Do()
	if err != nil {
		return fmt.Errorf("%s%v", prefixLog, err)
	}
	log.Printf("%sfile '%s' uploaded in '%s' folder\n", prefixLog, fileCreated.Name, file.Name)

	permission := &drive.Permission{Role: "writer", Type: "anyone"}
	permissions, err := u.service.Permissions.Create(fileCreated.Id, permission).Do()
	if err != nil {
		return fmt.Errorf("%s%v", prefixLog, err)
	}

	log.Printf("%sfile '%s' shared with link: %s\n", prefixLog, fileCreated.Name, permissions.DisplayName)
	return nil
}
