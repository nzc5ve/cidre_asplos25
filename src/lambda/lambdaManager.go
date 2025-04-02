package lambda

import (
	"container/list"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/open-lambda/open-lambda/ol/common"
	"github.com/open-lambda/open-lambda/ol/sandbox"
)

// LambdaMgr provides thread-safe getting of lambda functions and collects all
// lambda subsystems (resource pullers and sandbox pools) in one place
type LambdaMgr struct {
	// subsystems (these are thread safe)
	sbPool sandbox.SandboxPool
	*DepTracer
	*PackagePuller // depends on sbPool and DepTracer
	*ImportCache   // depends PackagePuller
	*HandlerPuller // depends on sbPool and ImportCache[optional]

	// storage dirs that we manage
	codeDirs    *common.DirMaker
	scratchDirs *common.DirMaker

	// thread-safe map from a lambda's name to its LambdaFunc
	mapMutex sync.Mutex
	lfuncMap map[string]*LambdaFunc
}

// represents an HTTP request to be handled by a lambda instance
type Invocation struct {
	w http.ResponseWriter
	r *http.Request

	// signal to client that response has been written to w
	done chan bool

	// queuing delay for the requests
	queueStart time.Time
	queuingMs int

	// how many milliseconds did ServeHTTP take?  (doesn't count
	// queue time or Sandbox init)
	execMs int
}

// NewLambdaMgr creates a new LambdaMgr instance and initializes its subsystems.
func NewLambdaMgr() (res *LambdaMgr, err error) {
	mgr := &LambdaMgr{
		lfuncMap: make(map[string]*LambdaFunc),
	}
	defer func() {
		if err != nil {
			log.Printf("Cleanup Lambda Manager due to error: %v", err)
			mgr.Cleanup()
		}
	}()

	mgr.codeDirs, err = common.NewDirMaker("code", common.Conf.Storage.Code.Mode())
	if err != nil {
		return nil, err
	}
	mgr.scratchDirs, err = common.NewDirMaker("scratch", common.Conf.Storage.Scratch.Mode())
	if err != nil {
		return nil, err
	}

	log.Printf("Creating SandboxPool")
	mgr.sbPool, err = sandbox.SandboxPoolFromConfig("sandboxes", common.Conf.Mem_pool_mb)
	if err != nil {
		return nil, err
	}

	log.Printf("Creating DepTracer")
	mgr.DepTracer, err = NewDepTracer(filepath.Join(common.Conf.Worker_dir, "dep-trace.json"))
	if err != nil {
		return nil, err
	}

	log.Printf("Creating PackagePuller")
	mgr.PackagePuller, err = NewPackagePuller(mgr.sbPool, mgr.DepTracer)
	if err != nil {
		return nil, err
	}

	if common.Conf.Features.Import_cache {
		log.Printf("Creating ImportCache")
		mgr.ImportCache, err = NewImportCache(mgr.codeDirs, mgr.scratchDirs, mgr.sbPool, mgr.PackagePuller)
		if err != nil {
			return nil, err
		}
	}

	log.Printf("Creating HandlerPuller")
	mgr.HandlerPuller, err = NewHandlerPuller(mgr.codeDirs)
	if err != nil {
		return nil, err
	}

	return mgr, nil
}

// Returns an existing instance (if there is one), or creates a new one
func (mgr *LambdaMgr) Get(name string) (f *LambdaFunc) {
	mgr.mapMutex.Lock()
	defer mgr.mapMutex.Unlock()

	f = mgr.lfuncMap[name]

	if f == nil {
		f = &LambdaFunc{
			lmgr:      mgr,
			name:      name,
			// make these configurable
			funcChan:  make(chan *Invocation, 256),
			instChan:  make(chan *Invocation, 256),
			doneChan:  make(chan *Invocation, 256),
			delyChan:  make(chan *Invocation, 256),
			instances: list.New(),
			coldPath:  true,
			killChan:  make(chan chan bool, 1),
		}

		go f.Task()
		mgr.lfuncMap[name] = f
	}

	return f
}

// Debug returns the debug information of the sandbox pool.
func (mgr *LambdaMgr) Debug() string {
	return mgr.sbPool.DebugString() + "\n"
}

// DumpStatsToLog logs the profiling information of the LambdaMgr.
func (mgr *LambdaMgr) DumpStatsToLog() {
	snapshot := common.SnapshotStats()

	sec := func(name string) (float64) {
		return float64(snapshot[name+".cnt"] * snapshot[name+".ms-avg"]) / 1000
	}

	time := func(indent int, name string, parent string) {
		selftime := sec(name)
		ptime := sec(parent)
		tabs := strings.Repeat("\t", indent)
		if ptime > 0 {
			log.Printf("%s%s: %.3f (%.1f%%)", tabs, name, selftime, selftime/ptime*100)
		} else {
			log.Printf("%s%s: %.3f", tabs, name, selftime)
		}
	}

	log.Printf("Request Profiling (cumulative seconds):")
	time(0, "LambdaFunc.Invoke", "")

	time(1, "LambdaInstance-WaitSandbox", "LambdaFunc.Invoke")
	time(2, "LambdaInstance-WaitSandbox-Unpause", "LambdaInstance-WaitSandbox")
	time(2, "LambdaInstance-WaitSandbox-NoImportCache", "LambdaInstance-WaitSandbox")
	time(2, "ImportCache.Create", "LambdaInstance-WaitSandbox")
	time(3, "ImportCache.root.Lookup", "ImportCache.Create")
	time(3, "ImportCache.createChildSandboxFromNode", "ImportCache.Create")
	time(4, "ImportCache.getSandboxInNode", "ImportCache.createChildSandboxFromNode")
	time(4, "ImportCache.createChildSandboxFromNode:childSandboxPool.Create",
		"ImportCache.createChildSandboxFromNode")
	time(4, "ImportCache.putSandboxInNode", "ImportCache.createChildSandboxFromNode")
	time(5, "ImportCache.putSandboxInNode:Lock", "ImportCache.putSandboxInNode")
	time(5, "ImportCache.putSandboxInNode:Pause", "ImportCache.putSandboxInNode")
	time(1, "LambdaInstance-ServeRequests", "LambdaFunc.Invoke")
	time(2, "LambdaInstance-RoundTrip", "LambdaInstance-ServeRequests")
}

// Cleanup performs cleanup operations for the LambdaMgr and its subsystems.
func (mgr *LambdaMgr) Cleanup() {
	mgr.mapMutex.Lock() // don't unlock, because this shouldn't be used anymore

	mgr.DumpStatsToLog()

	// HandlerPuller+PackagePuller requires no cleanup

	// 1. cleanup handler Sandboxes
	// 2. cleanup Zygote Sandboxes (after the handlers, which depend on the Zygotes)
	// 3. cleanup SandboxPool underlying both of above
	for _, f := range mgr.lfuncMap {
		log.Printf("Kill function: %s", f.name)
		f.Kill()
	}

	if mgr.ImportCache != nil {
		mgr.ImportCache.Cleanup()
	}

	if mgr.sbPool != nil {
		mgr.sbPool.Cleanup() // assumes all Sandboxes are gone
	}

	// cleanup DepTracer
	if mgr.DepTracer != nil {
		mgr.DepTracer.Cleanup()
	}

	if mgr.codeDirs != nil {
		mgr.codeDirs.Cleanup()
	}

	if mgr.scratchDirs != nil {
		mgr.scratchDirs.Cleanup()
	}
}
