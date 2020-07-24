package main

import (
	"fmt"
	"github.com/psmarcin/jira-versioner/pkg/git"
	"github.com/psmarcin/jira-versioner/pkg/jira"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

var (
	rootCmd = &cobra.Command{
		Use:   "jira-versioner",
		Short: "A simple version setter for Jira tasks since last version",
		Long: `A solution for automatically create version, 
link all issues from commits to newly created version. 
All automatically.`,
		Run: rootFunc,
	}
)

func init() {
	// get current directory path
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	pwd := filepath.Dir(ex)

	//rootCmd.Flags().StringP("verbose", "v", "info", "")
	rootCmd.Flags().StringP("jira-version", "v", "", "Version name for Jira")
	rootCmd.Flags().StringP("tag", "t", "", "Existing git tag")
	rootCmd.Flags().StringP("jira-email", "e", "", "Jira email")
	rootCmd.Flags().StringP("jira-token", "k", "", "Jira token/key")
	rootCmd.Flags().StringP("jira-project", "p", "", "Jira project, it has to be ID, example: 10003")
	rootCmd.Flags().StringP("jira-base-url", "u", "", "Jira service base url, example: https://example.atlassian.net")
	rootCmd.Flags().StringP("dir", "d", pwd, "Absolute directory path to git repository")
	rootCmd.MarkFlagRequired("tag")
	rootCmd.MarkFlagRequired("jira-email")
	rootCmd.MarkFlagRequired("jira-token")
	rootCmd.MarkFlagRequired("jira-project")
	rootCmd.MarkFlagRequired("jira-base-url")

	rootCmd.Example = "jira-versioner -e jira@example.com -k SOME_TOKEN -p 10003 -t v1.1.0 -u https://example.atlassian.net"
}

func main() {
	Execute()
}

func rootFunc(c *cobra.Command, args []string) {
	tag := c.Flag("tag").Value.String()

	version := c.Flag("jira-version").Value.String()
	if version == ""{
		version = tag
	}

	jiraEmail := c.Flag("jira-email").Value.String()
	jiraToken := c.Flag("jira-token").Value.String()
	jiraProject := c.Flag("jira-project").Value.String()
	jiraBaseUrl := c.Flag("jira-base-url").Value.String()
	gitDir := c.Flag("dir").Value.String()
	log.Printf("[JIRA-VERSIONER] git directory: %s", gitDir)

	g := git.New(gitDir)

	tasks, err := g.GetTasks(tag)
	if err != nil {
		log.Panicf("[GIT] error while getting tasks since latest commit %+v", err)
	}

	j, err := jira.New(jiraEmail, jiraToken, jiraProject, jiraBaseUrl)
	if err != nil {
		log.Panicf("[VERSION] error while connecting to jira server %+v", err)
	}

	_, err = j.CreateVersion(version)
	if err != nil {
		log.Panicf("[VERSION] error while creating version %+v", err)
	}

	err = j.LinkTasksToVersion(tasks)
	if err != nil {
		log.Panicf("[VERSION] can't update task to fix version %s", err)
	}

	log.Print("[JIRA-VERSIONER] Done ✅")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
