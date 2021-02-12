package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/psmarcin/jira-versioner/pkg/git"
	"github.com/psmarcin/jira-versioner/pkg/jira"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "jira-versioner",
		Short: "A simple version setter for Jira tasks since last version",
		Long: `A solution for automatically create version, 
link all issues from commits to newly created version. 
All automatically.`,
		Run: rootFunc,
	}
	// get current directory path
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	pwd := filepath.Dir(ex)

	rootCmd.Flags().StringP("jira-version", "v", "", "Version name for Jira")
	rootCmd.Flags().StringP("tag", "t", "", "Existing git tag")
	rootCmd.Flags().StringP("jira-email", "e", "", "Jira email")
	rootCmd.Flags().StringP("jira-token", "k", "", "Jira token/key/password")
	rootCmd.Flags().StringP("jira-project", "p", "", "Jira project, it has to be ID, example: 10003")
	rootCmd.Flags().StringP("jira-base-url", "u", "", "Jira service base url, example: https://example.atlassian.net")
	rootCmd.Flags().StringP("dir", "d", pwd, "Absolute directory path to git repository")
	_ = rootCmd.Flags().BoolP("dry-run", "", false, "Enable dry run mode")

	err = rootCmd.MarkFlagRequired("tag")
	if err != nil {
		fmt.Printf("err: %+v", err)
		os.Exit(1)
	}
	err = rootCmd.MarkFlagRequired("jira-email")
	if err != nil {
		fmt.Printf("err: %+v", err)
		os.Exit(1)
	}
	err = rootCmd.MarkFlagRequired("jira-token")
	if err != nil {
		fmt.Printf("err: %+v", err)
		os.Exit(1)
	}
	err = rootCmd.MarkFlagRequired("jira-project")
	if err != nil {
		fmt.Printf("err: %+v", err)
		os.Exit(1)
	}
	err = rootCmd.MarkFlagRequired("jira-base-url")
	if err != nil {
		fmt.Printf("err: %+v", err)
		os.Exit(1)
	}

	rootCmd.Example = "jira-versioner -e jira@example.com -k pa$$wor0 -p 10003 -t v1.1.0 -u https://example.atlassian.net"

	if err = rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func rootFunc(c *cobra.Command, _ []string) {
	log := zap.NewExample().Sugar()
	dryRun := false
	defer func() {
		_ = log.Sync()
	}()

	tag := c.Flag("tag").Value.String()

	version := c.Flag("jira-version").Value.String()
	if version == "" {
		version = tag
	}

	jiraEmail := c.Flag("jira-email").Value.String()
	jiraToken := c.Flag("jira-token").Value.String()
	jiraProject := c.Flag("jira-project").Value.String()
	jiraBaseURL := c.Flag("jira-base-url").Value.String()
	dryRunRaw := c.Flag("dry-run").Value.String()
	if dryRunRaw == "true" {
		dryRun = true
	}
	gitDir := c.Flag("dir").Value.String()

	log.Debugf(
		"[JIRA-VERSIONER] starting with parameters: %+v",
		map[string]interface{}{
			"jiraEmail":   jiraEmail,
			"jiraToken":   jiraToken,
			"jiraProject": jiraProject,
			"jiraBaseURL": jiraBaseURL,
			"gitDir":      gitDir,
			"tag":         tag,
			"version":     version,
			"dryRun":      dryRun,
		},
	)
	log.Infof("[JIRA-VERSIONER] git directory: %s", gitDir)

	g := git.New(gitDir, log)

	tasks, err := g.GetTasks(tag)
	if err != nil {
		log.Fatalf("[GIT] error while getting tasks since latest commit %+v", err) // nolint
	}

	var jiraConfig = jira.Config{
		Username:  jiraEmail,
		Token:     jiraToken,
		ProjectID: jiraProject,
		BaseURL:   jiraBaseURL,
		Log:       log,
		DryRun:    dryRun,
	}
	j, err := jira.New(&jiraConfig)
	if err != nil {
		log.Fatalf("[VERSION] error while connecting to jira server %+v", err)
	}

	_, err = j.CreateVersion(version)
	if err != nil {
		log.Fatalf("[VERSION] error while creating version %+v", err)
	}

	j.LinkTasksToVersion(tasks)

	log.Infof("[JIRA-VERSIONER] done ✅")
}
