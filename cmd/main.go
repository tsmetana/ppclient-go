package main

import (
	"fmt"
	"os"

	"github.com/tsmetana/ppclient-go/pkg/ppclient"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Endpoint argument required")
		os.Exit(1)
	}
	endpoint := os.Args[1]
	client := ppclient.NewPpClient(endpoint)
	releases, err := client.GetReleases("openshift")
	if err != nil {
		fmt.Printf("Error getting release list: %v\n", err)
		os.Exit(1)
	}
	for _, r := range releases {
		fmt.Printf("Version: %s, phase %s, z-stream: %v\n", r.GetVersion(), r.GetPhase(), r.IsZStream())
	}

	latest := releases.GetLatestVersion(false)
	fmt.Printf("Latest Y-Stream: %s\n", latest)
	latest = releases.GetLatestVersion(true)
	fmt.Printf("Latest: %s\n", latest)
}
