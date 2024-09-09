package main

import (
	"context"
	"fmt"

	"github.com/AlexGustafsson/cupdate/internal/source"
	"github.com/AlexGustafsson/cupdate/internal/source/k8s"
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
			entries, err := s.Entries(context.Background())
			if err != nil {
				return err
			}

			for _, e := range entries {
				// For now, let's focus on leaves - actually running pods. It's easier
				// as we don't have to deduplicate from templates
				switch o := e.Origin.(type) {
				case *k8s.Origin:
					if o.Container.Pod.IsTemplate {
						continue
					}
				}

				fmt.Printf("%s:%s\n", e.Image, e.Version)
				if e.ImageID != "" {
					fmt.Println("\tResolved:", e.ImageID)
				}

				switch o := e.Origin.(type) {
				case *k8s.Origin:
					fmt.Println("\tContainer name:", o.Container.Name)
					fmt.Println("\tIs template:", o.Container.Pod.IsTemplate)
					if o.Container.Pod.Name == "" {
						fmt.Println("\tPod name:", "-")
					} else {
						fmt.Println("\tPod name:", o.Container.Pod.Name)
					}
					if o.Container.Pod.Parent != nil {
						fmt.Println("\tParent:")
						fmt.Println("\t\tKind:", o.Container.Pod.Parent.ResourceKind)
						fmt.Println("\t\tNamespace:", o.Container.Pod.Parent.Namespace)
						fmt.Println("\t\tName:", o.Container.Pod.Parent.Name)
					}
				}
				fmt.Println()
			}

			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		panic(err.Error())
	}
}
