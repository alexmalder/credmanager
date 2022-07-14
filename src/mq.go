package src

import (
	"context"
	"fmt"
	"log"
	"time"

	zmq "github.com/go-zeromq/zmq4"
	"github.com/vmihailenco/msgpack/v5"
)

type Item struct {
	Elem int
	Foo  string
}

func ZMQServer() error {
	ctx := context.Background()
	// Socket to talk to clients
	socket := zmq.NewRep(ctx)
	defer socket.Close()
	if err := socket.Listen("tcp://*:5555"); err != nil {
		return fmt.Errorf("listening: %w", err)
	}
	for {
		msg, err := socket.Recv()
		if err != nil {
			return fmt.Errorf("receiving: %w", err)
		}
		var item Item
		err = msgpack.Unmarshal(msg.Bytes(), &item)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(item.Foo, item.Elem)

		//fmt.Println("Received ", msg)
		// Do some 'work'
		//reply := fmt.Sprintf("%s", item.Foo)
		if err := socket.Send(zmq.NewMsgString(item.Foo)); err != nil {
			return fmt.Errorf("sending reply: %w", err)
		}
	}
}

func ZMQClient(message string) error {
	ctx := context.Background()
	socket := zmq.NewReq(ctx, zmq.WithDialerRetry(time.Second))
	defer socket.Close()

	fmt.Printf("Connecting to hello world server...")
	if err := socket.Dial("tcp://localhost:5555"); err != nil {
		return fmt.Errorf("dialing: %w", err)
	}
    enc, err := EncTest(message)
    if err != nil {
        log.Fatal(err)
    }
	b, err := msgpack.Marshal(&Item{Foo: enc, Elem: 4})
	if err != nil {
		log.Fatal(err)
	}
	m := zmq.NewMsgFrom(b)
	fmt.Println("sending ", m)
	if err := socket.Send(m); err != nil {
		return fmt.Errorf("sending: %w", err)
	}
	// Wait for reply.
	r, err := socket.Recv()
	if err != nil {
		return fmt.Errorf("receiving: %w", err)
	}
    msg := fmt.Sprintf("%s", r.Frames[0])
    fmt.Println("Message ", msg)
    resp, err := DecTest(msg)
    fmt.Println("---")
    fmt.Println(resp)
	return nil
}
