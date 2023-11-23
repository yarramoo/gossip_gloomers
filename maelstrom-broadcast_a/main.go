package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()
	var broadcast_vals []any
	// var topology any

	handle_broadcast := func (msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		broadcast_vals = append(broadcast_vals, body["message"])

		return_body := make(map[string]any)
		return_body["type"] = "broadcast_ok"
		return n.Reply(msg, return_body)
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
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}	
		// topology = body["topology"]

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