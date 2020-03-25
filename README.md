# spotify-slack-status

[![Go Report Card](https://goreportcard.com/badge/github.com/adamveld12/spotify-slack-status)](https://goreportcard.com/report/github.com/adamveld12/spotify-slack-status)
[![Gocover](https://gocover.io/_badge/github.com/adamveld12/spotify-slack-status)](https://gocover.io/github.com/adamveld12/spotify-slack-status)
[![Go Walker](http://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/adamveld12/spotify-slack-status)
[![Build Status](https://semaphoreci.com/api/v1/adamveld12/spotify-slack-status/branches/master/badge.svg)](https://semaphoreci.com/adamveld12/spotify-slack-status)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

Uses DBus to read the currently playing song on Spotify and updates your status in Slack accordingly.

## How to use

- Goes without saying but this required DBus

- Install `go get -u github.com/adamveld12/spotify-slack-status`
  - make sure your `$GOPATH/bin` is in your path

- Run it: `spotify-slack-status -token $SLACK_TOKEN -playing-emoji "headphones"`

## LICENSE

GPL-V3
