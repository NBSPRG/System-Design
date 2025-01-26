# Lamport Logical Clocks Simulation with Multiprocessing

This project demonstrates the implementation of Lamport logical clocks in a multiprocessing environment. Lamport clocks are used to maintain causal ordering of events in a distributed system.

## Key Features
- **Logical Event Ordering**: Tracks the causal order of events in separate processes.
- **Message Passing**: Processes communicate timestamps to synchronize clocks.
- **Concurrency Handling**: Simulates real-time multiprocessing with independent processes.

## How It Works
1. Each process maintains its own logical clock.
2. Events are categorized as:
   - Local events (clock increment within a process).
   - Message passing events (sending/receiving messages).
3. When a message is sent, the sender includes its clock value.
4. The receiver updates its clock to the maximum of its own value and the received timestamp, ensuring causal consistency.

