package main

import (
	"encoding/json"
	// "fmt"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type TopologyMessageBody struct {
	Type string `json:"type"`
	Topology map[string][]string `json:"topology"`
}

type Server struct {
	n *maelstrom.Node
	messages map[int]struct{}
	neighbours []string
}

func main() {
	n := maelstrom.NewNode()
	var broadcast_vals []any
	var topology map[string][]string

	handle_broadcast := func (msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		broadcast_vals = append(broadcast_vals, body["message"])
		delete(body, "msg_id")

		for _, neighb := range topology[n.ID()] {
			neighb := neighb
			go func() {
				if err := n.Send(neighb, body); err != nil {
					panic(err)
				}
			}()
		}

		if msg.Src[0] != 'n' {
			return_body := make(map[string]any)
			return_body["type"] = "broadcast_ok"
			return n.Reply(msg, return_body)
		}
		return nil
	}

	handle_read := func (msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}	
		body["messages"] = broadcast_vals
		body["type"] = "read_ok"	
		return n.Reply(msg, body)
	}

	handle_topology := func (msg maelstrom.Message) error {
		var body TopologyMessageBody
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		topology = body.Topology

		return_body := make(map[string]any)
		return_body["type"] = "topology_ok"
		return n.Reply(msg, return_body)
	}

	n.Handle("broadcast", handle_broadcast)
	n.Handle("read", handle_read)
	n.Handle("topology", handle_topology)

	if err := n.Run(); n != nil {
		log.Fatal(err)
	}
}