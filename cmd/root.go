package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Error(err)
	}
	rootCmd.PersistentFlags().StringSlice("namespaces", namespaces, "List of namespaces for event grabbing")
	viper.BindPFlag("namespaces", rootCmd.Flags().Lookup("namespaces"))
	rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", homedir+"/.kube/config", "(optional) absolute path to the kubeconfig file")
	viper.BindPFlag("kubeconfig", rootCmd.Flags().Lookup("kubeconfig"))
	return rootCmd.Execute()
}

// If namespaces is empty list, default to all namespaces
func ns(namespaces []string, clientset *kubernetes.Clientset) []string {
	if len(namespaces) == 0 {
		namespaceList, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
		if err != nil {
			log.Error("namespaces flag not provided and cannot grab list of namespaces", err)
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
		kubeconfig = ""
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Error(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error(err)
	}
	return clientset
}

func main(cmd *cobra.Command, args []string) {
	log.SetFormatter(&log.JSONFormatter{})
	clientset := clientset()
	for _, namespace := range ns(namespaces, clientset) {
		events := clientset.CoreV1().Events(namespace)
		watch, _ := events.Watch(metav1.ListOptions{})
		results := <-watch.ResultChan()
		log.Info(results)
	}
}
