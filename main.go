package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "up":
		upCmd := flag.NewFlagSet("up", flag.ExitOnError)
		configFile := upCmd.String("config", "config.yaml", "Path to YAML config file")
		force := upCmd.Bool("force", false, "Force recreate containers if they exist")
		upCmd.Parse(os.Args[2:])
		if err := startServices(*configFile, *force); err != nil {
			fmt.Fprintf(os.Stderr, "Start failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("All services started successfully")
	case "down":
		downCmd := flag.NewFlagSet("down", flag.ExitOnError)
		configFile := downCmd.String("config", "config.yaml", "Path to YAML config file")
		downCmd.Parse(os.Args[2:])
		if err := stopServices(*configFile); err != nil {
			fmt.Fprintf(os.Stderr, "Stop failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("All services stopped and removed")
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("Usage: localdev [up|down] --config <config.yaml>")
	fmt.Println("  up: Start services from config (add --force to recreate if exist)")
	fmt.Println("  down: Stop and remove services from config")
}
