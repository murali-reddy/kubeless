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
	kubelessApi "github.com/kubeless/kubeless/pkg/apis/kubeless/v1beta1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"github.com/kubeless/kubeless/pkg/utils"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var createCmd = &cobra.Command{
	Use:   "create <trigger_name> FLAG",
	Short: "Create a trigger to Kubeless",
	Long:  `Create a trigger to Kubeless`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			logrus.Fatal("Need exactly one argument - trigger name")
		}
		triggerName := args[0]

		triggerHTTP, err := cmd.Flags().GetBool("trigger-http")
		if err != nil {
			logrus.Fatal(err)
		}

		schedule, err := cmd.Flags().GetString("schedule")
		if err != nil {
			logrus.Fatal(err)
		}

		if schedule != "" {
			if _, err := cron.ParseStandard(schedule); err != nil {
				logrus.Fatalf("Invalid value for --schedule. " + err.Error())
			}
		}

		ns, err := cmd.Flags().GetString("namespace")
		if err != nil {
			logrus.Fatal(err)
		}
		if ns == "" {
			ns = utils.GetDefaultNamespace()
		}

		topic, err := cmd.Flags().GetString("trigger-topic")
		if err != nil {
			logrus.Fatal(err)
		}

		functionName, err := cmd.Flags().GetString("function-name")
		if err != nil {
			logrus.Fatal(err)
		}

		headless, err := cmd.Flags().GetBool("headless")
		if err != nil {
			logrus.Fatal(err)
		}

		port, err := cmd.Flags().GetInt32("port")
		if err != nil {
			logrus.Fatal(err)
		}
		if port <= 0 || port > 65535 {
			logrus.Fatalf("Invalid port number %d specified", port)
		}

		kubelessClient, err := utils.GetKubelessClientOutCluster()
		if err != nil {
			logrus.Fatalf("Can not create out-of-cluster client: %v", err)
		}
		funcObj, err := utils.GetFunction(kubelessClient, functionName, ns)
		if err != nil {
			logrus.Fatalf("Unable to find the function %s in the namespace %s. Received %s: ", functionName, ns, err)
		}

		svcSpec := v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     "http-function-port",
					NodePort: 0,
					Protocol: v1.ProtocolTCP,
				},
			},
			Selector: funcObj.ObjectMeta.Labels,
			Type:     v1.ServiceTypeClusterIP,
		}

		if headless {
			svcSpec.ClusterIP = v1.ClusterIPNone
		}

		if port != 0 {
			svcSpec.Ports[0].Port = port
			svcSpec.Ports[0].TargetPort = intstr.FromInt(int(port))
		}

		trigger := kubelessApi.Trigger{}
		trigger.TypeMeta = metav1.TypeMeta{
			Kind:       "Trigger",
			APIVersion: "kubeless.io/v1beta1",
		}

		switch {
		case triggerHTTP:
			trigger.Spec.Type = "HTTP"
			trigger.Spec.Topic = ""
			trigger.Spec.Schedule = ""
			break
		case schedule != "":
			trigger.Spec.Type = "Scheduled"
			trigger.Spec.Schedule = schedule
			trigger.Spec.Topic = ""
			break
		case topic != "":
			trigger.Spec.Type = "PubSub"
			trigger.Spec.Topic = topic
			trigger.Spec.Schedule = ""
			break
		}
		trigger.Spec.ServiceSpec = svcSpec
		trigger.Spec.FunctionName = functionName
		trigger.ObjectMeta = metav1.ObjectMeta{
			Name:      triggerName,
			Namespace: ns,
			Labels: map[string]string{
				"created-by": "kubeless",
			},
		}

		logrus.Infof("Deploying trigger...")
		err = utils.CreateTriggerResource(kubelessClient, &trigger)
		if err != nil {
			logrus.Fatalf("Failed to deploy %s. Received:\n%s", triggerName, err)
		}
		logrus.Infof("Trigger %s submitted for deployment", triggerName)
		logrus.Infof("Check the deployment status executing 'kubeless trigger ls %s'", triggerName)
	},
}

func init() {
	createCmd.Flags().StringP("namespace", "", "", "Specify namespace for the function")
	createCmd.Flags().StringP("trigger-topic", "", "", "Deploy a pubsub function to Kubeless")
	createCmd.Flags().StringP("schedule", "", "", "Specify schedule in cron format for scheduled function")
	createCmd.Flags().Bool("trigger-http", false, "Deploy a http-based function to Kubeless")
	createCmd.Flags().StringP("function-name", "", "", "Name of the function to be associated with trigger")
	createCmd.Flags().Bool("headless", false, "Deploy http-based function without a single service IP and load balancing support from Kubernetes. See: https://kubernetes.io/docs/concepts/services-networking/service/#headless-services")
	createCmd.Flags().Int32("port", 8080, "Deploy http-based function with a custom port")
	createCmd.MarkFlagRequired("function-name")
}
