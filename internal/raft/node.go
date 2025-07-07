package raft

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	Follower  = "follower"
	Candidate = "candidate"
	Leader    = "leader"
)

type VoteRequest struct {
	Term        int
	CandidateID int
}

type VoteResponse struct {
	VoteGranted bool
}

type HeartbeatRequest struct {
	LeaderID int
	Term     int
}

type Node struct {
	ID          int
	Addr        string
	Peers       []string
	mu          sync.Mutex
	State       string // "follower", "candidate", "leader"
	CurrentTerm int
	VotedFor    int
	LeaderID    int
}

func NewNode(id int, addr string, peers []string) *Node {
	return &Node{
		ID:    id,
		Addr:  addr,
		Peers: peers,
		State: Follower,
	}
}

func (n *Node) logLeaderPeriodically() {
	for {
		n.mu.Lock()
		leader := "unknown"
		if n.State == Leader {
			leader = fmt.Sprintf("Node %d (myself)", n.ID)
		} else if n.LeaderID != 0 {
			leader = fmt.Sprintf("Node %d", n.LeaderID)
		}
		log.Printf("[Node %d] Current leader: %s", n.ID, leader)
		n.mu.Unlock()
		time.Sleep(3 * time.Second)
	}
}

func (n *Node) Start() {
	log.Printf("Node %d started at %s", n.ID, n.Addr)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from node %d!", n.ID)
	})
	http.HandleFunc("/request-vote", n.handleVoteRequest)
	http.HandleFunc("/heartbeat", n.handleHeartbeat)
	go n.runElectionTimer()
	go n.logLeaderPeriodically()
	log.Fatal(http.ListenAndServe(n.Addr, nil))
}

func (n *Node) handleVoteRequest(w http.ResponseWriter, r *http.Request) {
	var req VoteRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	granted := false

	n.mu.Lock()
	if req.Term > n.CurrentTerm {
		n.CurrentTerm = req.Term
		n.State = Follower
		n.VotedFor = 0
	}
	if req.Term == n.CurrentTerm && (n.VotedFor == 0 || n.VotedFor == req.CandidateID) {
		n.VotedFor = req.CandidateID
		granted = true
	}
	n.mu.Unlock()

	_ = json.NewEncoder(w).Encode(VoteResponse{VoteGranted: granted})
}

func (n *Node) handleHeartbeat(w http.ResponseWriter, r *http.Request) {
	var req HeartbeatRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	n.mu.Lock()
	if req.Term > n.CurrentTerm {
		n.CurrentTerm = req.Term
		n.State = Follower
		n.VotedFor = 0
	}
	if req.LeaderID != n.LeaderID {
		n.LeaderID = req.LeaderID
	}
	n.State = Follower
	n.mu.Unlock()
	select {
	case heartbeatChan <- struct{}{}:
	default:
	}
	w.WriteHeader(http.StatusOK)
}
