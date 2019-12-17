package main

import (
  "github.com/jessevdk/go-flags"
  "os"
  "crypto/sha256"
  "math/rand"
  "time"
  "fmt"
  "sync"
)


func init() {
    rand.Seed(time.Now().UnixNano())
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}

func main() {
  opts := &Options{}
  parser := flags.NewParser(opts, flags.Default)
  parser.Name = "Hash test"
  parser.Usage = "-i[--input] <input> -o[--output] <output> -c[--count] <the-number-of-threads>"
  _, err := parser.Parse()
  if err != nil {
    os.Exit(1)
  }

  optTCount := opts.ThreadCount
  optOFile := opts.Output
  optIFile := opts.Input

  if optOFile == "" {
    optOFile = "results"
  }

  if optIFile == "" {
    optIFile = RandStringBytes(1000)
  }

  result, err := Hashing(optIFile, optTCount, 1000)
  fmt.Print(result)
}


func Hashing(inputData string, threadCount int, oppSize int) (string, error) {
    ch := make(chan bool)
    var wg sync.WaitGroup

    for i := 0; i < threadCount; i++ {
      wg.Add(1)
      go func(input string, ch chan bool) {
        for {
          _, ok := <-ch
          if !ok {
            wg.Done()
            return
          }
          h := sha256.New()
          h.Write([]byte(input))
          h.Sum(nil)
        }
      }(inputData, ch)
    }
    for i := 0; i < oppSize; i++ {
      ch <- false
}
    close(ch)
    wg.Wait()

    return "success", nil
}
