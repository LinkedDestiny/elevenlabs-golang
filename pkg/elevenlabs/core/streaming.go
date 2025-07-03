package core

import (
	"bufio"
	"context"
	"io"
	"net/http"
)

// StreamChunk represents a chunk of streaming data
type StreamChunk struct {
	Data []byte
	Err  error
}

// StreamResponse streams an HTTP response body in chunks
func StreamResponse(resp *http.Response, chunkSize int) <-chan []byte {
	ch := make(chan []byte)

	go func() {
		defer close(ch)
		defer resp.Body.Close()

		buffer := make([]byte, chunkSize)
		for {
			n, err := resp.Body.Read(buffer)
			if n > 0 {
				// Make a copy of the data to avoid race conditions
				chunk := make([]byte, n)
				copy(chunk, buffer[:n])
				ch <- chunk
			}
			if err != nil {
				if err != io.EOF {
					// For non-EOF errors, we could send them through a separate error channel
					// For now, we'll just break out of the loop
				}
				break
			}
		}
	}()

	return ch
}

// StreamWithContext streams an HTTP response with context support
func StreamWithContext(ctx context.Context, resp *http.Response) <-chan StreamChunk {
	ch := make(chan StreamChunk)

	go func() {
		defer close(ch)
		defer resp.Body.Close()

		reader := bufio.NewReader(resp.Body)
		buffer := make([]byte, 1024)

		for {
			select {
			case <-ctx.Done():
				ch <- StreamChunk{Err: ctx.Err()}
				return
			default:
				n, err := reader.Read(buffer)
				if n > 0 {
					// Make a copy of the data
					chunk := make([]byte, n)
					copy(chunk, buffer[:n])
					ch <- StreamChunk{Data: chunk}
				}
				if err != nil {
					if err != io.EOF {
						ch <- StreamChunk{Err: err}
					}
					return
				}
			}
		}
	}()

	return ch
}

// StreamLines streams an HTTP response line by line
func StreamLines(ctx context.Context, resp *http.Response) <-chan StreamChunk {
	ch := make(chan StreamChunk)

	go func() {
		defer close(ch)
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				ch <- StreamChunk{Err: ctx.Err()}
				return
			default:
				line := scanner.Text()
				ch <- StreamChunk{Data: []byte(line)}
			}
		}

		if err := scanner.Err(); err != nil {
			ch <- StreamChunk{Err: err}
		}
	}()

	return ch
}
