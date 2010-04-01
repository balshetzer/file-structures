package main;

import . "block/byteslice"
import "bptree"
import "os"
import "bufio"
import "fmt"
import "json"
import "log"

type Metadata struct {
    Filename string
    Keysize uint32
    Fieldsizes []uint32
}

type Command struct {
    Op string
    LeftKey ByteSlice
    RightKey ByteSlice
    Fields []ByteSlice
}

func main() {
    // Read the string
    inputReader := bufio.NewReader(os.Stdin)
    
    var info = Metadata{"", uint32(0), nil}
    var cmd = Command{"", nil, nil, nil}
    
    infoJson, err := inputReader.ReadString('\n')
    if err != nil {
        log.Exit(err)
    } else {
        json.Unmarshal(infoJson, &info)
    }
    
    bptree, bperr := bptree.NewBpTree(info.Filename, info.Keysize, info.Fieldsizes)
    if bperr {
        log.Exit("Failed B+ tree creation")
    }
    
    alive := true;
    
    for alive {
        cmdJson, err := inputReader.ReadString('\n')
        if err != nil {
            log.Exit(err)
        }
        if cmdJson == "q\n" {
            alive = false
        } else {
            json.Unmarshal(cmdJson, &cmd)
            if cmd.Op == "insert" {
                result := bptree.Insert(cmd.LeftKey, cmd.Fields)
                fmt.Println(result)
            } else if cmd.Op == "find" {
                records, ack := bptree.Find(cmd.LeftKey, cmd.RightKey)
                for record := range records {
                    fmt.Println(record)
                    ack<-true;                              // ack<-true must be the last line of the loop.
                }
                fmt.Println("end")
            }
        }
    }
    fmt.Println("exited")
}

// Determine which file and schema is being opened
//  (filename string, keysize uint32, fields []uint32)

// insert(key Byteslice, record []Byteslice)
// find(leftkey, right key) returns channel with all matching keys+records (Record structs)
