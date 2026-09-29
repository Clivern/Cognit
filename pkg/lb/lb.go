// Copyright 2026 Cognit. All rights reserved.
// License can be found in the LICENSE file.

package lb

import (
	"cmp"
	"errors"
	"hash/fnv"
	"math/rand/v2"
	"slices"
	"sync"
)

var (
	ErrEmpty    = errors.New("lb: no instances")
	ErrStrategy = errors.New("lb: unknown strategy")
)

const (
	Random                Strategy = "random"
	RoundRobin            Strategy = "round_robin"
	LeastOutstandingTasks Strategy = "least_outstanding_tasks"
)

// Strategy is how Pick chooses among passing instances.
type Strategy string

// Instance is a live agent instance in the discovery pool.
type Instance struct {
	ID          string
	Datacenter  string
	Outstanding int
}

// Options control how Pick chooses an instance.
type Options struct {
	Strategy         Strategy
	Pool             string
	PreferDatacenter string
	StickyKey        string
}

// Picker selects an instance from a live pool.
type Picker struct {
	mu   sync.Mutex
	next map[string]int
	intn func(int) int
}

// New returns a picker.
func New() *Picker {
	return &Picker{
		next: make(map[string]int),
		intn: rand.IntN,
	}
}

// Pick returns one instance from the live pool.
// StickyKey pins a task to the same instance while that instance stays in the pool.
// PreferDatacenter keeps traffic in that DC when any instance is there.
func (p *Picker) Pick(instances []Instance, opts Options) (*Instance, error) {
	pool := filter(instances, opts.PreferDatacenter)
	if len(pool) == 0 {
		return nil, ErrEmpty
	}

	if opts.StickyKey != "" {
		picked := sticky(pool, opts.StickyKey)
		return &picked, nil
	}

	strategy := opts.Strategy
	if strategy == "" {
		strategy = Random
	}

	var picked Instance
	switch strategy {
	case Random:
		picked = pool[p.intn(len(pool))]
	case RoundRobin:
		picked = p.roundRobin(opts.Pool, pool)
	case LeastOutstandingTasks:
		picked = p.leastOutstanding(pool)
	default:
		return nil, ErrStrategy
	}
	return &picked, nil
}

// Shuffle returns a copy of instances in random order, like Consul DNS.
func (p *Picker) Shuffle(instances []Instance) []Instance {
	out := slices.Clone(instances)
	for i := len(out) - 1; i > 0; i-- {
		j := p.intn(i + 1)
		out[i], out[j] = out[j], out[i]
	}
	return out
}

func (p *Picker) roundRobin(pool string, instances []Instance) Instance {
	ordered := ordered(instances)

	p.mu.Lock()
	i := p.next[pool]
	p.next[pool] = i + 1
	p.mu.Unlock()

	return ordered[i%len(ordered)]
}

func (p *Picker) leastOutstanding(instances []Instance) Instance {
	min := instances[0].Outstanding
	for _, inst := range instances[1:] {
		if inst.Outstanding < min {
			min = inst.Outstanding
		}
	}

	var tied []Instance
	for _, inst := range instances {
		if inst.Outstanding == min {
			tied = append(tied, inst)
		}
	}
	return tied[p.intn(len(tied))]
}

func filter(instances []Instance, datacenter string) []Instance {
	if datacenter == "" || len(instances) == 0 {
		return slices.Clone(instances)
	}

	var matched []Instance
	for _, inst := range instances {
		if inst.Datacenter == datacenter {
			matched = append(matched, inst)
		}
	}
	if len(matched) == 0 {
		return slices.Clone(instances)
	}
	return matched
}

func ordered(instances []Instance) []Instance {
	out := slices.Clone(instances)
	slices.SortFunc(out, func(a, b Instance) int {
		return cmp.Compare(a.ID, b.ID)
	})
	return out
}

func sticky(instances []Instance, key string) Instance {
	pool := ordered(instances)
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return pool[int(h.Sum32())%len(pool)]
}
