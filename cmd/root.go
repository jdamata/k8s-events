package cmd

import (
	"encoding/json"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	namespaces []string
	kubeconfig string
	rootCmd    = &cobra.Command{
		Use:   "k8s-events",
		Short: "Logs kubernetes events",
		Long:  "Logs kubernetes events",
		Run:   main,
	}
)

// Execute executes the root command.
func Execute(version string) error {
	rootCmd.Version = version
	homedir, _ := os.UserHomeDir()
	rootCmd.PersistentFlags().StringSlice("namespaces", namespaces, "List of namespaces for event grabbing")
	viper.BindPFlag("namespaces", rootCmd.Flags().Lookup("namespaces"))
	rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", homedir+"/.kube/config", "(optional) absolute path to the kubeconfig file")
	viper.BindPFlag("kubeconfig", rootCmd.Flags().Lookup("kubeconfig"))
	return rootCmd.Execute()
}

// If namespaces is empty list, default to all namespaces
func ns(namespaces []string, clientset *kubernetes.Clientset) []string {
	if len(namespaces) == 0 {
		log.Debug("Namespaces flag not provided. Defaulting to all namespaces")
		namespaceList, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
		if err != nil {
			log.Fatal("Failed to grab list of namespaces: ", err)
		}
		for _, namespace := range namespaceList.Items {
			namespaces = append(namespaces, namespace.Name)
		}
	}
	return namespaces
}

// Authenticate to cluster
func clientset() *kubernetes.Clientset {
	if _, err := os.Stat(kubeconfig); err != nil {
		log.Debug("Kubeconfig not provided")
		kubeconfig = ""
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal("Cannot build kubeconfig for authentication: ", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal("Cannot create kubernetes client: ", err)
	}
	return clientset
}

// Grab kubernetes events from a namespace
func events(namespace string, clientset *kubernetes.Clientset) {
	log.Info("Starting watch on namespace: ", namespace)
	events := clientset.CoreV1().Events(namespace)
	watch, err := events.Watch(metav1.ListOptions{Watch: true})
	if err != nil {
		log.Error("Cannot create watch interface on namespace: ", namespace, err)
	}
	for {
		var event corev1.Event
		results := <-watch.ResultChan()
		data, _ := json.Marshal(results.Object)
		json.Unmarshal(data, &event)
		log.WithFields(log.Fields{
			"namespace":      namespace,
			"pod":            event.Name,
			"reason":         event.Reason,
			"event_level":    event.Type,
			"count":          event.Count,
			"FirstTimestamp": event.FirstTimestamp,
			"LastTimestamp":  event.LastTimestamp,
			//"event_time":     event.EventTime,
			//"controller": event.ReportingController,
		}).Info(event.Message)
	}
}

func main(cmd *cobra.Command, args []string) {
	log.SetFormatter(&log.JSONFormatter{})
	clientset := clientset()
	var wg sync.WaitGroup
	for _, namespace := range ns(namespaces, clientset) {
		wg.Add(1)
		go events(namespace, clientset)
	}
	wg.Wait()
}
