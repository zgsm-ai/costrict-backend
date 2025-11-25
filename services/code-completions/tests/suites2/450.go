package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "math/big"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type Blockchain struct {
    blocks     []*Block
    mutex      sync.Mutex
    difficulty int
}

type Block struct {
    Index        int       `json:"index"`
    Timestamp    time.Time `json:"timestamp"`
    Data         string    `json:"data"`
    PreviousHash string    `json:"previous_hash"`
    Hash         string    `json:"hash"`
    Nonce        int       `json:"nonce"`
}

type Transaction struct {
    From   string `json:"from"`
    To     string `json:"to"`
    Amount int    `json:"amount"`
}

func NewBlockchain(difficulty int) *Blockchain {
    genesisBlock := createGenesisBlock()
    return &Blockchain{
        blocks:     []*Block{genesisBlock},
        difficulty: difficulty,
    }
}

func createGenesisBlock() *Block {
    return &Block{
        Index:        0,
        Timestamp:    time.Now(),
        Data:         "Genesis Block",
        PreviousHash: "0",
        Hash:         "",
        Nonce:        0,
    }
}

func (bc *Blockchain) GetLatestBlock() *Block {
    bc.mutex.Lock()
    defer bc.mutex.Unlock()
    
    return bc.blocks[len(bc.blocks)-1]
}

func (bc *Blockchain) AddBlock(data string) error {
    previousBlock := bc.GetLatestBlock()
    
    newBlock := &Block{
        Index:        previousBlock.Index + 1,
        Timestamp:    time.Now(),
        Data:         data,
        PreviousHash: previousBlock.Hash,
        Hash:         "",
        Nonce:        0,
    }
    
    newBlock.Hash = bc.calculateHash(newBlock)
    newBlock = bc.proofOfWork(newBlock)
    
    bc.mutex.Lock()
    bc.blocks = append(bc.blocks, newBlock)
    bc.mutex.Unlock()
    
    return nil
}

func (bc *Blockchain) calculateHash(block *Block) string {
    record := fmt.Sprintf("%d%s%s%s%d",
        block.Index,
        block.Timestamp.String(),
        block.Data,
        block.PreviousHash,
        block.Nonce,
    )
    
    // Simple hash function for demonstration
    // In a real blockchain, you would use a cryptographic hash like SHA-256
    var hash big.Int
    hash.SetString(record, 10)
    
    return hash.Text(16)
}

func (bc *Blockchain) proofOfWork(block *Block) *Block {
    target := big.NewInt(1)
    target.Lsh(target, uint(256-bc.difficulty))
    
    var hashInt big.Int
    var hash string
    
    for {
        hash = bc.calculateHash(block)
        hashInt.SetString(hash, 16)
        
        if hashInt.Cmp(target) == -1 {
            break
        }
        
        block.Nonce++
    }
    
    return block
}

func (bc *Blockchain) IsValid() bool {
    bc.mutex.Lock()
    defer bc.mutex.Unlock()
    
    for i := 1; i < len(bc.blocks); i++ {
        currentBlock := bc.blocks[i]
        previousBlock := bc.blocks[i-1]
        
        if currentBlock.Hash != bc.calculateHash(currentBlock) {
            return false
        }
        
        if currentBlock.PreviousHash != previousBlock.Hash {
            return false
        }
    }
    
    return true
}

func (bc *Blockchain) SaveToFile(filename string) error {
    bc.mutex.Lock()
    defer bc.mutex.Unlock()
    
    data, err := json.MarshalIndent(bc.blocks, "", "  ")
    if err != nil {
        return err
    }
    
    // Create directory if it doesn't exist
    dir := filepath.Dir(filename)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }
    
    <｜fim▁hole｜>
    
    return nil
}

func (bc *Blockchain) LoadFromFile(filename string) error {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return err
    }
    
    var blocks []*Block
    if err := json.Unmarshal(data, &blocks); err != nil {
        return err
    }
    
    bc.mutex.Lock()
    bc.blocks = blocks
    bc.mutex.Unlock()
    
    return nil
}

func (bc *Blockchain) GetChainInfo() map[string]interface{} {
    bc.mutex.Lock()
    defer bc.mutex.Unlock()
    
    info := make(map[string]interface{})
    info["length"] = len(bc.blocks)
    info["difficulty"] = bc.difficulty
    info["valid"] = bc.IsValid()
    
    if len(bc.blocks) > 0 {
        info["latest_block"] = bc.blocks[len(bc.blocks)-1]
    }
    
    return info
}

func main() {
    // Create a new blockchain with difficulty 4
    blockchain := NewBlockchain(4)
    defer blockchain.SaveToFile("blockchain.json")
    
    fmt.Println("Blockchain Simulation")
    fmt.Println("==================")
    
    // Add some blocks
    transactions := []Transaction{
        {"Alice", "Bob", 100},
        {"Bob", "Charlie", 50},
        {"Charlie", "David", 25},
        {"David", "Alice", 10},
    }
    
    for i, tx := range transactions {
        txData, _ := json.Marshal(tx)
        err := blockchain.AddBlock(string(txData))
        if err != nil {
            fmt.Printf("Error adding block %d: %v\n", i+1, err)
            continue
        }
        
        fmt.Printf("Block %d added: %s -> %s: %d\n", 
            i+2, tx.From, tx.To, tx.Amount)
    }
    
    // Check if blockchain is valid
    fmt.Printf("\nBlockchain is valid: %t\n", blockchain.IsValid())
    
    // Display blockchain info
    info := blockchain.GetChainInfo()
    infoJSON, _ := json.MarshalIndent(info, "", "  ")
    fmt.Printf("\nBlockchain Info:\n%s\n", infoJSON)
    
    // Try to load from file
    newBlockchain := NewBlockchain(4)
    if err := newBlockchain.LoadFromFile("blockchain.json"); err != nil {
        fmt.Printf("Error loading blockchain: %v\n", err)
    } else {
        fmt.Println("\nBlockchain loaded from file successfully")
        newInfo := newBlockchain.GetChainInfo()
        newInfoJSON, _ := json.MarshalIndent(newInfo, "", "  ")
        fmt.Printf("\nLoaded Blockchain Info:\n%s\n", newInfoJSON)
    }
    
    // Simulate tampering with the blockchain
    if len(blockchain.blocks) > 1 {
        blockchain.mutex.Lock()
        blockchain.blocks[1].Data = "Tampered Data"
        blockchain.mutex.Unlock()
        
        fmt.Printf("\nBlockchain is valid after tampering: %t\n", blockchain.IsValid())
    }
}