package server

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

type message struct {
	id      uint32
	payload []byte
}

type server struct {
	listener net.Listener
}

func NewServer(port string) (*server, error) {
	l, err := net.Listen("tcp", port)
	if err != nil {
		return nil, err
	}
	return &server{listener: l}, nil
}

func (s *server) Start() error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *server) stop() error {
	return s.listener.Close()
}

func (s *server) handleConnection(conn net.Conn) {
	defer conn.Close()

	incoming := make(chan message)
	outgoing := make(chan []byte, 10)

	// Start a goroutine to read incoming messages from the client
	go s.readIncomingMessages(conn, incoming, outgoing)

	// Start a goroutine to write outgoing messages to the client
	go s.writeOutgoingMessages(conn, outgoing)

	// Start processing incoming messages from the client
	s.processIncomingMessages(conn, incoming, outgoing)
}

func (s *server) readIncomingMessages(conn net.Conn, incoming chan message, outgoing chan []byte) {
	defer close(incoming)
	defer close(outgoing)

	for {
		header := make([]byte, 8)
		_, err := io.ReadFull(conn, header)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed by client")
			} else {
				fmt.Println("Error reading header:", err.Error())
			}
			return
		}
		packetID := binary.BigEndian.Uint32(header[:4])
		payloadLength := binary.BigEndian.Uint32(header[4:])
		payload := make([]byte, payloadLength)
		_, err = io.ReadFull(conn, payload)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed by client")
			} else {
				fmt.Println("Error reading payload:", err.Error())
			}
			return
		}
		incoming <- message{id: packetID, payload: payload}
	}
}

func (s *server) writeOutgoingMessages(conn net.Conn, outgoing chan []byte) {
	for data := range outgoing {
		_, err := conn.Write(data)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed by client")
			} else {
				fmt.Println("Error writing data:", err.Error())
			}
			return
		}
	}
}

func (s *server) processIncomingMessages(conn net.Conn, incoming chan message, outgoing chan []byte) {
	for msg := range incoming {
		switch msg.id {
		case 1:
			payloadStr := string(msg.payload)
			fmt.Printf("Packet received (ID %d, payload length %d): %s\n", msg.id, len(msg.payload), payloadStr)
			// Send a response
			response := []byte("Hello, client!")
			outgoing <- append([]byte{0, 0, 0, 1}, uint32ToBytes(uint32(len(response)))...)
			outgoing <- response
		default:
			fmt.Printf("Unknown packet ID: %d\n", msg.id)
		}
	}
}

func uint32ToBytes(n uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, n)
	return b
}
