package main

import (
    "flag"
    "fmt"
    "os"

    "github.com/BITVEL22/r-uqny/internal/node"
    "github.com/BITVEL22/r-uqny/internal/version"
)

func main() {
    versionFlag := flag.Bool("version", false, "show version")
    nodeID := flag.String("node", "", "create a node with the given ID")

    flag.Parse()

    if *versionFlag {
        fmt.Printf("%s %s\n", version.Name, version.Version)
        return
    }

    if *nodeID != "" {
        n := node.New(*nodeID)
        fmt.Printf("Node ID: %s\n", n.ID)
        return
    }

    fmt.Printf("%s %s\n", version.Name, version.Version)
    fmt.Println("Use -version to show version.")
    fmt.Println("Use -node <id> to create a node.")

    os.Exit(0)
}
