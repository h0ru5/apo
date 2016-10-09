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

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "<user> <file> -- remove a user from the file",
	Long:  `delete a username from the file`,
	Run: func(cmd *cobra.Command, args []string) {
		user := cmd.Flag("user").Value.String()
		file := cmd.Flag("file").Value.String()
		fmt.Println("rm called on file ", file, " to remove user ", user)
		err := htpasswd.RemoveUser(file, user)
		if err == nil {
			fmt.Println("Success")
		} else {
			panic(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(rmCmd)
	rmCmd.PersistentFlags().StringP("user", "u", "", "user name")
}
