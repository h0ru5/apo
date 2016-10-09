// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/foomo/htpasswd"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list users in the file",
	Long:  `listing all users`,
	Run: func(cmd *cobra.Command, args []string) {
		file := cmd.Flag("file").Value.String()
		fmt.Print("getting passes from file ", file, " ...")
		passwords, err := htpasswd.ParseHtpasswdFile(file)
		if err == nil {
			fmt.Print("Success!\n")
		} else {
			panic(err)
		}
		fmt.Println("Found users:")
		for k := range passwords {
			fmt.Println(k)
		}

	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
