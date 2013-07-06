// Copyright (C) 2013 Max Riveiro
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// This package is just a tiny handler function for printing some statistics (mostly memory) for running application with just a tiny call of SIGINFO - ^T shortcut.
//
// Example
//
// func main() {
//   // Let's enable all the Stats
//   sigstats.EnableSigStats(&sigstats.SigStatsSettings{
//     NumCPU:       true,
//     NumCGOCalls:  true,
//     NumGoroutine: true,
//     MemStats:     true,
//     GoVersion:    true,
//     GCStats:      true,
//   })
//   select {}
// }
package sigstats

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"syscall"
	"time"
)

// Syscall signal number
const SIGINFO = syscall.Signal(29)

// Settings struct. You need it for enabling Stats' outputs.
type SigStatsSettings struct {
	// Booleans represintig corresponding Stats outputs
	NumCPU, NumCGOCalls, NumGoroutine, MemStats, GoVersion, GCStats bool
}

// Init handler on any system signal. Basic function.
func initSignalHandler(sig os.Signal, handler func()) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sig)

	go func() {
		for _ = range ch {
			handler()
		}
	}()
}

// Init custom SIGINFO handlers. May be helpful for custom Stats' outputs.
func InitInfoHandler(handler func()) {
	initSignalHandler(SIGINFO, handler)
}

// Init the SigStats SIGINFO handler
func EnableSigStats(s *SigStatsSettings) {
	InitInfoHandler(func() {
		fmt.Println()

		if s != nil {
			if s.GoVersion {
				fmt.Printf("Go version: %s\n\n", runtime.Version())
			}

			if s.NumCPU {
				fmt.Printf("Number of CPUs: %d\n\n", runtime.NumCPU())
			}

			if s.NumCGOCalls {
				fmt.Printf("Number of CGO Calls: %d\n\n", runtime.NumCgoCall())
			}

			if s.NumGoroutine {
				fmt.Printf("Number of existing Goroutines: %d\n\n", runtime.NumGoroutine())
			}

			if s.MemStats {
				m := &runtime.MemStats{}
				runtime.ReadMemStats(m)

				last_gc := time.Duration(m.LastGC) / time.Second

				fmt.Printf("Memstats — Alloc size: %d bytes\n", m.Alloc)
				fmt.Printf("Memstats — TotalAlloc size: %d bytes\n", m.TotalAlloc)
				fmt.Printf("Memstats — Sys size: %d bytes\n", m.Sys)
				fmt.Printf("Memstats — Lookups number: %d\n", m.Lookups)
				fmt.Printf("Memstats — Mallocs number: %d\n", m.Mallocs)
				fmt.Printf("Memstats — Frees number: %d\n", m.Frees)
				fmt.Println()

				fmt.Printf("Memstats - Heap — HeapAlloc size: %d bytes\n", m.HeapAlloc)
				fmt.Printf("Memstats - Heap — HeapSys size: %d bytes\n", m.HeapSys)
				fmt.Printf("Memstats - Heap — HeapIdle size: %d bytes\n", m.HeapIdle)
				fmt.Printf("Memstats - Heap — HeapInuse size: %d bytes\n", m.HeapInuse)
				fmt.Printf("Memstats - Heap — HeapReleased size: %d bytes\n", m.HeapReleased)
				fmt.Printf("Memstats - Heap — HeapObjects size: %d bytes\n", m.HeapObjects)
				fmt.Println()

				fmt.Printf("Memstats - Allocator — StackInuse size: %d bytes\n", m.StackInuse)
				fmt.Printf("Memstats - Allocator — StackSys size: %d bytes\n", m.StackSys)
				fmt.Printf("Memstats - Allocator — MSpanInuse size: %d bytes\n", m.MSpanInuse)
				fmt.Printf("Memstats - Allocator — MSpanSys size: %d bytes\n", m.MSpanSys)
				fmt.Printf("Memstats - Allocator — MCacheInuse size: %d bytes\n", m.MCacheInuse)
				fmt.Printf("Memstats - Allocator — MCacheSys size: %d bytes\n", m.MCacheSys)
				fmt.Printf("Memstats - Allocator — BuckHashSys size: %d bytes\n", m.BuckHashSys)
				fmt.Println()

				fmt.Printf("Memstats - GC — NextGC time: %d bytes\n", m.NextGC)
				fmt.Printf("Memstats - GC — LastGC run time: %s\n", last_gc)
				fmt.Printf("Memstats - GC — PauseTotalNs number: %d\n", m.PauseTotalNs)
				fmt.Printf("Memstats - GC — EnableGC: %t\n", m.EnableGC)
				fmt.Printf("Memstats - GC — DebugGC: %t\n", m.DebugGC)
				fmt.Println()
			}

			if s.GCStats {
				g := &debug.GCStats{}
				debug.ReadGCStats(g)

				pause_total := time.Duration(g.PauseTotal) / time.Microsecond
				pause_last := time.Duration(g.Pause[0]) / time.Microsecond

				fmt.Printf("GC - LastGC at: %s\n", g.LastGC)
				fmt.Printf("GC - NumGC number: %d\n", g.NumGC)
				fmt.Printf("GC - PauseTotal time: %s\n", pause_total)
				fmt.Printf("GC - Last Pause time: %s\n", pause_last)
			}
		}
	})
}
