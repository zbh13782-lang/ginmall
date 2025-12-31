package llo

import "fmt"

var res int

func main() {
	var n int
	n,_=fmt.Scan()
	res = 0
	a := make([]int,n)
	process(a, 0, n)
	fmt.Println(res)
}

func valid(a []int, i, j int) bool {
	for k := 0; k < i; k++ {
		if a[k] == j || (j-a[k]) == (i-k) || (j-a[k]) == (k-i) {
			return false
		}
	}
	return true
}

func process(a []int, i, n int) {
	if i == n {
		res++
		return
	}

	for j := 0; j < n; j++ {
		if valid(a, i, j) {
			a[i] = j
			process(a, i+i, n)
		}
	}
}
