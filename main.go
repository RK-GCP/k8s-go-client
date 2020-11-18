package main

import (
	"fmt"
	"log"
	"os"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// This program lists the pods in a cluster equivalent to
//
// kubectl get pods
//
func main() {
	rNs := os.Args[1]

	if rNs == "" {
		fmt.Println("no namespace requested")
	}

	// Bootstrap k8s configuration from local 	Kubernetes config file
	//kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	kubeconfig := os.Getenv("KUBECONFIG")
	log.Println("Using kubeconfig file: ", kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	// Create an rest client not targeting specific API version
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	//create namespace
	nsNew, err := clientset.CoreV1().Namespaces().Create(requestNamepace(rNs))
	if err != nil {
		if errors.IsAlreadyExists(err) {
			fmt.Println("namespace exist", rNs)
		} else {
			log.Fatalln("error creating namespace: ", err)
		}
	} else {
		fmt.Println("namespace created", nsNew.Name)
	}

	//get all namespaces
	allns, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to get namespace:", err)
	}

	//show namespaces
	for n, nns := range allns.Items {
		fmt.Printf("[%d] %s\n", n, nns.Name)
	}
}

func requestNamepace(namepace string) *v1.Namespace {
	if namepace == "" {
		return nil
	}

	return &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namepace,
		},
	}
}
