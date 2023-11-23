package main

import (
	"encoding/json"
	"strconv"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()
	counter := 0

	handle := func (msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		body["type"] = "generate_ok"

		// Calculate the new unique id
		counter_str := strconv.Itoa(counter)
		body["id"] = msg.Src + counter_str
		counter += 1

		return n.Reply(msg, body)
	}
	
	n.Handle("generate", handle)

	if err := n.Run(); n != nil {
		log.Fatal(err)
	}
}