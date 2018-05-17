package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

type stack []int

func newStack() stack {
	return make([]int, 0)
}

func (s *stack) push(x int) {
	*s = append(*s, x)
}

func (s *stack) empty() bool {
	return len(*s) == 0
}

func (s *stack) size() int {
	return len(*s)
}

func (s *stack) pop() (x int, err error) {
	if s.empty() {
		err = errors.New("no values in stack")
		return
	}
	x = (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return
}

func (s *stack) top() (x int, err error) {
	if s.empty() {
		err = errors.New("no values in stack")
		return
	}
	x = (*s)[len(*s)-1]
	return
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n'
}

func isNumber(r rune) bool {
	return '0' <= r && r <= '9'
}

func calc(in io.Reader, out io.Writer) error {
	stack := newStack()
	sequence := make([]byte, 100)
	for {
		n, errIn := in.Read(sequence)
		if errIn != nil {
			if errIn != io.EOF {
				return errIn
			}
			return nil
		}
		secRun := []rune(string(sequence[0:n]))
		length := len(secRun)
		i := 0
	loop:
		for {
			for i < length && isSpace(secRun[i]) {
				i++
			}
			switch {
			case i == length:
				break loop
			case isNumber(secRun[i]):
				num := 0
				for i < length && isNumber(secRun[i]) {
					num = num*10 + (int)(secRun[i]-'0')
					i++
				}
				stack.push(num)
			case secRun[i] == '+':
				if stack.size() >= 2 {
					x, _ := stack.pop()
					y, _ := stack.pop()
					stack.push(x + y)
					i++
				} else {
					return errors.New("not enough arguments for +")
				}
			case secRun[i] == '-':
				if stack.size() >= 2 {
					x, _ := stack.pop()
					y, _ := stack.pop()
					stack.push(y - x)
					i++
				} else {
					return errors.New("not enough arguments for -")
				}
			case secRun[i] == '*':
				if stack.size() >= 2 {
					x, _ := stack.pop()
					y, _ := stack.pop()
					stack.push(y * x)
					i++
				} else {
					return errors.New("not enough arguments for *")
				}
			case secRun[i] == '/':
				if stack.size() >= 2 {
					x, _ := stack.pop()
					if x == 0 {
						return errors.New("division by zero")
					}
					y, _ := stack.pop()
					stack.push(y / x)
					i++
				} else {
					return errors.New("not enough arguments for /")
				}
			case secRun[i] == '=':
				x, _ := stack.top()
				i++
				out.Write([]byte(strconv.Itoa(x)))
				out.Write([]byte("\n"))
			default:
				return errors.New("illigal symbol: " + (string)(secRun[i]))
			}
		}
	}
}

func main() {
	err := calc(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Println(err)
	}
}
