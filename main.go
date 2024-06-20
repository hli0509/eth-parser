package main

import (
	"log"

	ethparser "eth-parser/eth_parser"
	server "eth-parser/server"
)

func main() {
	parser := ethparser.NewEthParser("https://cloudflare-eth.com")

	srv := server.NewServer(parser)
	log.Fatal(srv.Start())
}
