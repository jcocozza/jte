package main

import "github.com/jcocozza/jte/internal/logger"

func main() {
    l := logger.NewLogger()
    l.Debug("hello")
}
