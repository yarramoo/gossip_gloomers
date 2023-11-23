package main

import (
	"encoding/json"
	"log"
	"math/rand"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main_() {
	n := maelstrom.NewNode()
	var id int64
	
	n.Handle("generate", func(msg maelstrom.Message) error {
		msg_id := int64(msg.Src[1] - '0')
		if msg_id != id {
			id = msg_id
			rand.Seed(id)
		}
		
		
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		body["type"] = "generate_ok"
		body["id"] = rand.Int()

		return n.Reply(msg, body)
	})

	if err := n.Run(); n != nil {
		log.Fatal(err)
	}
}