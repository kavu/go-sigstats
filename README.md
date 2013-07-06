# SigStats

Package sigstats is a tiny handler function for printing some statistics (mostly memory) for running application with just a tiny call of SIGINFO - ^T shortcut.

## Example

```go
func main() {
  // Let's enable all the Stats
  EnableSigStats(&SigStatsSettings{
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
