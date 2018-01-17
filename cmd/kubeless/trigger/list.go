/*
Copyright (c) 2016-2017 Bitnami

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package trigger

import (
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list FLAG",
	Aliases: []string{"ls"},
	Short:   "list all trigger deployed to Kubeless",
	Long:    `list all trigger deployed to Kubeless`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	listCmd.Flags().StringP("out", "o", "", "Output format. One of: json|yaml")
	listCmd.Flags().StringP("namespace", "n", "", "Specify namespace for the function")
}
