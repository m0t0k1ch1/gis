package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"code.google.com/p/goauth2/oauth"

	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
)

func main() {
	user, err := getUser()
	if err != nil {
		log.Fatal(err)
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
			Value: 1,
			Usage: "page of results to retrieve",
		},
		cli.IntFlag{
			Name:  "limit, l",
			Value: 50,
			Usage: "the number of results to include per page",
		},
	}
	app.Action = func(c *cli.Context) {
		token, err := getToken()
		if err != nil {
			log.Fatal(err)
		}
		owner, repo, err := getOwnerAndRepo()
		if err != nil {
			log.Fatal(err)
		}

		t := &oauth.Transport{
			Token: &oauth.Token{AccessToken: token},
		}
		client := github.NewClient(t.Client())

		opt := &github.IssueListByRepoOptions{
			ListOptions: github.ListOptions{
				Page:    c.Int("page"),
				PerPage: c.Int("limit"),
			},
		}
		if c.Bool("assignee") {
			opt.Assignee = user
		}
		if c.Bool("mentioned") {
			opt.Mentioned = user
		}

		issues, _, err := client.Issues.ListByRepo(owner, repo, opt)
		if err != nil {
			log.Fatal(err)
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

func getOwnerAndRepo() (string, string, error) {
	url, err := getGitConfig("remote.origin.url")
	if err != nil {
		return "", "", err
	}

	re, err := regexp.Compile(`git@github\.com:(.+)/(.+)\.git`)
	if err != nil {
		return "", "", err
	}

	matches := re.FindStringSubmatch(url)
	if len(matches) != 3 {
		return "", "", fmt.Errorf("can't parse remote.origin.url")
	}

	return matches[1], matches[2], nil
}

func getGitConfig(key string) (string, error) {
	cmd := exec.Command("git", "config", "--get", key)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}
