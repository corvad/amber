package main

import (
	"encoding/json"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

var (
	messages      map[int]int
	messagesMutex *sync.Mutex
	n             *maelstrom.Node
)

func broadcast(msg maelstrom.Message) error {
	// Unmarshal the message body as an loosely-typed map.
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	i := body["message"].(float64)
	messagesMutex.Lock()
	messages[int(i)]++
	messagesMutex.Unlock()
	// Update the message type to return back.
	body["type"] = "broadcast_ok"
	delete(body, "message")

	// Echo the original message back with the updated message type.
	return n.Reply(msg, body)
}

func read(msg maelstrom.Message) error {
	// Unmarshal the message body as an loosely-typed map.
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	messagesMutex.Lock()
	values := make([]int, 0, len(messages))
	for k, _ := range messages {
		values = append(values, k)
	}
	messagesMutex.Unlock()

	// Update the message type to return back.
	body["type"] = "read_ok"

	body["messages"] = values

	// Echo the original message back with the updated message type.
	return n.Reply(msg, body)
}

func topology(msg maelstrom.Message) error {
	// Unmarshal the message body as an loosely-typed map.
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	// Update the message type to return back.
	body["type"] = "topology_ok"
	delete(body, "topology")
	// Echo the original message back with the updated message type.
	return n.Reply(msg, body)
}

func main() {
	messagesMutex = &sync.Mutex{}
	messages = make(map[int]int)
	n = maelstrom.NewNode()
	n.Handle("broadcast", broadcast)
	n.Handle("read", read)
	n.Handle("topology", topology)
	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
