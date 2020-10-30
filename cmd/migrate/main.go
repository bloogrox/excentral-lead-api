package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gitlab.com/cpanova/excentral/domain/conversion"
	"gitlab.com/cpanova/excentral/domain/lead"
	"gitlab.com/cpanova/excentral/domain/partner"
	"gitlab.com/cpanova/excentral/domain/postback"
	"gitlab.com/cpanova/excentral/domain/sender"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		databaseURL := os.Getenv("DATABASE_URL")
		db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
		if err != nil {
			panic(err)
		}

		log.Println("Run Migrate")

		db.AutoMigrate(
			&lead.Lead{},
			&conversion.Conversion{},
			&postback.Postback{},
			&partner.Partner{},
			&sender.Sender{},
		)

		log.Println("Done")
	},
}

func main() {
	migrateCmd.Execute()
}
