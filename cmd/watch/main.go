package main

import (
	"github.com/imdevinc/bl3_twitter_shift/internal/monitor"
	"os"
)

func main() {
	client := monitor.NewTwitterClient(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"))
	client.StartMonitor()
}
