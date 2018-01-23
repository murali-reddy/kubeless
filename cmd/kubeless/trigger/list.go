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
	"encoding/json"
	"fmt"
	"io"

	"github.com/ghodss/yaml"
	"github.com/gosuri/uitable"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	kubelessApi "github.com/kubeless/kubeless/pkg/apis/kubeless/v1beta1"
	"github.com/kubeless/kubeless/pkg/client/clientset/versioned"
	"github.com/kubeless/kubeless/pkg/utils"
)

var listCmd = &cobra.Command{
	Use:     "list FLAG",
	Aliases: []string{"ls"},
	Short:   "list all trigger deployed to Kubeless",
	Long:    `list all trigger deployed to Kubeless`,
	Run: func(cmd *cobra.Command, args []string) {
		output, err := cmd.Flags().GetString("out")
		if err != nil {
			logrus.Fatal(err.Error())
		}
		ns, err := cmd.Flags().GetString("namespace")
		if err != nil {
			logrus.Fatal(err.Error())
		}
		if ns == "" {
			ns = utils.GetDefaultNamespace()
		}

		kubelessClient, err := utils.GetKubelessClientOutCluster()
		if err != nil {
			logrus.Fatalf("Can not create out-of-cluster client: %v", err)
		}

		apiV1Client := utils.GetClientOutOfCluster()

		if err := doList(cmd.OutOrStdout(), kubelessClient, apiV1Client, ns, output, args); err != nil {
			logrus.Fatal(err.Error())
		}
	},
}

func init() {
	listCmd.Flags().StringP("out", "o", "", "Output format. One of: json|yaml")
	listCmd.Flags().StringP("namespace", "n", "", "Specify namespace for the function")
}

func doList(w io.Writer, kubelessClient versioned.Interface, apiV1Client kubernetes.Interface, ns, output string, args []string) error {
	var list []*kubelessApi.Trigger
	if len(args) == 0 {
		triggerList, err := kubelessClient.KubelessV1beta1().Triggers(ns).List(metav1.ListOptions{})
		if err != nil {
			return err
		}
		list = triggerList.Items
	} else {
		list = make([]*kubelessApi.Trigger, 0, len(args))
		for _, arg := range args {
			t, err := kubelessClient.KubelessV1beta1().Triggers(ns).Get(arg, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("Error listing trigger %s: %v", arg, err)
			}
			list = append(list, t)
		}
	}

	return printTriggers(w, list, apiV1Client, output)
}

// printprintTriggersFunctions formats the output of function list
func printTriggers(w io.Writer, triggers []*kubelessApi.Trigger, cli kubernetes.Interface, output string) error {
	if output == "" || output == "wide" {
		table := uitable.New()
		table.MaxColWidth = 50
		table.Wrap = true
		table.AddRow("NAME", "NAMESPACE", "TYPE", "TOPIC", "SCHEDULE", "FUNCTION", "STATUS")
		for _, t := range triggers {
			name := t.ObjectMeta.Name
			eventType := t.Spec.Type
			topic := t.Spec.Topic
			schedule := t.Spec.Schedule
			ns := t.ObjectMeta.Namespace
			functionName := t.Spec.FunctionName
			status, err := getDeploymentStatus(cli, functionName, ns)
			if err != nil && k8sErrors.IsNotFound(err) {
				status = "MISSING: Check controller logs"
			} else if err != nil {
				return err
			}
			table.AddRow(name, ns, eventType, topic, schedule, functionName, status)
		}
		fmt.Fprintln(w, table)
	} else {
		for _, trigger := range triggers {
			switch output {
			case "json":
				b, err := json.MarshalIndent(trigger, "", "  ")
				if err != nil {
					return err
				}
				fmt.Fprintln(w, string(b))
			case "yaml":
				b, err := yaml.Marshal(trigger)
				if err != nil {
					return err
				}
				fmt.Fprintln(w, string(b))
			default:
				return fmt.Errorf("Wrong output format. Please use only json|yaml")
			}
		}
	}
	return nil
}

func getDeploymentStatus(cli kubernetes.Interface, funcName, ns string) (string, error) {
	dpm, err := cli.ExtensionsV1beta1().Deployments(ns).Get(funcName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	status := fmt.Sprintf("%d/%d", dpm.Status.ReadyReplicas, dpm.Status.Replicas)
	if dpm.Status.ReadyReplicas > 0 {
		status += " READY"
	} else {
		status += " NOT READY"
	}
	return status, nil
}
