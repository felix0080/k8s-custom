package main

import (
	"flag"
	"github.com/felix0080/k8s-custom/pkg/signals"
	"k8s.io/client-go/tools/clientcmd"
	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
	clientset "github.com/felix0080/k8s-custom/pkg/client/clientset/versioned"
	informers "github.com/felix0080/k8s-custom/pkg/client/informers/externalversions"
	"time"
)

var(
	masterURL string
	kubeconfig string
)
func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}

func main() {
	flag.Parse()
	// set up signals so we handle the first shutdown signal gracefully
	stopCh:=signals.SetupSignalHandler()
	config,err:=clientcmd.BuildConfigFromFlags(masterURL,kubeconfig)
	if err != nil {
		glog.Errorf("Error building kubeconfig: %s",err.Error())
	}
	kubeClient,err:=kubernetes.NewForConfig(config)
	if err != nil {
		glog.Errorf("Error building kubernetes clientset: %s",err.Error())
	}
	networkClient,err:=clientset.NewForConfig(config)
	if err != nil {
		glog.Fatalf("Error building example clientset: %s", err.Error())
	}
	networkInformerFactory := informers.NewSharedInformerFactory(networkClient,time.Second*30)

	controller:=NewController(kubeClient,networkClient,networkInformerFactory.Samplecrd().V1().Networks())
	go networkInformerFactory.Start(stopCh)

	err=controller.Run(2,stopCh)
	if err != nil {
		glog.Errorf("Error running controller: %s", err.Error())
	}

}