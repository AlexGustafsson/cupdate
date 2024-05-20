package main

import (
	"context"
	"fmt"

	"github.com/AlexGustafsson/k8s-image-feed/internal/source"
	"github.com/AlexGustafsson/k8s-image-feed/internal/source/k8s"
	"golang.org/x/sync/errgroup"
	"k8s.io/client-go/rest"
)

func main() {
	// config, err := rest.InClusterConfig()
	// if err != nil {
	// 	panic(err.Error())
	// }
	config := &rest.Config{
		Host: "http://localhost:8001",
	}

	k8sSource, err := k8s.New(config)
	if err != nil {
		panic(err.Error())
	}

	sources := []source.Source{k8sSource}

	var wg errgroup.Group
	for _, s := range sources {
		s := s
		wg.Go(func() error {
			return s.EachListItem(context.Background(), func(e source.Entry) error {
				fmt.Printf("%s@%s\n", e.Image, e.Version)

				switch o := e.Origin.(type) {
				case *k8s.Origin:
					fmt.Println("\tKind:", o.ResourceKind)
					fmt.Println("\tNamespace:", o.Namespace)
					fmt.Println("\tName:", o.Name)
					fmt.Println("\tContainerName:", o.ContainerName)
					fmt.Println("\tCreated:", o.Created)
				}
				fmt.Println()

				return nil
			})
		})
	}

	if err := wg.Wait(); err != nil {
		panic(err.Error())
	}
}
