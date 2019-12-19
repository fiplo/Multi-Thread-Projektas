package main

import (
  "github.com/jessevdk/go-flags"
  "github.com/wcharczuk/go-chart"
  //"bytes"
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

type record struct {
    threadCount int
    dataSize int
    time int64 
}

type chartData struct{
    threads []int
    time []int64
    dataSize int
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

  //optTCount := opts.ThreadCount
  optOFile := opts.Output
  optIFile := opts.Input

  if optOFile == "" {
    optOFile = "results"
  }

  if optIFile == "" {
    optIFile = RandStringBytes(1000)
  }


  for i := 1; i < 11; i++ {
    var result record
    oppSize := (10^i)
    var thrdarr []float64
    var timearr []float64
    filewrite := fmt.Sprint("output",oppSize,".png")
    for y := 0; y < 40; y++ {
      result = Hashing(optIFile, y+1, oppSize)
      thrdarr = append(thrdarr, float64(result.threadCount))
      timearr = append(timearr, float64(result.time))
    }
    graph := chart.Chart{
        Series: []chart.Series{
          chart.ContinuousSeries{
            XValues: thrdarr,
            YValues: timearr,
          },
        },
    }
    f,_ := os.Create(filewrite)
    defer f.Close()
    graph.Render(chart.PNG, f)
  } 
}


func Hashing(inputData string, threadCount int, oppSize int) (record) {
    ch := make(chan bool)
    var wg sync.WaitGroup

    start := time.Now();
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
    elapsed := time.Since(start)
    end := record{threadCount, oppSize, elapsed.Microseconds()}

    return end
}
