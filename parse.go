package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

const marker = "go 1.11 trace\x00\x00\x00"

func run() error {
	f, err := os.Open("trace.out")
	if err != nil {
		return err
	}
	defer f.Close()
	res, err := Parse(f, "")
	if err != nil {
		return err
	}
	for i, ev := range res.Events {
		desc := EventDescriptions[ev.Type]
		fmt.Printf("%3d: %-12s G%-2d P%-2d ", i+1, desc.Name, ev.G, ev.P)
		stk := ev.Stk
		if len(stk) == 0 {
			stk = res.Stacks[ev.StkID]
		}
		if len(stk) > 0 {
			frame1 := stk[len(stk)-1]
			frame2 := stk[0]
			if frame1 != frame2 {
				fmt.Printf("%s:%d -> %s:%d ", frame1.File, frame1.Line, frame2.File, frame2.Line)
			} else {
				fmt.Printf("%s:%d ", frame1.File, frame1.Line)
			}
		}
		for i, a := range desc.Args {
			fmt.Printf("%s=%d ", a, ev.Args[i])
		}
		for i, a := range desc.SArgs {
			fmt.Printf("%s=%s ", a, ev.SArgs[i])
		}
		if ev.Link != nil {
			for j, ev2 := range res.Events {
				if ev2 == ev.Link {
					desc = EventDescriptions[ev2.Type]
					fmt.Printf("linked %s@%d ", desc.Name, j+1)
					break
				}
			}
		}
		fmt.Println()
	}
	fmt.Println()
	var buf strings.Builder
	g := func(v int64, m string) {
		if v > 0 {
			buf.WriteString(fmt.Sprintf("%s=%d ", m, v))
		}
	}
	stats := GoroutineStats(res.Events)
	keys := make([]int, 0, len(stats))
	for k := range stats {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	for _, i := range keys {
		buf.Reset()
		v := stats[uint64(i)]
		stat := v.GExecutionStat
		g(stat.ExecTime, "Exec")
		g(stat.SchedWaitTime, "SchedWait")
		g(stat.IOTime, "IO")
		g(stat.BlockTime, "Block")
		g(stat.SyscallTime, "Syscall")
		g(stat.GCTime, "GC")
		g(stat.SweepTime, "Sweep")
		g(stat.TotalTime, "Total")

		fmt.Printf("%2d: %s%s\n", v.ID, buf.String(), v.Name)
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
