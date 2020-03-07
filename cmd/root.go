package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
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
	rootCmd.PersistentFlags().StringVarP(&kubeconfig, "kubeconfig", "k", homedir+"/.kube/config", "(optional) absolute path to the kubeconfig file")
	viper.BindPFlag("kubeconfig", rootCmd.Flags().Lookup("kubeconfig"))
	return rootCmd.Execute()
}

func logging() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
}

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
	clientset := clientset()
	events := clientset.CoreV1().Events("default")
	opts := v1.ListOptions{}
	watchInt, err := events.Watch(opts)
	if err != nil {
		log.Error(err)
	}
	results := watchInt.ResultChan()
	fmt.Println(results)
}
