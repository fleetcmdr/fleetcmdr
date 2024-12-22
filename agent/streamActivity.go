package main

func (d *agentDaemon) startStreamingActivity() {
	// begin sending updates every second to server

	// CPU (consumed of total)
	// Memory (consumed of total)
	// Disk IO (io per second, latency?)
	// Network (bytes/sec)
	// Network latency? (tell server when we intended to send the data, and calculate the diff?)
}
