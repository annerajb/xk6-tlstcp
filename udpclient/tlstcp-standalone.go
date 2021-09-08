package main

// import "fmt"

// func main() {
// 	fmt.Println("Hello, World!")
// }

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		log.Fatal("Please provide a host:port string")
	}
	//load public key
	content, err := ioutil.ReadFile("ntp.pubkey")
	if err != nil {
		log.Fatalf("pubkey load failure: %s", err)
	}

	CONNECT := arguments[1]

	s, err := net.ResolveUDPAddr("udp4", CONNECT)
	if err != nil {
		log.Fatalf("dns resolution fail: %s", err.Error())
	}
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Sending to UDP server at: %s\n", c.RemoteAddr().String())
	defer c.Close()

	//all except the challenge
	b64originalpacket := "GwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAARwAAJA=="
	decodeNtp, _ := base64.StdEncoding.DecodeString(b64originalpacket)
	fmt.Printf("request beginning len: %d\n", len(decodeNtp))

	b53 := "GwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAARwAAJAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	decodedString, err := base64.StdEncoding.DecodeString(b53)
	if err != nil {
		log.Fatalf("Error Found: %s\n", err)
	}

	fmt.Printf("fullworking beginning len: %d\n", len(decodedString))

	ntp_request := make([]byte, 84)

	key := [32]byte{}

	_, err = rand.Read(key[:])
	if err != nil {
		log.Fatalf("Error Found: %s\n", err)
	}

	fmt.Printf("challenge: %s\n", strings.ToUpper(hex.EncodeToString(key[:8])))
	fmt.Printf("           %s\n", strings.ToUpper(hex.EncodeToString(key[8:16])))
	fmt.Printf("           %s\n", strings.ToUpper(hex.EncodeToString(key[16:24])))
	fmt.Printf("           %s\n", strings.ToUpper(hex.EncodeToString(key[24:32])))

	copy(ntp_request, decodeNtp[:])
	copy(ntp_request[52:], key[:])
	fmt.Printf("sending size: %d\n", len(ntp_request))
	//sends
	for i := 0; i < 7; i++ {
		_, err = c.Write(ntp_request)

		if err != nil {
			log.Fatalf("write failure: %s\n", err.Error())
		}
	}
	//receives
	buffer := make([]byte, 1024)

	validResponses := 0
	for validResponses = 0; validResponses < 7; {
		n, _, err := c.ReadFromUDP(buffer)
		if err != nil {
			log.Fatalf("read failure: %s", err.Error())
		}
		fmt.Printf("Received Challenge size: %d\n", n)
		if n != 148 {
			fmt.Printf("received packet mismatched: %d != %d\n", n, 148)
			break
		}
		ntp_response_bytes := buffer[:148]
		signature := ntp_response_bytes[0:64]
		fmt.Printf("signature: %s\n", hex.EncodeToString(signature))
		signedPacket := ntp_response_bytes[64:]
		fmt.Printf("signature: %s\n", hex.EncodeToString(signedPacket))
		// response := ntp_response_bytes[64:]
		// fmt.Printf("Signed: %s\n", hex.EncodeToString(response))

		if !ed25519.Verify(content, signedPacket, signature) {
			fmt.Print("signature verify failed\n")
			break
		}

		validResponses++
	}
	//verify / parse
	for i := 0; i < validResponses; i++ {
		//crypto verify

	}
}
