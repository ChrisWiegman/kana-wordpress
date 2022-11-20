package database

import (
	"fmt"
	"os"
	"path"

	"github.com/ChrisWiegman/kana-cli/internal/site"
)

func Import(site *site.Site, file string, preserve bool, replaceDomain string) error {

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	rawImportFile := path.Join(cwd, file)
	if _, err = os.Stat(rawImportFile); os.IsNotExist(err) {
		return fmt.Errorf("the specified sql file does not exist. Please enter a valid file to import")
	}

	kanaImportFile := path.Join(site.StaticConfig.SiteDirectory, "import.sql")

	err = copyFile(rawImportFile, kanaImportFile)
	if err != nil {
		return err
	}

	if !preserve {

		fmt.Println("Dropping the existing database.")

		dropCommand := []string{
			"db",
			"drop",
			"--yes",
		}

		createCommand := []string{
			"db",
			"create",
		}

		_, err := site.RunWPCli(dropCommand)
		if err != nil {
			return err
		}

		_, err = site.RunWPCli(createCommand)
		if err != nil {
			return err
		}
	}

	fmt.Println("Importing the database file.")

	importCommand := []string{
		"db",
		"import",
		"/Site/import.sql",
	}

	_, err = site.RunWPCli(importCommand)
	if err != nil {
		return err
	}

	if replaceDomain != "" {

		fmt.Println("Replacing the old domain name")

		replaceCommand := []string{
			"search-replace",
			replaceDomain,
			fmt.Sprintf("%s.%s", site.StaticConfig.SiteName, site.StaticConfig.AppDomain),
			"--all-tables",
		}

		_, err := site.RunWPCli(replaceCommand)
		if err != nil {
			return err
		}
	}

	return nil
}
