# 23-08-2026

## What is the difference between `sync.Mutex` and `sync.RWMutex`?

- **sync.Mutex**
    - `sync.Mutex` is a basic mutual exclusion lock that allows only one goroutine to hold the lock at a time
    - All goroutines (readers and writers) block equally
    - Exclusive lock - only one reader OR writer can hold it, so no distinction between read and write operations
- **sync.RWMutex**
    - `sync.RWMutex` is a read-write mutual exclusion lock that allows multiple readers to hold the lock at the same time
    - Multiple readers can hold the lock simultaneously but one writer can hold the lock exclusively
    - Writers block all readers and other writers (Writers have priority to prevent **writer starvation**)
    - More overhead than sync.Mutex

### Best Practices

- Start with sync.Mutex: It's simpler and often sufficient
- Profile before optimizing: Only use RWMutex if you've measured and confirmed read-heavy workloads
- Keep critical sections small: Reduce lock contention regardless of type
- Use defer for unlocking: Prevents deadlocks from panics or early returns
