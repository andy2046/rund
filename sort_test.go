package rund

import (
	"strings"
	"testing"
)

func TestTopoSort(t *testing.T) {
	expectedOrder := "4-5-0-2-3-1"
	graph, ops := make(map[string][]string), make(map[string]Operator)
	for _, n := range []string{"0", "1", "2", "3", "4", "5"} {
		ops[n] = nil
	}
	for from, to := range map[string][]string{
		"5": []string{"2", "0"},
		"4": []string{"0", "1"},
		"2": []string{"3"},
		"3": []string{"1"},
	} {
		graph[from] = append(graph[from], to...)
	}

	result := topoSort(graph, ops)
	resultOrder := strings.Join(result, "-")
	if resultOrder != expectedOrder {
		t.Errorf("want %s got %s", expectedOrder, resultOrder)
	}
}
