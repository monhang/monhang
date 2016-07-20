// Copyright 2016 Thiago Cangussu de Castro Gomes. All rights reserved.
// Use of this source code is governed by a GNU General Public License
// version 3 that can be found in the LICENSE file.

package main

var cmdBoot = &Command{
	Name:  "boot",
	Args:  "[configfile]",
	Short: "bootstrap a component and its dependencies",
	Long: `
Boot fetches and setups the workspace for the component described in the given configuration file.
`,
	OnlyInGOPATH: true,
}

var bootF = cmdBoot.Flag.String("f", "<defaultconfig>", "configuration file")

func findBootDesc() string {
	if *bootF != "<defaultconfig>" {
		return *bootF
	}
	return "./monhang.json"
}

func runBoot(cmd *Command, args []string) {
	config, err := parseProjectFile(findBootDesc())
	if err != nil {
		check(err)
	}

	config.Fetch() // fetch toplevel component

	// Fetch build dependencies
	for _, dep := range config.Deps.Build {
		log.Debug("Processing dependency ", dep.Name)
		if dep.Repoconfig == nil {
			log.Debug("Adding toplevel repoconfig to dep:", *config.Repoconfig)
			dep.Repoconfig = config.Repoconfig
		}
		dep.Fetch()
	}
}

func init() {
	cmdBoot.Run = runBoot // break init loop
}
