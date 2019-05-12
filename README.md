# rund

[![GoDoc](https://godoc.org/github.com/andy2046/rund?status.svg)](https://godoc.org/github.com/andy2046/rund)

rund is a task dependency scheduler to run all tasks in parallel topological order, using the semantic of Directed Acyclic Graph.

## Install

```
go get github.com/andy2046/rund
```

## Example

```go
r := rund.New()

op1 := rund.NewFuncOperator(func() error {
	fmt.Println("1 will run before 2 and 3")
	return nil
})
op2 := rund.NewFuncOperator(func() error {
	fmt.Println("2 and 3 will run in parallel before 4")
	return nil
})
op3 := rund.NewFuncOperator(func() error {
	fmt.Println("2 and 3 will run in parallel before 4")
	return nil
})
op4 := rund.NewFuncOperator(func() error {
	fmt.Println("4 will run after 2 and 3")
	return errors.New("4 is broken")
})
op5 := rund.NewFuncOperator(func() error {
	fmt.Println("5 will never run")
	return nil
})

r.AddNode("1", op1)
r.AddNode("2", op2)
r.AddNode("3", op3)
r.AddNode("4", op4)
r.AddNode("5", op5)

r.AddEdge("1", "2")
r.AddEdge("1", "3")
r.AddEdge("2", "4")
r.AddEdge("3", "4")
r.AddEdge("4", "5")

//     2
//   /   \
//  1     4 - 5
//   \   /
//     3

fmt.Printf("the result: %v\n", r.Run())
// Output:
// 1 will run before 2 and 3
// 2 and 3 will run in parallel before 4
// 2 and 3 will run in parallel before 4
// 4 will run after 2 and 3
// the result: 4 is broken

```
