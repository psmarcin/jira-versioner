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

### How does it work

Here is our git log history:

```shell script
05e5705322cc2d9daf7fb376a8c5e9cbd039b257 (HEAD -> master, tag: v2.1.0, origin/master) chore: remove unnecessary string conversion
9bf13576317845cd7d10980d62afe719872ceb01 feat: error logs contains command output JR-4
831e4c253829dbc12683baa5b4d494aa3524f39f feat: jira version not required, default tag JR-13
533569497a68f04674e43e23d06fb9c1f0b3b958 docs: update readme.md with new name JR-2
1e1dd3131aeed3611e70d5f329989c1a09371822 (tag: v2.0.1) chore: rename to jira-versioner JR-3
aeae65755553d03920d7cd7c4a5fdb40a02d7c57 docs: update command name JR-3
e874a9c6162fd102b9de926397a855c1b0dbd880 docs: README.md file JR-2
d15916037b0a6ca04776e474ac461e767631c838 (tag: v2.0.0) feat: consistent arguments name
cb59ea7f0bc3efb8b92de87cd88b589024d18ee7 (tag: v1.1.0) feat: JR-40 argument to run git commands in different path
2e6d61dee0c4ed3a0f7f887973dbc326a487675b (tag: v1.0.0) feat: github JR-1 release action
```

For simplification in examples I omitted jira configs 

1. `jira-versioner -t v2.1.0`
    
    Found commits:
    1. `05e5705322cc2d9daf7fb376a8c5e9cbd039b257`
    1. `9bf13576317845cd7d10980d62afe719872ceb01`
    1. `831e4c253829dbc12683baa5b4d494aa3524f39f`
    1. `533569497a68f04674e43e23d06fb9c1f0b3b958`
    
    Found tasks:
    1. `JR-4`
    1. `JR-13`
    1. `JR-2`
    
    Results: 
    1. New version created(if not already exists) - `v2.1.0`
    1. If task was found in commits and exists in Jira it will set fixed version for it
    
1. `jira-versioner -t v2.0.1`
    
    Found commits:
    1. `1e1dd3131aeed3611e70d5f329989c1a09371822`
    1. `aeae65755553d03920d7cd7c4a5fdb40a02d7c57`
    1. `e874a9c6162fd102b9de926397a855c1b0dbd880`
    
    Found tasks:
    1. `JR-3`
    1. `JR-2`
    
    Results: 
    1. New version created(if not already exists) - `v2.0.1`
    1. If task was found in commits and exists in Jira it will set fixed version for it
    1. Task from commits are always unique (no deuplicates)
    
1. `jira-versioner -t v2.0.0`
    
    Found commits:
    1. `d15916037b0a6ca04776e474ac461e767631c838`
    
    Found tasks:
    none
    
    Results: 
    1. New version created(if not already exists) - `v2.0.0`
    1. No tasks found, no links = empty version
    
1. `jira-versioner -t v1.1.0 -v 1000.000.000`
    
    Found commits:
    1. `cb59ea7f0bc3efb8b92de87cd88b589024d18ee7`
    
    Found tasks:
    1. `JR-40`
    
    Results: 
    1. New version created(if not already exists) - `1000.000.000`
    1. If task was found in commits and exists in Jira it will set fixed version for it
    1. It looks for task id in whole commit message


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


