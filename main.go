package main

import (
	"fmt"

	"github.com/clusterService/server"
)

func main() {

	fmt.Printf("Initializing server...")
	server := &server.Server{}
	server.Initialize()

}

// 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
// remote k8 master ip->34.198.214.90//
// local minikube ip ->192.168.59.100
