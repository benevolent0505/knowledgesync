package main

import (
	"github.com/benevolent0505/knowledgesync/command"
	"github.com/mitchellh/cli"
)

func Commands(meta *command.Meta, conf *command.Config) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"pull": func() (cli.Command, error) {
			return &command.PullCommand{
				Config: *conf,
				Meta: *meta,
			}, nil
		},
		"push": func() (cli.Command, error) {
			return &command.PushCommand{
				Meta: *meta,
			}, nil
		},
		"post": func() (cli.Command, error) {
			return &command.PostCommand{
				Meta: *meta,
			}, nil
		},

		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				Meta:     *meta,
				Version:  Version,
				Revision: GitCommit,
				Name:     Name,
			}, nil
		},
	}
}
