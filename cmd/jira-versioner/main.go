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
		Use:   "jira-releaser",
		Short: "A simple version setter for Jira tasks since last version",
		Long: `A solution for automatically create release, 
link all issues from commits to newly created release. 
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
	log.Printf("[JIRA-VERSIONER] git directory: %s", pwd)

	//rootCmd.Flags().StringP("verbose", "v", "info", "")
	rootCmd.Flags().StringP("version", "", "", "release name, required, must be unique")
	rootCmd.Flags().StringP("tag", "t", "", "release name, required, must be unique")
	rootCmd.Flags().StringP("jira-email", "e", "", "Jira email")
	rootCmd.Flags().StringP("jira-token", "k", "", "Jira token/key")
	rootCmd.Flags().StringP("jira-project", "p", "", "Jira project")
	rootCmd.Flags().StringP("jira-base-url", "u", "", "Jira service base url")
	rootCmd.Flags().StringP("dir", "d", pwd, "absolute directory path to git repository")
	rootCmd.MarkFlagRequired("version")
	rootCmd.MarkFlagRequired("tag")
	rootCmd.MarkFlagRequired("jira-email")
	rootCmd.MarkFlagRequired("jira-token")
	rootCmd.MarkFlagRequired("jira-project")
	rootCmd.MarkFlagRequired("jira-base-url")
}

func main() {
	Execute()
}

func rootFunc(c *cobra.Command, args []string) {
	version := c.Flag("version").Value.String()
	tag := c.Flag("tag").Value.String()
	jiraEmail := c.Flag("jira-email").Value.String()
	jiraToken := c.Flag("jira-token").Value.String()
	jiraProject := c.Flag("jira-project").Value.String()
	jiraBaseUrl := c.Flag("jira-base-url").Value.String()
	gitDir := c.Flag("dir").Value.String()

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
