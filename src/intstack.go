package main

type IntStack []int

func (s *IntStack) Push(v int) {
	*s = append(*s, v)
}

func (s *IntStack) Pop() int {
	res := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return res
}

func (s *IntStack) Top() int {
	return (*s)[len(*s)-1]
}

func (s *IntStack) IsEmpty() bool {
	return len(*s) == 0
}
