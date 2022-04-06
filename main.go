package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/civo/civogo"
)

func main() {
	civoKey, ok := os.LookupEnv("CIVO_API_KEY")
	if !ok {
		log.Fatalln("missing required CIVO_API_KEY environment variable")
	}

	client, err := civogo.NewClient(civoKey, "")
	if err != nil {
		log.Fatalf("failed creating new Civo client: %v", err)
	}

	clusters, err := client.ListKubernetesClusters()
	if err != nil {
		log.Fatalf("failed listing all instances: %v", err)
	}

	// Print a list of clusters and ask which ones to delete.
	indexedClusters := make(map[int]string, len(clusters.Items))
	for i, cluster := range clusters.Items {
		fmt.Printf("%v: %v\n", i, cluster.Name)
		indexedClusters[i] = cluster.ID
	}
	selection := StringPrompt("Enter cluster index number to remove: ")

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
