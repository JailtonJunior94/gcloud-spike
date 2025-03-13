package main

import (
	"log"
	"os"

	"github.com/jailtonjunior94/gcloud-spike/cmd/drive"

	"github.com/spf13/cobra"
)

func main() {
	credentials, err := os.ReadFile("../credentials.json")
	if err != nil {
		log.Fatalf("error reading credentials file: %v", err)
	}

	root := &cobra.Command{
		Use:   "gcloud-spikes",
		Short: "gcloud-spikes is a collection of spikes for Google Cloud Platform",
	}

	upload := &cobra.Command{
		Use:   "upload",
		Short: "drive is a spike for Google Drive",
		Run: func(cmd *cobra.Command, args []string) {
			upload, err := drive.NewUpload(credentials)
			if err != nil {
				log.Fatal(err)
			}

			if err := upload.Upload(); err != nil {
				log.Fatal(err)
			}
		},
	}

	chat := &cobra.Command{
		Use:   "chat",
		Short: "chat is a spike for Google Chat",
	}

	root.AddCommand(upload, chat)
	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}
