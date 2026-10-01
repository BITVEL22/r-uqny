package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/BITVEL22/r-uqny/internal/config"
	"github.com/BITVEL22/r-uqny/internal/identity"
	"github.com/BITVEL22/r-uqny/internal/logging"
	"github.com/BITVEL22/r-uqny/internal/node"
	"github.com/BITVEL22/r-uqny/internal/version"
)

func main() {
	versionFlag := flag.Bool("version", false, "show version")
	listenAddr := flag.String("listen", config.DefaultListenAddress, "listen address")

	flag.Parse()

	if *versionFlag {
		fmt.Printf("%s %s\n", version.Name, version.Version)
		return
	}

	logging.Init(slog.LevelInfo)

	id, err := identity.Generate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate identity: %v\n", err)
		os.Exit(1)
	}

	cfg := config.Default()
	cfg.NodeID = id.NodeID
	cfg.ListenAddress = *listenAddr

	n := node.New(cfg, id)

	if err := n.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to start node: %v\n", err)
		os.Exit(1)
	}

	slog.Info("node started",
		"node_id", n.ID,
		"listen", n.Listener.Address().String(),
		"version", version.Version,
	)

	go func() {
		if err := n.Serve(); err != nil {
			slog.Error("serve error", "error", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	slog.Info("shutting down")

	if err := n.Stop(); err != nil {
		slog.Error("failed to stop node", "error", err)
		os.Exit(1)
	}

	slog.Info("node stopped")
}

