# SigStats

Package sigstats is a tiny handler function for printing some statistics (mostly memory) for running application with just a tiny call of **SIGINFO** - **^T** shortcut.

You can view documentation on [godoc.org](http://godoc.org/github.com/kavu/go-sigstats "go-sigstats documentation").

## Example

```go

package main

import "github.com/kavu/go-sigstats"

func main() {
  // Let's enable all the Stats
  sigstats.EnableSigStats(&sigstats.SigStatsSettings{
    NumCPU:       true,
    NumCGOCalls:  true,
    NumGoroutine: true,
    MemStats:     true,
    GoVersion:    true,
    GCStats:      true,
  })

  select {}
}
```
