package database

import (
	"bufio"         // Package bufio implements buffered I/O
	"encoding/json" // Package json implements encoding and decoding of JSON
	"fmt"           // Package fmt implements formatted I/O with functions analogous to C's printf and scanf
	"os"            // Package os provides a platform-independent interface to operating system functionality
	"path/filepath" // Package filepath implements utility routines for manipulating filename paths
)

type State struct { // declare a struct type State
	Balances  map[Account]uint // map of account to balance
	txMempool []Tx             // array of Transactions

	dbFile *os.File // pointer to tx.db file, where we will store the transactions
}

func NewStateFromDisk() (*State, error) { // declare a function NewStateFromDisk which returns a pointer to State and an error
	cwd, err := os.Getwd() // Getwd returns a rooted path name corresponding to the current directory

	if err != nil {
		return nil, err
	}

	gen, err := loadGenesis(filepath.Join(cwd, "database", "genesis.json"))
	// loadGenesis returns a genesis details(genesis.json) and an error
	if err != nil {
		return nil, err
	}

	balances := make(map[Account]uint) // declaring a map of account to track balances

	for account, balance := range gen.Balances {
		// copying all the balances from genesis to the balances map
		balances[account] = balance
	}

	f, err := os.OpenFile(filepath.Join(cwd, "database", "tx.db"), os.O_APPEND|os.O_RDWR, 0600)
	// f points to the tx.db file so that we can access all the transactions happened till now
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)
	// NewScanner returns a new Scanner to read from the tx.db

	state := &State{balances, make([]Tx, 0), f}
	// declaring a pointer to State and initializing it with balances, empty txMempool and f
	// we will be returning this state after filling it with data

	for scanner.Scan() { // scan the tx.db file
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		var tx Tx                                   // a dummy tx variable to store the one transaction at a time from tx.db file
		err := json.Unmarshal(scanner.Bytes(), &tx) // Unmarshal parses the JSON-encoded data and stores the decoded result in tx
		if err != nil {
			return nil, err
		}

		if err := state.apply(tx); err != nil { // calls the aplly function, which will verify the transaction and append it to the state
			return nil, err
		}
	}

	return state, nil
}

func (s *State) Add(tx Tx) error {

	// verify the transaction
	if err := s.apply(tx); err != nil {
		return err
	}

	// append the transaction to the txMempool
	s.txMempool = append(s.txMempool, tx)

	return nil
}

func (s *State) Persist() error {
	// it will write all the transactions from txMempool to the tx.db file and clear the txMempool
	// first it will make a copy of the txMempool of a particular instance of State, so that the original txMempool is not affected
	mempool := make([]Tx, len(s.txMempool))
	copy(mempool, s.txMempool)

	for i := 0; i < len(mempool); i++ {
		// for each transaction in the mempool, it will write the transaction to the tx.db file
		txJson, err := json.Marshal(mempool[i])
		// txMempool stores all transaction in tx format while tx.db stores it in JSON format hence we have to convert it to JSON format
		if err != nil {
			return err
		}

		// appending the transaction to the tx.db file
		if _, err = s.dbFile.Write(append(txJson, '\n')); err != nil {
			return err
		}

		// removing the transaction from the txMempool from front
		s.txMempool = s.txMempool[1:]
	}

	return nil
}

func (s *State) Close() {
	// close the tx.db file
	s.dbFile.Close()
}

func (s *State) apply(tx Tx) error {
	// it will verify the transaction and append it to the state

	if tx.IsReward() {
		// if it is a request for reward, then it will add the reward to the account
		// this type of transactions are done by the miners, to reward them as we are not deducing any amount from any account, we will only append the amt to the account
		s.Balances[tx.To] += tx.Value
		return nil
	}

	if s.Balances[tx.From] < tx.Value {
		// if the balance of the account is less than the value of the transaction, then it will return an error
		return fmt.Errorf("insufficient balance")
	}

	// if the balance is sufficient, then it will deduct the value from the sender's account and add it to the receiver's account
	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value

	return nil
}
