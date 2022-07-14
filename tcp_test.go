package echo_test

import (
	"net"
	"runtime"
	"sync"
	"testing"
)

var workers = 4 * runtime.NumCPU()

func BenchmarkAdaptive(b *testing.B) {
	// pushing this higher causes the benchmark not terminate
	const batch = 2 * 1024

	var conns []net.Conn
	for i := 0; i < 1000; i++ {
		conn, err := net.Dial("tcp", "127.0.0.1:1122")
		if err != nil {
			b.Error(err)
			return
		}
		conns = append(conns, conn)
	}
	defer func() {
		for _, conn := range conns {
			conn.Close()
		}
	}()

	var (
		wg      sync.WaitGroup
		ch      = make(chan []byte)
		payload = make([]byte, batch)
		maxGort = 0
	)
	for i := 0; i < b.N; i += batch {
		select {
		case ch <- payload:
			continue
		default:
		}

		if gort := runtime.NumGoroutine(); maxGort < gort {
			maxGort = gort
		}

		wg.Add(1)
		go func(conn net.Conn) {
			defer wg.Done()

			buf := make([]byte, batch)
			for v := range ch {
				_, err := conn.Write(v)
				if err != nil {
					b.Error(err)
					break
				}
				receive := len(v)
				for receive > 0 {
					n, err := conn.Read(buf)
					if err != nil {
						b.Error(err)
						return
					}
					receive -= n
				}
				if receive < 0 {
					b.Errorf("received %d extra bytes", -receive)
				}
			}
		}(conns[i%len(conns)])

		ch <- payload
	}
	close(ch)
	wg.Wait()

	b.Logf("max goroutines: %d\n", maxGort)
}

func BenchmarkKnownCount(b *testing.B) {
	const batch = 64 * 1024

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		count := b.N / workers
		if i <= b.N%workers {
			count++
		}

		wg.Add(1)
		go func(count int) {
			defer wg.Done()

			conn, err := net.Dial("tcp", "127.0.0.1:1122")
			if err != nil {
				b.Error(err)
				return
			}
			defer conn.Close()

			buf := make([]byte, batch)
			for i := 0; i < count; i += batch {
				_, err := conn.Write(buf)
				if err != nil {
					b.Error(err)
					return
				}
				receive := batch
				for receive > 0 {
					n, err := conn.Read(buf)
					if err != nil {
						b.Error(err)
						return
					}
					receive -= n
				}
				if receive < 0 {
					b.Errorf("received %d extra bytes", -receive)
				}
			}
		}(count)
	}
	wg.Wait()
}

func BenchmarkFixedWorkers(b *testing.B) {
	const batch = 64 * 1024

	var (
		wg sync.WaitGroup
		ch = make(chan []byte)
	)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			conn, err := net.Dial("tcp", "127.0.0.1:1122")
			if err != nil {
				b.Error(err)
				return
			}
			defer conn.Close()

			buf := make([]byte, batch)
			for val := range ch {
				_, err = conn.Write(val)
				if err != nil {
					b.Error(err)
					return
				}

				receive := len(val)
				for receive > 0 {
					n, err := conn.Read(buf)
					if err != nil {
						b.Error(err)
						return
					}
					receive -= n
				}
				if receive < 0 {
					b.Errorf("received %d extra bytes", -receive)
				}
			}
		}()
	}

	for i := 0; i < b.N; i += batch {
		ch <- make([]byte, batch)
	}
	close(ch)
	wg.Wait()
}
