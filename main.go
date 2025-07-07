package main

import (
	"flag"
	"fmt"
	raft "node/internal/raft"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	// Önce environment değişkenlerini oku
	idEnv := os.Getenv("NODE_ID")
	addrEnv := os.Getenv("NODE_ADDR")
	peersEnv := os.Getenv("NODE_PEERS")

	id := 1
	if idEnv != "" {
		if parsed, err := strconv.Atoi(idEnv); err == nil {
			id = parsed
		}
	}

	addr := ":8001"
	if addrEnv != "" {
		addr = addrEnv
	}

	var peers []string
	if peersEnv != "" {
		for _, p := range strings.Split(peersEnv, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				peers = append(peers, p)
			}
		}
	}

	// Eğer environment yoksa eski flag'leri kullan
	if idEnv == "" || addrEnv == "" || peersEnv == "" {
		flagId := flag.Int("id", 1, "Node ID")
		port := flag.Int("port", 8080, "Port number")
		flag.Parse()
		host := "localhost"
		addr = fmt.Sprintf("%s:%d", host, *port)
		allPeers := []string{"localhost:8001", "localhost:8002", "localhost:8003"}
		peers = nil
		for _, p := range allPeers {
			if p != addr {
				peers = append(peers, p)
			}
		}
		id = *flagId
	}

	node := raft.NewNode(id, addr, peers)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go node.Start()
	fmt.Printf("Raft node %d is running at %s...\n", id, addr)
	fmt.Printf("Peers: %v\n", peers)

	<-sigs
	fmt.Println("Node shutting down...")
}
