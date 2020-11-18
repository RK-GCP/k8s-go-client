package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// This program lists the pods in a cluster equivalent to
//
// kubectl get pods
//
func main() {
	var ns string
	flag.StringVar(&ns, "namespace", "", "namespace")

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

	nsNew, err := clientset.CoreV1().Namespaces().Create(requestNamepace())
	if err != nil {
		log.Fatalln("error creating namespace: ", err)
	}
	fmt.Println("namespace created", nsNew.Name)

	allns, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to get namespace:", err)
	}

	for n, nns := range allns.Items {
		fmt.Printf("[%d] %s\n", n, nns.Name)
	}

	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to get pods:", err)
	}

	// print pods
	for i, pod := range pods.Items {
		fmt.Printf("[%d] %s\n", i, pod.GetName())
	}
}

func requestNamepace() *v1.Namespace {
	return &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "gc-namepace",
		},
	}
}
