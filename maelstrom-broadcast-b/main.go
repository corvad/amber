package main

import (
	"encoding/json"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

var (
	messages       map[int]int
	messagesMutex  *sync.Mutex
	neighbors      []string
	neighborsMutex *sync.Mutex
	n              *maelstrom.Node
)

func broadcast(msg maelstrom.Message) error {
	// Unmarshal the message body as an loosely-typed map.
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	i := body["message"].(float64)
	messagesMutex.Lock()
	_, repeat := messages[int(i)]
	if !repeat {
		messages[int(i)]++
		messagesMutex.Unlock()
		neighborsMutex.Lock()
		for _, node := range neighbors {
			if err := n.Send(node, msg.Body); err != nil {
				log.Fatal(err)
			}
		}
		neighborsMutex.Unlock()
	} else {
		messagesMutex.Unlock()
	}
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
	type Topology struct {
		Type     string              `json:"type"`
		Topology map[string][]string `json:"topology"`
	}
	var body Topology
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	neighborsMutex.Lock()
	nodes := body.Topology
	neighbors = nodes[n.ID()]
	neighborsMutex.Unlock()

	type Reply struct {
		Type string `json:"type"`
	}

	return n.Reply(msg, Reply{"topology_ok"})
}

func main() {
	messagesMutex = &sync.Mutex{}
	messages = make(map[int]int)
	neighborsMutex = &sync.Mutex{}
	neighbors = make([]string, 0, 10)
	n = maelstrom.NewNode()
	n.Handle("broadcast", broadcast)
	n.Handle("read", read)
	n.Handle("topology", topology)
	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
