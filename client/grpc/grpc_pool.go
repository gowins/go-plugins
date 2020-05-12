package grpc

import (
	"sync"
	"time"

	"google.golang.org/grpc"
)

const (
	// 偿试清理间隔
	cleanupInterval = time.Minute * 10
)

type pool struct {
	size int
	ttl  int64

	sync.Mutex
	conns map[string]*poolManager
}

type poolConn struct {
	*grpc.ClientConn

	// 创建状态
	newCreated bool
	created    time.Time

	// 使用中的引用计数
	refCount int64

	// 用来决定是否关闭
	closable bool

	next *poolConn
	prev *poolConn
}

func newPool(size int, ttl time.Duration) *pool {
	pool := &pool{
		size:  size,
		ttl:   int64(ttl.Seconds()),
		conns: make(map[string]*poolManager),
	}

	go pool.cleanup()

	return pool
}

func (p *pool) getConn(addr string, opts ...grpc.DialOption) (*poolConn, error) {
	return p.getManager(addr).get(opts...)
}

func (p *pool) release(addr string, conn *poolConn, err error) {
	// otherwise put it back for reuse
	p.getManager(addr).put(conn, err)
}

func (p *pool) getManager(addr string) *poolManager {
	p.Lock()
	defer p.Unlock()

	manager := p.conns[addr]
	if manager == nil {
		manager = newManager(addr, p.size, p.ttl)
		p.conns[addr] = manager
	}

	return manager
}

func (p *pool) cleanup() {
	timer := time.NewTicker(cleanupInterval)
	for range timer.C {
		cleans := p.findCanCleanups()

		if len(cleans) == 0 {
			continue
		}

		for _, manager := range cleans {
			manager.cleanup()
		}
	}
}

func (p *pool) findCanCleanups() map[string]*poolManager {
	p.Lock()
	defer p.Unlock()

	cleans := make(map[string]*poolManager)

	for addr, manager := range p.conns {
		if manager != nil && manager.canCleanup() {
			cleans[addr] = manager
			delete(p.conns, addr)
		}
	}

	return cleans
}
