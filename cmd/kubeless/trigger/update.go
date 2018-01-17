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

var updateCmd = &cobra.Command{
	Use:   "update <function_name> FLAG",
	Short: "update a trigger on Kubeless",
	Long:  `update a trigger on Kubeless`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	updateCmd.Flags().StringP("namespace", "", "", "Specify namespace for the function")
	updateCmd.Flags().StringP("trigger-topic", "", "", "Deploy a pubsub function to Kubeless")
	updateCmd.Flags().StringP("schedule", "", "", "Specify schedule in cron format for scheduled function")
	updateCmd.Flags().Bool("trigger-http", false, "Deploy a http-based function to Kubeless")
}
