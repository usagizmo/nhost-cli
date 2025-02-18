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
	"os"
	"os/exec"
	"runtime"

	"github.com/manifoldco/promptui"
	"github.com/nhost/cli/nhost"
	"github.com/nhost/cli/util"
	"github.com/spf13/cobra"
)

var (

	//  initialize flags for every choice
	issue   bool
	chat    bool
	discuss bool
	//  wiki    bool

	noBrowser bool
)

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "support",
	Short: "Reach out to us",
	Long: `Launches URL in browser to allow
you to open issues and submit bug reports
in case you encounter something broken with this CLI.

Or even chat with our team and start a new discussion.`,
	Run: func(cmd *cobra.Command, args []string) {

		options := []map[string]interface{}{
			{
				"text":  "Report bugs & open feature requests",
				"value": fmt.Sprintf("https://github.com/%v/issues", nhost.REPOSITORY),
				"flag":  issue,
			},
			{
				"text":  "Chat with our team",
				"value": "https://discord.com/invite/9V7Qb2U",
				"flag":  chat,
			},
			{
				"text":  "Start a new discussion",
				"value": fmt.Sprintf("https://github.com/%v/discussions/new", nhost.REPOSITORY),
				"flag":  discuss,
			},
			/*
				{
					"text":  "Advanced usage, deployment help, example apps & more",
					"value": fmt.Sprintf("https://github.com/%v/wiki", nhost.REPOSITORY),
					"flag":  wiki,
				},
			*/
		}

		//  if the user has passed the flag for any option,
		//  launch those directly
		//  bypassing selection prompt

		ok := false
		for _, item := range options {

			if item["flag"].(bool) {

				ok = true
				if noBrowser {
					status.Info(fmt.Sprintf(item["text"].(string), " @ ", util.Bold, item["value"], util.Reset))
				} else if err := openbrowser(item["value"].(string)); err != nil {
					log.Debug(err)
					status.Errorln("Failed to launch browser")
					status.Info(fmt.Sprintf(item["text"].(string), " @ ", item["value"]))
				}
			}
		}

		if ok {
			os.Exit(0)
		}

		//  configure interactive prompt template
		templates := promptui.SelectTemplates{
			Active:   `{{ "✔" | green | bold }} {{ .text | cyan | bold }} {{ .value | faint }}`,
			Inactive: `   {{ .text | cyan | bold }} `,
			Selected: `{{ "✔" | green | bold }} {{ "Selected" | bold }}: {{ .text | cyan | bold }}`,
		}

		//  configure interative prompt
		prompt := promptui.Select{
			Label:     "Select Option",
			Items:     options,
			Templates: &templates,
		}

		index, _, err := prompt.Run()
		if err != nil {
			os.Exit(0)
		}

		selected := options[index]

		if noBrowser {
			status.Info(selected["text"].(string) + " @ " + util.Bold + selected["value"].(string) + util.Reset)
		} else {
			//  launch browser
			if err := openbrowser(selected["value"].(string)); err != nil {
				log.Debug(err)
				status.Errorln("Failed to launch browser")
				status.Info(selected["text"].(string) + " @ " + selected["value"].(string))
			}
		}
	},
}

func linuxOpener() (string, error) {
	cmd := "xdg-open"
	if _, err := exec.LookPath(cmd); err != nil {
		// try "wslview", maybe we're on WSL
		if _, err := exec.LookPath("wslview"); err == nil {
			cmd = "wslview"
		} else {
			return "", err
		}
	}

	return cmd, nil
}

func openbrowser(url string) error {
	switch runtime.GOOS {
	case "linux":
		openCmd, err := linuxOpener()
		if err != nil {
			return err
		}

		return exec.Command(openCmd, url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return fmt.Errorf("unsupported platform")
	}
}

func init() {
	rootCmd.AddCommand(reportCmd)

	//  Here you will define your flags and configuration settings.

	//  Cobra supports Persistent Flags which will work for this command
	//  and all subcommands, e.g.:
	//  reportCmd.PersistentFlags().String("foo", "", "A help for foo")

	//  Cobra supports local flags which will only run when this command
	//  is called directly, e.g.:
	reportCmd.Flags().BoolVar(&noBrowser, "no-browser", false, "Don't open in browser")
	reportCmd.Flags().BoolVar(&issue, "issue", false, "Open Issue on GitHub")
	reportCmd.Flags().BoolVar(&chat, "chat", false, "Launch Nhost Discord Server")
	reportCmd.Flags().BoolVar(&discuss, "discuss", false, "Launch GitHub Discussions")
	//  reportCmd.Flags().BoolVar(&wiki, "wiki", false, "Launch wiki in browser")
}
