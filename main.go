package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

// import (
// 	"fmt"
// 	"log"
// 	"regexp"
// 	"strconv"
// )

// type TreeNode struct {
// 	Val   int
// 	Left  *TreeNode
// 	Right *TreeNode
// }

// func recoverFromPreorder(traversal string) *TreeNode {
// 	values, err := parse(traversal)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var root *TreeNode

// 	for _, v := range values {
// 		if root == nil {
// 			root = newNode(v)
// 		} else {
// 			insert(root, v)
// 		}
// 	}

// 	return root
// }

// func newNode(val int) *TreeNode {
// 	return &TreeNode{Val: val}
// }

// func insert(root *TreeNode, v int) {
// 	if v >= root.Val {
// 		if root.Left != nil {
// 			insert(root.Left, v)
// 		} else {
// 			root.Left = newNode(v)
// 		}
// 	} else {
// 		if root.Right != nil {
// 			insert(root.Right, v)
// 		} else {
// 			root.Right = newNode(v)
// 		}
// 	}
// }

// func parse(input string) ([]int, error) {
// 	reg := regexp.MustCompile("[-]+")
// 	rawValues := reg.Split(input, -1)

// 	values := make([]int, 0, len(rawValues))

// 	for _, raw := range rawValues {
// 		if raw == "" {
// 			continue
// 		}
// 		v, err := strconv.Atoi(raw)
// 		if err != nil {
// 			return nil, err
// 		}
// 		values = append(values, v)
// 	}

// 	return values, nil
// }

// func main() {
// 	root := recoverFromPreorder("1-401--349---90--88")
// 	bfs(root)
// }

// func bfs(root *TreeNode) {
// 	if root == nil {
// 		fmt.Print("null ")
// 		return
// 	}

// 	fmt.Printf("%v ", root.Val)

// 	bfs(root.Left)
// 	bfs(root.Right)
// }

func read(ctx context.Context, n int) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:29092"},
		Topic:   "messages",
		// GroupID: "clients",
	})
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("%v: %v\n", n, string(msg.Value))
	}
}

func main() {
	ctx := context.Background()

	var wg sync.WaitGroup
	wg.Add(2)

	go read(ctx, 1)
	go read(ctx, 2)
	go read(ctx, 3)
	go read(ctx, 4)

	go func() {
		producer := kafka.NewWriter(kafka.WriterConfig{
			Brokers:  []string{"localhost:29092"},
			Topic:    "messages",
			Balancer: &kafka.RoundRobin{},
		})
		defer producer.Close()

		for {
			time.Sleep(2 * time.Second)
			msg := kafka.Message{
				Value: []byte("Hell"),
			}
			if err := producer.WriteMessages(ctx, msg); err != nil {
				log.Fatal(err)
			}
		}
	}()

	wg.Wait()

	log.Println("done")
}
