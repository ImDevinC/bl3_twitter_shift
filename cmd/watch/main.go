package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/imdevinc/bl3_twitter_shift/internal/htmlupdater"

	"github.com/imdevinc/bl3_twitter_shift/internal/monitor"
)

func main() {
	client := monitor.NewTwitterClient(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"), os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_SECRET"))
	client.StartMonitor(os.Getenv("TWITTER_USER"), func(keys []string, timestamp string) {
		htmlupdater.AddKeys(os.Getenv("ORCZ_TITLE"), os.Getenv("ORCZ_USERNAME"), os.Getenv("ORCZ_PASSWORD"), keys, timestamp)
	})
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
	client.StopMonitor()
}
