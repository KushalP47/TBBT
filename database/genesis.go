package database

import (
	"encoding/json" // Package json implements encoding and decoding of JSON
	"os"            // Package os provides a platform-independent interface to operating system functionality
)

type genesis struct {
	Balances map[Account]uint `json:"balances"` // map of account to balance
}

// loadGenesis returns a genesis and an error, it will give all the details of the genesis block from the genesis.json file
func LoadGenesis(path string) (genesis, error) {
	content, err := os.ReadFile(path) // ReadFile reads the file named by filename and returns the contents

	if err != nil { // if error is not nil
		return genesis{}, err
	}

	var loadedGenesis genesis                     // declare a variable loadedGenesis of type genesis
	err = json.Unmarshal(content, &loadedGenesis) // Unmarshal parses the JSON-encoded data and stores the decoded result in the given variable
	if err != nil {
		return genesis{}, err
	}

	return loadedGenesis, nil // return loadedGenesis and nil
}
