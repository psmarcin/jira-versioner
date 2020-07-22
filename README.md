# Jira-versioner
Base on commit messages and tags you can now create versions in Jira. 

## Getting Started

These instructions will show you how to use jira-versioner in simple steps.

### Prerequisites

Things that you need to have before we start:

* Jira service
* Jira project
* Jira email and token with rights to write to Jira project
* Git repository
* At least two Git tags in rage

### Installing

Just grab one file from download page: https://github.com/psmarcin/jira-versioner/releases/latest.

## Examples 

```shell script
jira-versioner -e jira@example.com -k SOME_TOKEN -p 10003 -v v1.1.0 -t v1.1.0 -u https://example.atlassian.net
```

Help: 
```shell script
Usage:
  jira-versioner [flags]

Examples:
jira-versioner -e jira@example.com -k SOME_TOKEN -p 10003 -v v1.1.0 -t v1.1.0 -u https://example.atlassian.net

Flags:
  -d, --dir string             Absolute directory path to git repository (default "/Users/psmarcin/projects/jira-releaser")
  -h, --help                   help for jira-versioner
  -u, --jira-base-url string   Jira service base url, example: https://example.atlassian.net
  -e, --jira-email string      Jira email
  -p, --jira-project string    Jira project, it has to be ID, example: 10003
  -k, --jira-token string      Jira token/key
  -v, --jira-version string    Version name for Jira
  -t, --tag string             Existing git tag

required flag(s) "jira-base-url", "jira-email", "jira-project", "jira-token", "jira-version", "tag"
```

## Contributing

TODO:

Please read [CONTRIBUTING.md](https://gist.github.com/PurpleBooth/b24679402957c63ec426) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/psmarcin/jira-versioner/tags). 

## Authors

* **psmarcin** - [psmarcin](https://github.com/psmarcin)

See also the list of [contributors](https://github.com/psmarcin/jira-versioner/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details


