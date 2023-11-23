package main

import (
	"encoding/json"
	// "fmt"
	"sync"

	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type TopologyMessageBody struct {
	Type string `json:"type"`
	Topology map[string][]string `json:"topology"`
}

func main() {
	n := maelstrom.NewNode()
	s := &Server{n: n, ids: make(map[float64]struct{}), neighbours: make([]string, 0)}

	n.Handle("broadcast", s.handle_broadcast)
	n.Handle("read", s.handle_read)
	n.Handle("topology", s.handle_topology)

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

/// State of the node
type Server struct {
	n *maelstrom.Node
	ids map[float64]struct{}
	ids_mu sync.RWMutex
	neighbours []string
}

/// Handler for broadcast message
func (s *Server) handle_broadcast(msg maelstrom.Message) error {
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}	
	// Get the message value
	id := body["message"].(float64)
	s.ids_mu.Lock()
	if _, exists := s.ids[id]; exists {
		s.ids_mu.Unlock()
		return nil
	}
	s.ids[id] = struct{}{}
	// Propogate broadcast to neighbours
	for _, neighbour := range s.neighbours {
		neighbour := neighbour
		go func() {
			_ = s.n.Send(neighbour, body)
		}()
	}
	s.ids_mu.Unlock()
	// Return
	return s.n.Reply(msg, map[string]any{
		"type": "broadcast_ok",
	})
}

/// Handler for read message
func (s *Server) handle_read(msg maelstrom.Message) error {
	// Build the id list
	s.ids_mu.RLock()
	id_list := make([]float64, 0, len(s.ids))
	for id := range s.ids {
		id_list = append(id_list, id)
	}
	s.ids_mu.RUnlock()
	// Return
	return s.n.Reply(msg, map[string]any{
		"type": "read_ok",
		"messages": id_list,
	})
}

/// Handler for topology message
func (s *Server) handle_topology(msg maelstrom.Message) error {
	// Parse topology field
	var body TopologyMessageBody
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}
	s.neighbours = body.Topology[s.n.ID()]
	// Return
	return s.n.Reply(msg, map[string]any{
		"type": "topology_ok",
	});
}