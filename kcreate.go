package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/tools/clientcmd"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	f, err := os.Open("manifest.yaml")
	if err != nil {
		panic(err.Error())
	}

	home := os.Getenv("HOME")
	configpath := filepath.Join(home, ".kube/config")
	config, err := clientcmd.BuildConfigFromFlags("", configpath)
	if err != nil {
		panic(err)
	}

	scheme := runtime.NewScheme()
	client, err := crclient.New(config, crclient.Options{Scheme: scheme})
	if err != nil {
		panic(err)
	}

	decoder := yaml.NewYAMLOrJSONDecoder(f, 65536)
	for {
		u := unstructured.Unstructured{}
		err = decoder.Decode(&u)
		if err == io.EOF {
			return
		}
		if err != nil {
			panic(err.Error())
		}
		u.SetNamespace("default")

		err = client.Create(context.TODO(), &u)
		if err != nil {
			panic(err)
		}
		fmt.Printf("created %v\n", u.GroupVersionKind())
	}
}
