package db

import (
	"fmt"
	"os"

	"github.com/catalyzeio/cli/models"
	"github.com/catalyzeio/cli/prompts"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "db",
	ShortHelp: "Tasks for databases",
	LongHelp:  "Tasks for databases",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Command(BackupSubCmd.Name, BackupSubCmd.ShortHelp, BackupSubCmd.CmdFunc(settings))
			cmd.Command(DownloadSubCmd.Name, DownloadSubCmd.ShortHelp, DownloadSubCmd.CmdFunc(settings))
			cmd.Command(ExportSubCmd.Name, ExportSubCmd.ShortHelp, ExportSubCmd.CmdFunc(settings))
			cmd.Command(ImportSubCmd.Name, ImportSubCmd.ShortHelp, ImportSubCmd.CmdFunc(settings))
			cmd.Command(ListSubCmd.Name, ListSubCmd.ShortHelp, ListSubCmd.CmdFunc(settings))
		}
	},
}

var BackupSubCmd = models.Command{
	Name:      "backup",
	ShortHelp: "Create a new backup",
	LongHelp:  "Create a new backup",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			databaseName := subCmd.StringArg("DATABASE_NAME", "", "The name of the database service to create a backup for (i.e. 'db01')")
			skipPoll := subCmd.BoolOpt("s skip-poll", false, "Whether or not to wait for the backup to finish")
			subCmd.Action = func() {
				id := New(settings, prompts.New(), *databaseName, "", "", "", "", *skipPoll, false, 0, 0)
				err := id.Backup()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}
			subCmd.Spec = "DATABASE_NAME [-s]"
		}
	},
}

var DownloadSubCmd = models.Command{
	Name:      "download",
	ShortHelp: "Download a previously created backup",
	LongHelp:  "Download a previously created backup",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			databaseName := subCmd.StringArg("DATABASE_NAME", "", "The name of the database service which was backed up (i.e. 'db01')")
			backupID := subCmd.StringArg("BACKUP_ID", "", "The ID of the backup to download (found from \"catalyze backup list\")")
			filePath := subCmd.StringArg("FILEPATH", "", "The location to save the downloaded backup to. This location must NOT already exist unless -f is specified")
			force := subCmd.BoolOpt("f force", false, "If a file previously exists at \"filepath\", overwrite it and download the backup")
			subCmd.Action = func() {
				id := New(settings, prompts.New(), *databaseName, *backupID, *filePath, "", "", false, *force, 0, 0)
				err := id.Download()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}
			subCmd.Spec = "DATABASE_NAME BACKUP_ID FILEPATH [-f]"
		}
	},
}

var ExportSubCmd = models.Command{
	Name:      "export",
	ShortHelp: "Export data from a database",
	LongHelp:  "Export data from a database",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			databaseName := subCmd.StringArg("DATABASE_NAME", "", "The name of the database to export data from (i.e. 'db01')")
			filePath := subCmd.StringArg("FILEPATH", "", "The location to save the exported data. This location must NOT already exist unless -f is specified")
			force := subCmd.BoolOpt("f force", false, "If a file previously exists at `filepath`, overwrite it and export data")
			subCmd.Action = func() {
				id := New(settings, prompts.New(), *databaseName, "", *filePath, "", "", false, *force, 0, 0)
				err := id.Export()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}
			subCmd.Spec = "DATABASE_NAME FILEPATH [-f]"
		}
	},
}

var ImportSubCmd = models.Command{
	Name:      "import",
	ShortHelp: "Import data into a database",
	LongHelp:  "Import data into a database",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			databaseName := subCmd.StringArg("DATABASE_NAME", "", "The name of the database to import data to (i.e. 'db01')")
			filePath := subCmd.StringArg("FILEPATH", "", "The location of the file to import to the database")
			mongoCollection := subCmd.StringOpt("c mongo-collection", "", "If importing into a mongo service, the name of the collection to import into")
			mongoDatabase := subCmd.StringOpt("d mongo-database", "", "If importing into a mongo service, the name of the database to import into")
			subCmd.Action = func() {
				id := New(settings, prompts.New(), *databaseName, "", *filePath, *mongoCollection, *mongoDatabase, false, false, 0, 0)
				err := id.Import()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}
			subCmd.Spec = "DATABASE_NAME FILEPATH [-d [-c]]"
		}
	},
}

var ListSubCmd = models.Command{
	Name:      "list",
	ShortHelp: "List created backups",
	LongHelp:  "List created backups",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			databaseName := subCmd.StringArg("DATABASE_NAME", "", "The name of the database service to list backups for (i.e. 'db01')")
			page := subCmd.IntOpt("p page", 1, "The page to view")
			pageSize := subCmd.IntOpt("n page-size", 10, "The number of items to show per page")
			subCmd.Action = func() {
				id := New(settings, prompts.New(), *databaseName, "", "", "", "", false, false, *page, *pageSize)
				err := id.List()
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}
			subCmd.Spec = "DATABASE_NAME [-p] [-n]"
		}
	},
}

// IDb
type IDb interface {
	Backup() error
	Download() error
	Export() error
	Import() error
	List() error
}

// SDb is a concrete implementation of IDb
type SDb struct {
	Settings *models.Settings
	Prompts  prompts.IPrompts

	DatabaseName    string
	BackupID        string
	FilePath        string
	MongoCollection string
	MongoDatabase   string
	SkipPoll        bool
	Force           bool
	Page            int
	PageSize        int
}

// New returns an instance of IDb
func New(settings *models.Settings, prompts prompts.IPrompts, databaseName, backupID, filePath, mongoCollection, mongoDatabase string, skipPoll, force bool, page, pageSize int) IDb {
	return &SDb{
		Settings:        settings,
		Prompts:         prompts,
		DatabaseName:    databaseName,
		BackupID:        backupID,
		FilePath:        filePath,
		MongoCollection: mongoCollection,
		MongoDatabase:   mongoDatabase,
		SkipPoll:        skipPoll,
		Force:           force,
		Page:            page,
		PageSize:        pageSize,
	}
}
