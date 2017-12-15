package main

import (
	"errors"
	"fmt"
	"sync"
	"time"

	tm "github.com/buger/goterm"
	"github.com/spf13/cobra"
)

var (
	COUNT            = 5
	TIMEOUT          = 800 * time.Millisecond
	INTERVAL         = 100 * time.Millisecond
	MAX_HOPS         = 64
	RING_BUFFER_SIZE = 50
)

// rootCmd represents the root command
var RootCmd = &cobra.Command{
	Use: "mtr TARGET",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("No target provided")
		}
		fmt.Println("Start:", time.Now())
		m, ch := NewMTR(args[0], TIMEOUT, INTERVAL)
		tm.Clear()
		mu := &sync.Mutex{}
		go func(ch chan struct{}) {
			for {
				mu.Lock()
				<-ch
				render(m)
				mu.Unlock()
			}
		}(ch)
		for i := 0; i < COUNT; i++ {
			m.Run(ch)
		}
		close(ch)
		mu.Lock()
		render(m)
		mu.Unlock()
		return nil
	},
}

func render(m *MTR) {
	tm.MoveCursor(1, 1)
	m.Render(1)
	tm.Flush() // Call it every time at the end of rendering
}

func init() {
	RootCmd.Flags().IntVarP(&COUNT, "count", "c", COUNT, "Amount of pings per target")
	RootCmd.Flags().DurationVarP(&TIMEOUT, "timeout", "t", TIMEOUT, "ICMP reply timeout")
	RootCmd.Flags().DurationVarP(&INTERVAL, "interval", "i", INTERVAL, "Wait time between icmp packets before sending new one")
	RootCmd.Flags().IntVar(&MAX_HOPS, "max-hops", MAX_HOPS, "Maximal TTL count")
	RootCmd.Flags().IntVar(&RING_BUFFER_SIZE, "buffer-size", RING_BUFFER_SIZE, "Cached packet buffer size")
}
