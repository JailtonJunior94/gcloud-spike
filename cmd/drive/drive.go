package drive

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"strconv"

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

	if err := writer.WriteAll(generateRecords(1_00000)); err != nil {
		return fmt.Errorf("%s%v", prefixLog, err)
	}

	gDriveFile := &drive.File{
		Name:     "customers.csv",
		MimeType: "text/csv",
		Parents:  []string{file.Id},
	}

	fileCreated, err := u.service.Files.Create(gDriveFile).Media(&buffer).Do()
	if err != nil {
		return fmt.Errorf("%s%v", prefixLog, err)
	}

	permission := &drive.Permission{Role: "writer", Type: "anyone"}
	_, err = u.service.Permissions.Create(fileCreated.Id, permission).Do()
	if err != nil {
		return fmt.Errorf("%s%v", prefixLog, err)
	}

	log.Printf("%sfile '%s' uploaded in '%s' folder\n", prefixLog, fileCreated.Name, file.Name)
	return nil
}

func generateRecords(numRecords int) [][]string {
	records := [][]string{{"Name", "Age"}}
	for i := 1; i <= numRecords; i++ {
		records = append(records, []string{fmt.Sprintf("Name%d", i), strconv.Itoa(20 + (i % 30))})
	}
	return records
}
