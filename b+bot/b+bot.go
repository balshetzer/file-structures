package main;

import . "file-structures/block/byteslice"
import "file-structures/bptree"
import "os"
import "bufio"
import "fmt"
import "encoding/json"

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
    FileName string
}

func main() {
    // Read the string
    inputReader := bufio.NewReader(os.Stdin)

    var info = Metadata{"", uint32(0), nil}
    var cmd = Command{"", nil, nil, nil, ""}

    infoJson, err := inputReader.ReadBytes('\n')
    if err != nil {
        panic(err)
    } else {
        json.Unmarshal(infoJson, &info)
    }

    bpt, bperr := bptree.NewBpTree(info.Filename, info.Keysize, info.Fieldsizes)
    if !bperr {
        panic("Failed B+ tree creation")
    } else {
        fmt.Println("ok")
    }

    alive := true;

    for alive {
        cmdJson, err := inputReader.ReadBytes('\n')
        if err != nil {
            alive = false
            break
        }
        if cmdJson[0] == 'q' && cmdJson[1] == '\n' {
            alive = false
        } else {
            json.Unmarshal(cmdJson, &cmd)
            if cmd.Op == "insert" {
                result := bpt.Insert(cmd.LeftKey, cmd.Fields)
                fmt.Println(result)
            } else if cmd.Op == "find" {
                records, ack := bpt.Find(cmd.LeftKey, cmd.RightKey)
                for record := range records {
                    if bytes, err := json.Marshal(map[string]interface{}{
                      "key": record.GetKey(),
                      "value": record.AllFields()}); err != nil {
                        panic(err)
                    } else {
                      os.Stdout.Write(bytes)
                    }
                    fmt.Println()
                    ack<-true
                }
                fmt.Println("end")
            } else if cmd.Op == "visualize" {
                bptree.Dotty(cmd.FileName, bpt)
            } else if cmd.Op == "prettyprint" {
                s := fmt.Sprintln(bpt)
                f, _ := os.Create(cmd.FileName)
                f.Write([]byte(s))
                f.Close()
            }
        }
    }
    fmt.Println("exited")
}

// Determine which file and schema is being opened
//  (filename string, keysize uint32, fields []uint32)

// insert(key Byteslice, record []Byteslice)
// find(leftkey, right key) returns channel with all matching keys+records (Record structs)

