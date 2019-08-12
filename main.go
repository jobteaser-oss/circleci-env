/*
 * Copyright 2019 Jobteaser <opensource@jobteaser.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"github.com/galdor/go-cmdline"
	"github.com/jobteaser/circleci-env/circleci"
	"os"
)

var token string

func main() {
	cli := cmdline.New()
	cli.AddOption("t", "token", "", "circleci token")
	cli.AddOption("v", "vcs-type", "github", "the vcs type of the project")
	cli.AddOption("u", "username", "", "the username who host the project")
	cli.AddOption("p", "project", "", "the circleci project name")
	cli.AddCommand("list", "list env")
	cli.AddCommand("get", "get env")
	cli.AddCommand("set", "set env")
	cli.AddCommand("del", "delete env")
	cli.Parse(os.Args)

	if !cli.IsOptionSet("t") {
		fmt.Println("error: the circle token is required")
		os.Exit(1)
	}

	client, err := circleci.NewClient(cli.OptionValue("t"))
	if err != nil {
		fmt.Printf("error: cannot create circleci client: %v\n", err)
		os.Exit(1)
	}

	var cmdFun func(*circleci.Client, string, string, string, []string)
	switch cli.CommandName() {
	case "get":
		cmdFun = cmdGet
	case "list":
		cmdFun = cmdList
	case "set":
		cmdFun = cmdSet
	case "del":
		cmdFun = cmdDel
	}

	cmdFun(
		client,
		cli.OptionValue("vcs-type"),
		cli.OptionValue("username"),
		cli.OptionValue("project"),
		cli.CommandNameAndArguments(),
	)
}

func cmdGet(client *circleci.Client, vcsType, username, project string, args []string) {
	cli := cmdline.New()
	cli.AddArgument("key", "the name of the environment variable")
	cli.Parse(args)

	env, err := client.GetEnv(
		vcsType,
		username,
		project,
		cli.ArgumentValue("key"),
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s=%s\n", env.Key, env.Value)
}

func cmdList(client *circleci.Client, vcsType, username, project string, args []string) {
	cli := cmdline.New()
	cli.Parse(args)

	envs, err := client.ListEnv(vcsType, username, project)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	for _, env := range envs {
		fmt.Printf("%s=%s\n", env.Key, env.Value)
	}
}

func cmdSet(client *circleci.Client, vcsType, username, project string, args []string) {
	cli := cmdline.New()
	cli.AddArgument("key", "the name of the environment variable")
	cli.AddArgument("value", "the value of the environment variable")
	cli.Parse(args)

	err := client.SetEnv(
		vcsType,
		username,
		project,
		cli.ArgumentValue("key"),
		cli.ArgumentValue("value"),
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func cmdDel(client *circleci.Client, vcsType, username, project string, args []string) {
	cli := cmdline.New()
	cli.AddArgument("key", "the name of the environment variable")
	cli.Parse(args)

	err := client.DeleteEnv(
		vcsType,
		username,
		project,
		cli.ArgumentValue("key"),
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
