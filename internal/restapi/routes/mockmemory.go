package routes

import (
	"net/http"
	restapiutils "quote-service/internal/restapi/utils"
	"quote-service/pkg/logger"
	"strconv"
	"sync"
	"time"
)

// Global variables to hold memory allocations during request processing
var (
	activeAllocations = make(map[string][]byte)
	allocationsMutex  sync.Mutex
)

// HandleAutoScalingDemo
// /api/mock-memory
// This endpoint allocates 10MB of memory and holds it for 10 seconds
// to demonstrate auto-scaling behavior
func HandleAutoScalingDemo(logger logger.Logger) http.HandlerFunc {
	type Response struct {
		Message           string `json:"message"`
		MemoryMB          int    `json:"memory_mb"`
		DurationS         int    `json:"duration_seconds"`
		Timestamp         string `json:"timestamp"`
		RequestID         string `json:"request_id"`
		ActiveAllocations int    `json:"active_allocations"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		requestID := strconv.FormatInt(startTime.UnixNano(), 10)

		// Log the start of autoscaling demo
		logger.InfoWithCtx(r.Context(), "Mock-memory endpoint called",
			"timestamp", startTime.Format(time.RFC3339), "request_id", requestID)

		// Allocate 10MB of memory
		logger.InfoWithCtx(r.Context(), "Allocating 10MB memory for mock-memory demo",
			"request_id", requestID)

		memoryAllocation := make([]byte, 10*1024*1024) // 10MB

		// Fill memory with patterns to prevent optimization and ensure allocation
		for i := 0; i < len(memoryAllocation); i += 4096 { // Page-sized chunks
			for j := 0; j < 4096 && (i+j) < len(memoryAllocation); j++ {
				memoryAllocation[i+j] = byte((i + j) % 256)
			}
		}

		// Store allocation in global map to prevent GC
		allocationsMutex.Lock()
		activeAllocations[requestID] = memoryAllocation
		currentAllocations := len(activeAllocations)
		allocationsMutex.Unlock()

		// Perform CPU work to simulate load
		var checksum uint64
		for i := 0; i < 1000000; i++ {
			checksum += uint64(memoryAllocation[i%len(memoryAllocation)])
		}

		logger.Debug("Memory checksum: " + strconv.FormatUint(checksum, 10))

		// Hold the memory for 10 seconds while doing periodic work
		logger.InfoWithCtx(r.Context(), "Holding memory for 10 seconds",
			"request_id", requestID, "active_allocations", strconv.Itoa(currentAllocations))

		// Periodic work to keep memory active
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		done := make(chan bool, 1)
		go func() {
			for {
				select {
				case <-ticker.C:
					// Access memory periodically to keep it active
					_ = memoryAllocation[0]
					_ = memoryAllocation[len(memoryAllocation)-1]
				case <-done:
					return
				}
			}
		}()

		time.Sleep(10 * time.Second)
		done <- true

		// Clean up: remove from global map
		allocationsMutex.Lock()
		delete(activeAllocations, requestID)
		allocationsMutex.Unlock()

		// Force a small GC cycle by creating temporary allocation
		tempMemory := make([]byte, 1024) // 1KB
		for i := range tempMemory {
			tempMemory[i] = byte(i)
		}
		tempMemory = nil // This can be GC'd

		duration := time.Since(startTime)
		logger.InfoWithCtx(r.Context(), "Mock-memory demo completed",
			"request_id", requestID, "duration_ms", strconv.FormatInt(duration.Milliseconds(), 10),
			"active_allocations", strconv.Itoa(currentAllocations-1))

		resp := Response{
			Message:           "Mock-memory demo completed successfully",
			MemoryMB:          10,
			DurationS:         10,
			Timestamp:         startTime.Format(time.RFC3339),
			RequestID:         requestID,
			ActiveAllocations: currentAllocations,
		}

		restapiutils.WriteJSONResponse(w, http.StatusOK, resp)
	}
}
