package main

import (
	"os"

	"github.com/imdevinc/bl3_twitter_shift/internal/monitor"
)

func main() {
	client := monitor.NewTwitterClient(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"), os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_SECRET"))
	client.StartMonitor(os.Getenv("TWITTER_USER"))
}
