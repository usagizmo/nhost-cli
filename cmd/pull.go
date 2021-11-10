/*
MIT License

Copyright (c) Nhost

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package cmd

import (
	"fmt"
	"os/exec"

	"github.com/nhost/cli/hasura"
	"github.com/nhost/cli/nhost"
	"github.com/spf13/cobra"
)

//  pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:     "pull",
	Aliases: []string{"p"},
	Short:   "Pull migrations from remote",
	Hidden:  true,
	Long: `Pull latest migrations and metadata changes from remote
and sync them with your local app.`,
	/*
		Run: func(cmd *cobra.Command, args []string) {

			//  warn the user of upcoming dangers
			log.Warn("This can potentially break your local changes")
			status.Info("Backup your local changes before proceeding ahead")

			//  configure interative prompt
			confirmationPrompt := promptui.Prompt{
				Label:     "Are you sure you want to continue",
				IsConfirm: true,
			}

			_, err := confirmationPrompt.Run()
			if err != nil {
				log.Debug(err)
				os.Exit(0)
			}

			//  validate authentication
			user, err := validateAuth(nhost.AUTH_PATH)
			if err != nil {
				log.Debug(err)
				status.Errorln("Failed to validate authentication")

				//  begin the login procedure
				loginCmd.Run(cmd, args)
			}

			var projects []nhost.Project
			projects = append(projects, user.Projects...)
			for _, item := range user.Teams {
				projects = append(projects, item.Team.Projects...)
			}

			info, err := nhost.Info()
			if err != nil {
				log.Debug(err)
				log.Fatal("Failed to read saved Nhost app information")
			}

			var linkedProject nhost.Project

			for _, project := range projects {
				if project.ID == info.ProjectID {
					linkedProject = project
				}
			}

			//  intialize common options
			hasuraEndpoint := "https://" + linkedProject.ProjectDomains.Hasura
			adminSecret := linkedProject.HasuraGQEAdminSecret

			commonOptions := []string{"--endpoint", hasuraEndpoint, "--admin-secret", adminSecret, "--skip-update-check"}

			//  create migration
			//  and notify remote to skip it

			//  test new hasura client
			hasuraClient := hasura.Client{}
			hasuraClient.Init(hasuraEndpoint, adminSecret, nil)

			_, err = pullMigration(hasuraClient, "pulled_from_remote", commonOptions)
			if err != nil {
				log.Debug(err)
				log.Fatal("Failed to create migration from remote")
			}


			status.Info("Migrations pulled from remote")
			status.Info("To apply the new migrations use `nhost dev`")
		},
	*/
}

func pullMigration(client hasura.Client, name string) (hasura.Migration, error) {

	log.Debugf("Creating migration '%s'", name)

	migration := hasura.Migration{
		Name: name,
	}

	migration = migration.Init()
	/*

		metadata, err := client.GetMetadata()
		if err != nil {
			return migration, err
		}

		//  Fetch list of all ALLOWED schemas before applying
		schemas, err := client.GetSchemas()
		if err != nil {
			log.Debug("Failed to get list of schemas")
			return migration, err
		}

		migrationTables := getMigrationTables(schemas, metadata.Tables)
	*/

	//  fetch migrations
	log.Debug("Creating initial migration")
	args := []string{client.CLI, "migrate", "create", name, "--from-server"}
	//	migrationArgs = append(migrationArgs, getMigrationTables(schemas, tables)...)
	args = append(args, client.CommonOptions...)

	execute := exec.Cmd{
		Path: client.CLI,
		Args: args,
		Dir:  nhost.NHOST_DIR,
	}

	output, err := execute.CombinedOutput()
	if err != nil {
		log.Debug(string(output))
		return migration, err
	}

	log.Debug("Clearing remote migrations")

	args = []string{client.CLI, "migrate", "delete", "--all", "--server", "--force"}
	args = append(args, client.CommonOptions...)

	execute = exec.Cmd{
		Path: client.CLI,
		Args: args,
		Dir:  nhost.NHOST_DIR,
	}

	output, err = execute.CombinedOutput()
	if err != nil {
		log.Debug(string(output))
		return migration, err
	}

	log.Debug("Applying migrations")

	args = []string{client.CLI, "migrate", "apply", "--skip-execution"}
	args = append(args, client.CommonOptions...)

	execute = exec.Cmd{
		Path: client.CLI,
		Args: args,
		Dir:  nhost.NHOST_DIR,
	}

	output, err = execute.CombinedOutput()
	if err != nil {
		log.Debug(string(output))
		return migration, err
	}

	log.Debug("Export metadata")

	args = []string{client.CLI, "metadata", "export"}
	args = append(args, client.CommonOptionsWithoutDB...)

	execute = exec.Cmd{
		Path: client.CLI,
		Args: args,
		Dir:  nhost.NHOST_DIR,
	}

	output, err = execute.CombinedOutput()
	if err != nil {
		log.Debug(string(output))
		return migration, err
	}

	return migration, nil
}

/*
func copySeedToMigration(seed, migration string) error {

	//  first, search for the newly created seed file
	seedFile, err := searchFile(seed, nhost.SEEDS_DIR)
	if err != nil {
		return err
	}

	//  now read it's data
	data, err := os.ReadFile(filepath.Join(nhost.SEEDS_DIR, seedFile.Name()))
	if err != nil {
		return err
	}

	//  finally append the seed data to migration file
	if err = writeToFile(filepath.Join(nhost.MIGRATIONS_DIR, migration), string(data), "end"); err != nil {
		return err
	}

	//  delete the seed file so that it's not applied again
	if err = deleteAllPaths(filepath.Join(nhost.SEEDS_DIR, seedFile.Name())); err != nil {
		return err
	}

	return nil
}

func searchFile(name, directory string) (fs.DirEntry, error) {

	migrations, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for _, item := range migrations {
		if strings.Contains(item.Name(), name) {
			return item, nil
		}
	}

	return nil, errors.New("failed to find file %v in %v")
}


func getFormattedSchemas(list []string) []string {

	var response []string

	for _, item := range list {
		response = append(response, "--schema")
		response = append(response, item)
	}
	return response
}
*/

func getMigrationTables(schemas []string, tables []hasura.TableEntry) []string {

	var response []string

	for _, table := range tables {
		for _, schema := range schemas {
			if table.Table.Schema == schema {
				response = append(response, "--table")
				response = append(response, fmt.Sprintf(
					`%s.%s`,
					schema,
					table.Table.Name,
				))
			}
		}
	}

	/*
		for _, value := range filteredValues {
			if value != "public.users" {
				fromTables = append(fromTables, "--table")
				fromTables = append(fromTables, value)
			}
		}
	*/
	return response
}

func init() {
	//  rootCmd.AddCommand(pullCmd)

	//  Here you will define your flags and configuration settings.

	//  Cobra supports Persistent Flags which will work for this command
	//  and all subcommands, e.g.:
	//  pullCmd.PersistentFlags().String("foo", "", "A help for foo")

	//  Cobra supports local flags which will only run when this command
	//  is called directly, e.g.:
	//  pullCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
