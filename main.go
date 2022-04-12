package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	civo "github.com/civo/civogo"
	"github.com/digitalocean/godo"
	linode "github.com/linode/linodego"

	"golang.org/x/oauth2"
)

func main() {
	for {
		selection := StringPrompt("\n(1) Civo\n(2) Linode\n(3) Digital Ocean\n(q) quit\n\nChoose provider: ")
		if selection == "q" {
			return
		}

		index, err := strconv.Atoi(selection)
		if err != nil {
			log.Fatalf("invalid entry %v selection must be a number", selection)
		}

		fmt.Println()

		switch index {
		case 1:
			kaasKillCivo()
		case 2:
			kaasKillLinode()
		case 3:
			kaasKillDigitalOcean()
		default:
			log.Fatalf("invalid entry %v selection must be a number", selection)
		}
	}
}

// StringPrompt asks for a string value using the label
func StringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label)
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

func kaasKillCivo() {
	apikey, ok := os.LookupEnv("CIVO_API_KEY")
	if !ok {
		log.Fatalln("missing required CIVO_API_KEY environment variable")
	}

	client, err := civo.NewClient(apikey, "")
	if err != nil {
		log.Fatalf("failed creating new Civo client: %v", err)
	}

	clusters, err := client.ListKubernetesClusters()
	if err != nil {
		log.Fatalf("failed listing all instances: %v", err)
	}
	if len(clusters.Items) == 0 {
		fmt.Println("There are no clusters currently under this account")
		return
	}

	for {
		// Print a list of clusters and ask which ones to delete.
		indexedClusters := make(map[int]string, len(clusters.Items))
		for i, cluster := range clusters.Items {
			fmt.Printf("(%v) %v\n", i, cluster.Name)
			indexedClusters[i] = cluster.ID
		}
		fmt.Printf("(b) go back\n\n")
		selection := StringPrompt("Enter cluster index number to remove: ")
		if selection == "b" {
			return
		}

		// Try to delete selected cluster index number.
		index, err := strconv.Atoi(selection)
		if err != nil {
			log.Fatalf("invalid entry %v selection must be a number", selection)
		}

		id, ok := indexedClusters[index]
		if !ok {
			log.Fatalf("invalid cluster index: should be between 0-%v", len(clusters.Items)-1)
		}
		client.DeleteKubernetesCluster(id)
		fmt.Println("Successfully deleted cluster", clusters.Items[index].Name)
	}
}

func kaasKillLinode() {
	apikey, ok := os.LookupEnv("LINODE_API_KEY")
	if !ok {
		log.Fatalln("missing required LINODE_API_KEY environment variable")
	}

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: apikey})

	oauth2Client := &http.Client{
		Transport: &oauth2.Transport{
			Source: tokenSource,
		},
	}

	client := linode.NewClient(oauth2Client)

	for {
		clusters, err := client.ListLKEClusters(context.Background(), nil)
		if err != nil {
			log.Fatalf("failed listing all instances: %v", err)
		}

		if len(clusters) == 0 {
			fmt.Println("There are no clusters currently under this account")
			return
		}

		for i, cluster := range clusters {
			fmt.Printf("(%v) %v\n", i+1, cluster.Label)
		}
		fmt.Printf("(b) go back\n\n")
		selection := StringPrompt("Enter cluster index number to remove: ")
		if selection == "b" {
			return
		}

		// Try to delete selected cluster index number.
		index, err := strconv.Atoi(selection)
		if err != nil {
			log.Fatalf("invalid entry %v selection must be a number", selection)
		}

		err = client.DeleteLKECluster(context.Background(), clusters[index-1].ID)
		if err != nil {
			log.Fatalf("failed deleting instance: %v", err)
		}
	}
}

func kaasKillDigitalOcean() {
	apikey, ok := os.LookupEnv("DIGITAL_OCEAN_API_KEY")
	if !ok {
		log.Fatalln("missing required DIGITAL_OCEAN_API_KEY environment variable")
	}

	client := godo.NewFromToken(apikey)
	ctx := context.TODO()

	for {
		clusters, _, err := client.Kubernetes.List(ctx, nil)
		if err != nil {
			log.Fatalf("failed listing all instances: %v", err)
		}

		if len(clusters) == 0 {
			fmt.Println("There are no clusters currently under this account")
			return
		}

		for i, cluster := range clusters {
			fmt.Printf("(%v) %v\n", i+1, cluster.Name)
		}
		fmt.Printf("(b) go back\n\n")
		selection := StringPrompt("Enter cluster index number to remove: ")
		if selection == "b" {
			return
		}

		// Try to delete selected cluster index number.
		// index, err := strconv.Atoi(selection)
		// if err != nil {
		// 	log.Fatalf("invalid entry %v selection must be a number", selection)
		// }

		// err = client.DeleteLKECluster(context.Background(), clusters[index-1].ID)
		// if err != nil {
		// 	log.Fatalf("failed deleting instance: %v", err)
		// }
	}
}
