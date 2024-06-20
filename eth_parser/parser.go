package ethparser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type ethParser struct {
	rpcURL          string
	lastSyncedBlock int
	transactions    map[string][]Transaction
	lock            sync.RWMutex
}

func NewEthParser(rpcURL string) Parser {
	parser := &ethParser{
		rpcURL:       rpcURL,
		transactions: make(map[string][]Transaction),
	}
	parser.lastSyncedBlock = parser.GetCurrentBlock()
	// sync every 5 seconds
	go func() {
		tick := time.NewTicker(5 * time.Second)
		for range tick.C {
			parser.sync()
		}
	}()
	return parser
}

func (p *ethParser) rpcRequest(method string, params []any, response any) error {
	// rpc request implementation
	requestBody, err := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      1,
	})
	if err != nil {
		return err
	}
	resp, err := http.Post(p.rpcURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		return err
	}
	return nil
}

func (p *ethParser) GetCurrentBlock() int {
	var result RpcResp[string]
	err := p.rpcRequest("eth_blockNumber", []any{"latest"}, &result)
	if err != nil {
		return -1
	}
	return HexToDec(result.Result)
}

func (p *ethParser) Subscribe(address string) bool {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.transactions[address] == nil {
		p.transactions[address] = make([]Transaction, 0)
	}
	return true
}

func (p *ethParser) GetTransactions(address string) []Transaction {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.transactions[address]
}

func (p *ethParser) fetchTransactions(blockNumber int) error {
	var result RpcResp[Block]
	err := p.rpcRequest("eth_getBlockByNumber", []any{DecToHex(blockNumber), true}, &result)
	if err != nil {
		return fmt.Errorf("failed to fetch block %d: %v", blockNumber, err)
	}

	p.lock.Lock()
	defer p.lock.Unlock()
	for _, tx := range result.Result.Transactions {
		// only store transactions from/to subscribed addresses
		if _, ok := p.transactions[tx.From]; ok {
			p.transactions[tx.From] = append(p.transactions[tx.From], tx)
		}
		if _, ok := p.transactions[tx.To]; ok {
			p.transactions[tx.To] = append(p.transactions[tx.To], tx)
		}
	}
	return nil
}

func (p *ethParser) sync() {
	currentBlock := p.GetCurrentBlock()
	if currentBlock == -1 {
		return
	}
	log.Printf("Syncing from block %d to %d\n", p.lastSyncedBlock+1, currentBlock)
	for i := p.lastSyncedBlock + 1; i <= currentBlock; i++ {
		err := p.fetchTransactions(i)
		if err != nil {
			log.Println(err)
		}
	}
	p.lastSyncedBlock = currentBlock
}
