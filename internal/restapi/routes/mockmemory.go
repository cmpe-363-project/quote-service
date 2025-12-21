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
// This endpoint allocates dynamic memory and holds it for specified duration
// to demonstrate auto-scaling behavior. Query parameters:
//   - memory_mb: memory to allocate in MB (1-1000, default: 10)
//   - duration_seconds: duration to hold memory in seconds (1-300, default: 10)
func HandleAutoScalingDemo(logger logger.Logger) http.HandlerFunc {
	type Response struct {
		Message           string `json:"message"`
		MemoryMB          int    `json:"memory_mb"`
		DurationS         int    `json:"duration_seconds"`
		Timestamp         string `json:"timestamp"`
		RequestID         string `json:"request_id"`
		ActiveAllocations int    `json:"active_allocations"`
	}

	type RequestError struct {
		Error   string `json:"error"`
		Code    int    `json:"code"`
		Details string `json:"details,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		requestID := strconv.FormatInt(startTime.UnixNano(), 10)

		// Parse query parameters
		memoryMBStr := r.URL.Query().Get("memory_mb")
		durationSecondsStr := r.URL.Query().Get("duration_seconds")

		// Set defaults
		memoryMB := 10  // Default 10MB
		durationS := 10 // Default 10 seconds

		var err error

		// Parse and validate memory_mb parameter
		if memoryMBStr != "" {
			memoryMB, err = strconv.Atoi(memoryMBStr)
			if err != nil || memoryMB <= 0 {
				logger.WarnWithCtx(r.Context(), "Invalid memory_mb parameter",
					"memory_mb", memoryMBStr, "request_id", requestID)
				errorResp := RequestError{
					Error:   "Invalid memory_mb parameter",
					Code:    400,
					Details: "memory_mb must be a positive integer",
				}
				restapiutils.WriteJSONResponse(w, http.StatusBadRequest, errorResp)
				return
			}
			if memoryMB > 1000 { // Limit to 1GB max
				logger.WarnWithCtx(r.Context(), "Memory limit exceeded",
					"memory_mb", strconv.Itoa(memoryMB), "request_id", requestID)
				errorResp := RequestError{
					Error:   "Memory limit exceeded",
					Code:    400,
					Details: "memory_mb cannot exceed 1000MB (1GB)",
				}
				restapiutils.WriteJSONResponse(w, http.StatusBadRequest, errorResp)
				return
			}
		}

		// Parse and validate duration_seconds parameter
		if durationSecondsStr != "" {
			durationS, err = strconv.Atoi(durationSecondsStr)
			if err != nil || durationS <= 0 {
				logger.WarnWithCtx(r.Context(), "Invalid duration_seconds parameter",
					"duration_seconds", durationSecondsStr, "request_id", requestID)
				errorResp := RequestError{
					Error:   "Invalid duration_seconds parameter",
					Code:    400,
					Details: "duration_seconds must be a positive integer",
				}
				restapiutils.WriteJSONResponse(w, http.StatusBadRequest, errorResp)
				return
			}
			if durationS > 300 { // Limit to 5 minutes max
				logger.WarnWithCtx(r.Context(), "Duration limit exceeded",
					"duration_seconds", strconv.Itoa(durationS), "request_id", requestID)
				errorResp := RequestError{
					Error:   "Duration limit exceeded",
					Code:    400,
					Details: "duration_seconds cannot exceed 300 (5 minutes)",
				}
				restapiutils.WriteJSONResponse(w, http.StatusBadRequest, errorResp)
				return
			}
		}

		// Log the start of autoscaling demo
		logger.InfoWithCtx(r.Context(), "Mock-memory endpoint called",
			"timestamp", startTime.Format(time.RFC3339), "request_id", requestID,
			"memory_mb", strconv.Itoa(memoryMB), "duration_seconds", strconv.Itoa(durationS))

		// Allocate memory based on parameter
		logger.InfoWithCtx(r.Context(), "Allocating memory for mock-memory demo",
			"request_id", requestID, "memory_mb", strconv.Itoa(memoryMB))

		memoryAllocation := make([]byte, memoryMB*1024*1024) // Dynamic memory allocation

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

		// Hold the memory for specified duration while doing periodic work
		logger.InfoWithCtx(r.Context(), "Holding memory for specified duration",
			"request_id", requestID, "active_allocations", strconv.Itoa(currentAllocations),
			"duration_seconds", strconv.Itoa(durationS))

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

		time.Sleep(time.Duration(durationS) * time.Second)
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
			MemoryMB:          memoryMB,
			DurationS:         durationS,
			Timestamp:         startTime.Format(time.RFC3339),
			RequestID:         requestID,
			ActiveAllocations: currentAllocations,
		}

		restapiutils.WriteJSONResponse(w, http.StatusOK, resp)
	}
}
