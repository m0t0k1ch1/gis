package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"golang.org/x/oauth2"

	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
)

const (
	DefaultPage  = 1
	DefaultLimit = 50
)

var (
	re = regexp.MustCompile(`^(?:(?:ssh://)?git@github\.com(?::|/)|https://github\.com/)([^/]+)/([^/]+?)(?:\.git)?$`)
)

func main() {
	user, err := getUser()
	if err != nil {
		exit(err)
	}

	app := cli.NewApp()
	app.Name = "gis"
	app.Usage = "show GitHub issue list"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "user, u",
			Value: user,
			Usage: "target user",
		},
		cli.BoolFlag{
			Name:  "assignee, a",
			Usage: "filter issues based on their assignee",
		},
		cli.BoolFlag{
			Name:  "mentioned, m",
			Usage: "filter issues to those mentioned target user",
		},
		cli.IntFlag{
			Name:  "page, p",
			Value: DefaultPage,
			Usage: "page of results to retrieve",
		},
		cli.IntFlag{
			Name:  "limit, l",
			Value: DefaultLimit,
			Usage: "the number of results to include per page",
		},
	}
	app.Action = func(c *cli.Context) {
		token, err := getToken()
		if err != nil {
			exit(err)
		}
		owner, repo, err := getOwnerAndRepo()
		if err != nil {
			exit(err)
		}

		ctx := context.Background()

		ts := oauth2.StaticTokenSource(
			&oauth2.Token{
				AccessToken: token,
			},
		)

		client := github.NewClient(oauth2.NewClient(ctx, ts))

		opt := &github.IssueListByRepoOptions{
			ListOptions: github.ListOptions{
				Page:    c.Int("page"),
				PerPage: c.Int("limit"),
			},
		}
		if c.Bool("assignee") {
			opt.Assignee = c.String("user")
		}
		if c.Bool("mentioned") {
			opt.Mentioned = c.String("user")
		}

		issues, _, err := client.Issues.ListByRepo(ctx, owner, repo, opt)
		if err != nil {
			exit(err)
		}

		for _, issue := range issues {
			fmt.Printf("%d\t%s\t%s\n", *issue.Number, *issue.HTMLURL, *issue.Title)
		}
	}

	app.Run(os.Args)
}

func getUser() (string, error) {
	return getGitConfig("user.name")
}

func getToken() (string, error) {
	return getGitConfig("gis.token")
}

func getOwnerAndRepo() (owner, repo string, err error) {
	url, err := getGitConfig("remote.origin.url")
	if err != nil {
		return
	}

	matches := re.FindStringSubmatch(url)
	if len(matches) != 3 {
		err = fmt.Errorf("can't parse remote.origin.url")
		return
	}

	owner = matches[1]
	repo = matches[2]
	return
}

func getGitConfig(key string) (val string, err error) {
	cmd := exec.Command("git", "config", "--get", key)
	var out bytes.Buffer

	cmd.Stdout = &out
	if err = cmd.Run(); err != nil {
		err = fmt.Errorf("can't get git config: %s", key)
		return
	}

	val = strings.TrimSpace(out.String())
	return
}

func exit(err error) {
	fmt.Println(err)
	os.Exit(1)
}
