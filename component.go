// Copyright 2016 Thiago Cangussu de Castro Gomes. All rights reserved.
// Use of this source code is governed by a GNU General Public License
// version 3 that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
)

// ComponentRef is the configuration block that references a component.
type ComponentRef struct {
	Name       string      `json:"name"`
	Version    string      `json:"version"`
	Repo       string      `json:"repo"`
	Repoconfig *RepoConfig `json:"repoconfig"`
}

// Dependency is the configuration block that defines a dependency.
// There are three types of dependencies: build, runtime and intall
type Dependency struct {
	Build   []ComponentRef `json:"build"`
	Runtime []ComponentRef `json:"runtime"`
	Intall  []ComponentRef `json:"install"`
}

// Config is the toplevel struct that represents a configuration file
type Config struct {
	ComponentRef
	Deps Dependency
}

// RepoConfig defines the configuration for a repository
type RepoConfig struct {
	Type string `json:"type"`
	Base string `json:"base"`
}

// func (config Config) Fetch() {
// 	repo := "../test/" + config.Repo
// 	args := []string{"clone", repo, config.Name}
// 	git(args)
// }

func git(args []string) {
	log.Noticef("Executing: git %s\n", args)
	_, err := exec.Command("git", args...).Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			msg := string(ee.Stderr[:])
			log.Fatal("Error executing: ", msg)
		}

		log.Fatal(err)
	}
}

// Fetch the specified component
func (comp ComponentRef) Fetch() {
	var repo string
	if comp.Repoconfig != nil {
		repo += comp.Repoconfig.Base + comp.Repo
	} else {
		repo = comp.Repo
	}

	args := []string{"clone", repo, comp.Name}
	git(args)
}

func parseProjectFile(filename string) (*Config, error) {
	var data []byte
	data, err := ioutil.ReadFile(filename)
	if ee, ok := err.(*os.PathError); ok {
		log.Error("Error: ", ee)
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	return &config, err
}
