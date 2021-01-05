# gis

a command line GitHub issue list viewer written in Go

ref. [#8 Golang 製 コマンドラインから GitHub の issue / p-r をよしなに閲覧できる君](http://tech.kayac.com/archive/8_golang_github_issue.html)

![gis](https://github.com/m0t0k1ch1/gis/blob/master/img/demo.gif)

## Installation

### go get

``` sh
$ go get -u github.com/m0t0k1ch1/gis
```

### Create personal access token

Create [personal access token](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line).

In the token creation process, select __repo__ scope.

### Set gis.token

``` sh
$ git config --global gis.token [your personal access token]
```
