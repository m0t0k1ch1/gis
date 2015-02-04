gis
===

command line GitHub issue list viewer

ref. [#8 Golang 製 コマンドラインから GitHub の issue / p-r をよしなに閲覧できる君](http://tech.kayac.com/archive/8_golang_github_issue.html)

[![Gyazo](http://i.gyazo.com/2c45b0a0a24b95aaa8d005c4d14d190f.gif)](http://gyazo.com/2c45b0a0a24b95aaa8d005c4d14d190f)

## Installation

### go get

``` sh
$ go get github.com/m0t0k1ch1/gis
```

### Create personal access token

Create [personal access token](https://help.github.com/articles/creating-an-access-token-for-command-line-use).

[![Gyazo](http://i.gyazo.com/f66e8571a00e7a8590c04db62c4df744.png)](http://gyazo.com/f66e8571a00e7a8590c04db62c4df744)

### Set gis.token

``` sh
$ git config --global gis.token [your personal access token]
```
