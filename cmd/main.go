package main

import (
	"encoding/base64"
	"log"

	"github.com/jailtonjunior94/gcloud-spike/cmd/chat"
	"github.com/jailtonjunior94/gcloud-spike/cmd/drive"
	"github.com/jailtonjunior94/gcloud-spike/configs"

	"github.com/spf13/cobra"
)

func main() {
	cfg, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	serviceAccount, err := base64.StdEncoding.DecodeString(cfg.GCloudAPIKey)
	if err != nil {
		log.Fatalf("failed to decode Base64 string: %v", err)
	}

	root := &cobra.Command{
		Use:   "gcloud-spikes",
		Short: "gcloud-spikes is a collection of spikes for Google Cloud Platform",
	}

	upload := &cobra.Command{
		Use:   "upload",
		Short: "drive is a spike for Google Drive",
		Run: func(cmd *cobra.Command, args []string) {
			upload, err := drive.NewUpload(serviceAccount)
			if err != nil {
				log.Fatal(err)
			}

			if err := upload.Upload(cfg.GDriveFolderID); err != nil {
				log.Fatal(err)
			}
		},
	}

	chat := &cobra.Command{
		Use:   "chat",
		Short: "chat is a spike for Google Chat",
		Run: func(cmd *cobra.Command, args []string) {
			chat, err := chat.NewChat(serviceAccount)
			if err != nil {
				log.Fatal(err)
			}

			if err := chat.SendMessage("", ""); err != nil {
				log.Fatal(err)
			}
		},
	}

	root.AddCommand(upload, chat)
	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}
