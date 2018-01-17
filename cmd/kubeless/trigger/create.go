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

var createCmd = &cobra.Command{
	Use:   "deploy <function_name> FLAG",
	Short: "deploy a trigger to Kubeless",
	Long:  `deploy a trigger to Kubeless`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	createCmd.Flags().StringP("namespace", "", "", "Specify namespace for the function")
	createCmd.Flags().StringP("trigger-topic", "", "", "Deploy a pubsub function to Kubeless")
	createCmd.Flags().StringP("schedule", "", "", "Specify schedule in cron format for scheduled function")
	createCmd.Flags().Bool("trigger-http", false, "Deploy a http-based function to Kubeless")
}
