/*
   Monhang - component management tool
   Copyright (C) 2016  Thiago Cangussu de Castro Gomes

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"flag"
	"fmt"
	"github.com/op/go-logging"
	"os"
)

var log = logging.MustGetLogger("monhang")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func setupLog() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backend)
	backendLeveled.SetLevel(logging.ERROR, "")
	logging.SetBackend(backendFormatter)
}

// Command is an implementation of a godep command
// like godep save or godep go.
type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(cmd *Command, args []string)

	// Name of the command
	Name string

	// Args the command would expect
	Args string

	// Short is the short description shown in the 'godep help' output.
	Short string

	// Long is the long message shown in the
	// 'godep help <this-command>' output.
	Long string

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet

	// OnlyInGOPATH limits this command to being run only while inside of a GOPATH
	OnlyInGOPATH bool
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func version() {
	fmt.Println("monhang v0.0.1")
}

func usageExit() {
	version()
	fmt.Println(`
Usage:

	monhang command [arguments]

The commands are:

	boot        bootstraps a workspace
	version     print monhang version

Use "monhang help [command]" for more information about a command.
`)
	os.Exit(0)
}

var cmdHelp = &Command{
	Name: "help",
	Run: func(cmd *Command, args []string) {
		// TODO(cangussu): print the help for the command given in args
		usageExit()
	},
}

var commands = []*Command{
	cmdBoot,
	cmdHelp,
}

func init() {
	setupLog()
}

func main() {
	flag.Usage = usageExit
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("You must tell monhang what to do!")
		usageExit()
	}

	for _, cmd := range commands {
		if cmd.Name == args[0] {
			cmd.Flag.Parse(args[1:])
			cmd.Run(cmd, cmd.Flag.Args())
		}
	}
}
