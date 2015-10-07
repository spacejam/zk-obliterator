package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/samuel/go-zookeeper/zk"
	"github.com/spacejam/loghisto"
)

var reqCnt = int32(0)

func gen(size int) []byte {
	ret := ""
	for i := 0; i < size; i++ {
		ret += "A"
	}
	return []byte(ret)
}

func bench(conn *zk.Conn, value []byte, ratio float64) {
	choice := rand.Float64()
	atomic.AddInt32(&reqCnt, 1)
	if choice >= ratio {
		// write
		conn.Delete("/_test_load_obliterator", 0)
		conn.Create("/_test_load_obliterator", value, 1, zk.WorldACL(zk.PermAll))
	} else {
		// read
		conn.Get("/_test_load_obliterator")
	}
}

func main() {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	obliteration := flag.Int("concurrency", 10, "threads and connections to use for load generation")
	host := flag.String("zk", "master.mesos:2181", "host:port for zk")
	size := flag.Int("size", 1024, "bytes per key written")
	ratio := flag.Float64("ratio", 0.2, "0 to 1 ratio of reads to writes.  0 is all writes, 1 is all reads.")
	flag.Parse()

	value := gen(*size)

	conns := []*zk.Conn{}
	for i := 0; i < *obliteration; i++ {
		cli, _, err := zk.Connect([]string{*host}, 5*time.Second)
		if err != nil {
			fmt.Printf("error connecting to zk: %v\n", err)
			os.Exit(1)
		}
		conns = append(conns, cli)
	}

	doRpc := func() {
		cli := conns[rand.Intn(len(conns))]
		bench(cli, value, *ratio)
	}

	loghisto.PrintBenchmark("benchmark1234", uint(*obliteration), doRpc)
}
