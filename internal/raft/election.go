package raft

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var heartbeatReceived = false
var heartbeatChan = make(chan struct{})

func (n *Node) runElectionTimer() {
	for {
		timeout := time.Duration(800+rand.Intn(800)) * time.Millisecond
		timer := time.NewTimer(timeout)
		select {
		case <-timer.C:
			n.mu.Lock()
			if !heartbeatReceived && n.State != Leader {
				n.LeaderID = 0
				n.startElection()
			}
			heartbeatReceived = false
			n.mu.Unlock()

		case <-heartbeatChan:
			heartbeatReceived = true
			timer.Stop()

		}
	}
}

func (n *Node) startElection() {
	n.State = Candidate
	n.CurrentTerm++
	n.VotedFor = n.ID
	votes := 1
	log.Printf("[Node %d] Starting election for term %d", n.ID, n.CurrentTerm)

	for _, peer := range n.Peers {
		req := VoteRequest{Term: n.CurrentTerm, CandidateID: n.ID}
		buf := new(bytes.Buffer)
		_ = json.NewEncoder(buf).Encode(req)
		resp, err := http.Post(peer+"/request-vote", "application/json", buf)
		if err != nil {
			log.Printf("[Node %d] Error sending vote request to %s: %v", n.ID, peer, err)
			continue
		}
		var voteResp VoteResponse
		err = json.NewDecoder(resp.Body).Decode(&voteResp)
		if err != nil {
			log.Printf("[Node %d] Error decoding vote response from %s: %v", n.ID, peer, err)
			continue
		}
		if voteResp.VoteGranted {
			votes++
			log.Printf("[Node %d] Received vote from %s", n.ID, peer)
		} else {
			log.Printf("[Node %d] Vote denied by %s", n.ID, peer)
		}
	}
	log.Printf("[Node %d] Total votes: %d (need >%d)", n.ID, votes, len(n.Peers)/2)
	if votes > len(n.Peers)/2 {
		n.State = Leader
		log.Printf("[Node %d] Became leader for term %d", n.ID, n.CurrentTerm)
		go n.sendHeartbeats()
	}
}

func (n *Node) sendHeartbeats() {
	for {
		n.mu.Lock()
		if n.State != Leader {
			n.mu.Unlock()
			return
		}
		leaderID := n.ID
		n.mu.Unlock()
		for _, peer := range n.Peers {
			req := HeartbeatRequest{LeaderID: leaderID, Term: n.CurrentTerm}
			buf := new(bytes.Buffer)
			_ = json.NewEncoder(buf).Encode(req)
			_, _ = http.Post(peer+"/heartbeat", "application/json", buf)
		}
		time.Sleep(50 * time.Millisecond)
	}
}
