package ethparser

import "strings"

func ValidateEthAddr(arg string) bool {
	return strings.HasPrefix(arg, "0x") && len(arg) == 42
}

type RpcResp[Result any] struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  Result `json:"result"`
	Id      int    `json:"id"`
}

type Block struct {
	Number string `json:"number"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	From        string `json:"from"`
	To          string `json:"to"`
	Value       string     `json:"value"`
	BlockNumber string     `json:"blockNumber"`
	Hash        string     `json:"hash"`
}

type Parser interface {
	GetCurrentBlock() int
	Subscribe(address string) bool
	GetTransactions(address string) []Transaction
}
