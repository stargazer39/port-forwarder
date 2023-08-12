package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
)

type IncomingProxyConnection struct {
	ID         string
	Type       int // Type 0 is initial control connection, Type 1 is proxy connection
	Connection net.Conn
	Locked     bool
}

type IncomingProxyMessage struct {
	ID   string
	Type int
}

type MessageToClient struct {
	Message string
}

func ReadJSON(reader io.Reader, i any) error {
	var length uint64
	lengthBuffer := make([]byte, 8)

	_, err := io.ReadFull(reader, lengthBuffer)

	if err != nil {
		return err
	}

	length = binary.BigEndian.Uint64(lengthBuffer)

	messageBuffer := make([]byte, length)

	nRead, err2 := io.ReadFull(reader, messageBuffer)

	if err2 != nil {
		return err
	}

	if nRead != len(messageBuffer) {
		return fmt.Errorf("Couldn't read the whole message.")
	}

	return json.Unmarshal(messageBuffer, i)
}

func WriteJSON(writer io.Writer, i any) error {
	b, err := json.Marshal(i)
	var length uint64
	var messageBuffer []byte

	if err != nil {
		return err
	}

	length = uint64(len(b))
	lengthBuffer := make([]byte, 8)

	binary.BigEndian.PutUint64(lengthBuffer, length)

	messageBuffer = append(messageBuffer, lengthBuffer...)
	messageBuffer = append(messageBuffer, b...)

	nWritten, err := writer.Write(messageBuffer)

	if err != nil {
		return err
	}

	if nWritten != len(messageBuffer) {
		return fmt.Errorf("Couldn't write the whole message.")
	}

	return nil
}

func HandleTCPConn(src net.Conn, dst net.Conn) {
	done := make(chan struct{})

	go func() {
		defer src.Close()
		defer dst.Close()
		io.Copy(dst, src)
		done <- struct{}{}
	}()

	go func() {
		defer src.Close()
		defer dst.Close()
		io.Copy(src, dst)
		done <- struct{}{}
	}()

	<-done
	<-done
}
