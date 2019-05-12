// Package rund implements a DAG task dependency scheduler.
package rund

import "errors"

type (
	// Rund adds `Operator` into the graph as node and edge,
	// `Operator` runs in parallel topological order.
	Rund struct {
		ops    map[string]Operator
		graph  map[string][]string
		status Status
	}

	// Status for `Operator`.
	Status int

	result struct {
		name string
		err  error
	}
)

// `Operator` Status List.
const (
	PENDING Status = iota
	RUNNING
	SKIPPED
	SUCCESS
	FAILED
)

var (
	errMissingOperator     = errors.New("missing Operator")
	errCircularDepDetected = errors.New("circular dependency detected")
)

// New returns a new `Rund`.
func New() *Rund {
	return &Rund{
		ops:    make(map[string]Operator),
		graph:  make(map[string][]string),
		status: PENDING,
	}
}

// AddNode adds an `Operator` as a node in the graph,
// `Operator` name should be unique.
func (d *Rund) AddNode(name string, op Operator) {
	if d.ops == nil {
		d.ops = make(map[string]Operator)
	}
	d.ops[name] = op
}

// AddEdge establishes a dependency between two nodes in the graph,
// `from` node will be executed before `to` node.
func (d *Rund) AddEdge(from, to string) {
	if d.graph == nil {
		d.graph = make(map[string][]string)
	}
	d.graph[from] = append(d.graph[from], to)
}

// Run validates the graph then runs each node in parallel topological order,
// if any node Operator returns error, no more node will be scheduled.
func (d *Rund) Run() error {
	d.status = RUNNING

	if len(d.ops) == 0 {
		d.status = SKIPPED
		return nil
	}

	deps := make(map[string]int)
	for node, edges := range d.graph {
		// each node has one associated Operator
		if _, ok := d.ops[node]; !ok {
			d.status = FAILED
			return errMissingOperator
		}
		for _, node := range edges {
			if _, ok := d.ops[node]; !ok {
				d.status = FAILED
				return errMissingOperator
			}
			deps[node]++
		}
	}

	if d.hasCircularDep() {
		d.status = FAILED
		return errCircularDepDetected
	}

	runningOperator, resultCh := 0, make(chan result, len(d.ops))

	for name := range d.ops {
		if deps[name] == 0 {
			runningOperator++
			start(name, d.ops[name], resultCh)
		}
	}

	var err error
	// wait for all running Operator to complete
	for runningOperator > 0 {
		res := <-resultCh
		runningOperator--

		if res.err != nil && err == nil {
			d.status = FAILED
			err = res.err
		}

		if err != nil {
			continue
		}

		// start node with deps fully resolved
		for _, node := range d.graph[res.name] {
			deps[node]--
			if deps[node] == 0 {
				runningOperator++
				start(node, d.ops[node], resultCh)
			}
		}
	}

	if err == nil {
		d.status = SUCCESS
	}
	return err
}

func start(name string, op Operator, resultCh chan<- result) {
	go func() {
		resultCh <- result{
			name: name,
			err:  op.Run(),
		}
	}()
}

func (d *Rund) hasCircularDep() bool {
	visited := make(map[string]bool)
	recurStack := make(map[string]bool)

	for node := range d.graph {
		if !visited[node] {
			if d.circularDepHelper(node, visited, recurStack) {
				return true
			}
		}
	}
	return false
}

func (d *Rund) circularDepHelper(node string, visited, recurStack map[string]bool) bool {
	visited[node] = true
	recurStack[node] = true

	for _, v := range d.graph[node] {
		if !visited[v] {
			if d.circularDepHelper(v, visited, recurStack) {
				return true
			}
		} else if recurStack[v] {
			// if we already visited this node in this recursion stack, we have a cycle
			return true
		}
	}

	recurStack[node] = false
	return false
}
