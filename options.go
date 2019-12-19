package main

type Options struct {
  ThreadCount int     `short:"c" long:"tcount" description:"Thread count for operation"`
  Input     string  `short:"i" long:"input" description:"Input data that will be hashed"`
  Output      string  `short:"o" long:"output" description:"Output of given operation status, by default outputs to terminal"`
}
