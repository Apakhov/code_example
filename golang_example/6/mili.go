package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Set struct {
	parent *Set
	val    int
	depth  int
}

func (s *Set) initiate(val int) {
	s.parent = s
	s.val = val
	s.depth = 0
}

func (s *Set) findRoot() *Set {
	if s.parent == s {
		return s
	}
	s.parent = s.parent.findRoot()
	return s.parent
}

func (s *Set) unite(s2 *Set) {
	rootS1 := s.findRoot()
	rootS2 := s2.findRoot()
	if rootS1.depth < rootS2.depth {
		rootS1.parent = rootS2
	} else {
		rootS2.parent = rootS1
		if rootS1.depth == rootS2.depth && rootS1 != rootS2 {
			rootS1.depth++
		}
	}
}

type Mili struct {
	amStates        int
	inAlpSize       int
	startState      int
	transMatrix     [][]int
	outSignalMatrix [][]string
}

func (m *Mili) initiate(reader *bufio.Reader) {
	temp, _, _ := reader.ReadLine()
	m.amStates, _ = strconv.Atoi(strings.TrimSpace(string(temp)))
	temp, _, _ = reader.ReadLine()
	m.inAlpSize, _ = strconv.Atoi(strings.TrimSpace(string(temp)))
	temp, _, _ = reader.ReadLine()
	m.startState, _ = strconv.Atoi(strings.TrimSpace(string(temp)))
	m.transMatrix = make([][]int, m.amStates)
	for i := 0; i < m.amStates; i++ {
		m.transMatrix[i] = make([]int, m.inAlpSize)
		temp, _, _ = reader.ReadLine()
		tempSl := strings.Split(strings.TrimSpace(string(temp)), " ")
		for j, state := range tempSl {
			m.transMatrix[i][j], _ = strconv.Atoi(state)
		}
	}
	m.outSignalMatrix = make([][]string, m.amStates)
	for i := 0; i < m.amStates; i++ {
		m.outSignalMatrix[i] = make([]string, m.inAlpSize)
		temp, _, _ = reader.ReadLine()
		tempSl := strings.Split(strings.TrimSpace(string(temp)), " ")
		for j, signal := range tempSl {
			m.outSignalMatrix[i][j] = signal
		}
	}
}

func (m *Mili) canonize() {
	mapping := make([]int, m.amStates)
	for i := 0; i < m.amStates; i++ {
		mapping[i] = -1
	}
	m.recCanonizer(&mapping, m.startState, 0)
	var temp Mili
	temp.amStates = 0
	for _, state := range mapping {
		if state != -1 {
			temp.amStates++
		}
	}
	temp.inAlpSize = m.inAlpSize
	temp.startState = 0
	temp.outSignalMatrix = make([][]string, temp.amStates)
	temp.transMatrix = make([][]int, temp.amStates)
	for i, state := range mapping {
		if state == -1 {
			continue
		}
		temp.outSignalMatrix[state] = m.outSignalMatrix[i]
		temp.transMatrix[state] = make([]int, m.inAlpSize)
		for j := 0; j < m.inAlpSize; j++ {
			temp.transMatrix[state][j] = mapping[m.transMatrix[i][j]]
		}
	}
	*m = temp
}

func (m *Mili) recCanonizer(mapping *[]int, curState int, count int) int {
	(*mapping)[curState] = count
	count++
	for _, state := range m.transMatrix[curState] {
		if (*mapping)[state] == -1 {
			count = m.recCanonizer(mapping, state, count)
		}
	}
	return count
}

func (m *Mili) print() {
	fmt.Println(m.amStates)
	fmt.Println(m.inAlpSize)
	fmt.Println(m.startState)
	for _, x := range m.transMatrix {
		for _, y := range x {
			fmt.Print(y, " ")
		}
		fmt.Println()
	}
	for _, x := range m.outSignalMatrix {
		for _, y := range x {
			fmt.Print(y, " ")
		}
		fmt.Println()
	}
}

func (m *Mili) getGraphViz() string {
	buff := bytes.NewBuffer(make([]byte, 0))
	buff.WriteString(`digraph {
	rankdir = LR
	dummy [label = "", shape = none]
`)
	var a byte = 'a'
	for i := 0; i < m.amStates; i++ {
		buff.WriteString(`	`)
		buff.WriteString(strconv.Itoa(i))
		buff.WriteString(` [shape = circle]
`)
	}
	buff.WriteString(`	dummy -> `)
	buff.WriteString(strconv.Itoa(m.startState))
	buff.WriteString("\n")
	for i := 0; i < m.amStates; i++ {
		for j := 0; j < m.inAlpSize; j++ {
			buff.WriteString(`	`)
			buff.WriteString(strconv.Itoa(i))
			buff.WriteString(" -> ")
			buff.WriteString(strconv.Itoa(m.transMatrix[i][j]))
			buff.WriteString(` [label = "`)
			buff.WriteByte(a + byte(j))
			buff.WriteByte('(')
			buff.WriteString(m.outSignalMatrix[i][j])
			buff.WriteString(`)"]
`)

		}
	}
	buff.WriteByte('}')
	return buff.String()
}

func (m *Mili) split1() (n int, res []int) {
	n = m.amStates
	states := make([]Set, n)
	for i := 0; i < n; i++ {
		states[i].initiate(i)
	}
	for i := 0; i < m.amStates; i++ {
		for j := i + 1; j < m.amStates; j++ {
			if states[i].findRoot() != states[j].findRoot() {
				eq := true
				for z := 0; z < m.inAlpSize; z++ {
					if m.outSignalMatrix[i][z] != m.outSignalMatrix[j][z] {
						eq = false
						break
					}
				}
				if eq {
					states[i].unite(&states[j])
					n--
				}
			}
		}
	}
	res = make([]int, m.amStates)
	for i, state := range states {
		res[i] = state.findRoot().val
	}
	return
}

func (m *Mili) split(res *[]int) (n int) {
	n = m.amStates
	states := make([]Set, n)
	for i := 0; i < n; i++ {
		states[i].initiate(i)
	}
	for i := 0; i < m.amStates; i++ {
		for j := i + 1; j < m.amStates; j++ {
			if (*res)[i] == (*res)[j] && states[i].findRoot() != states[j].findRoot() {
				eq := true
				for z := 0; z < m.inAlpSize; z++ {
					w1 := m.transMatrix[i][z]
					w2 := m.transMatrix[j][z]
					if (*res)[w1] != (*res)[w2] {
						eq = false
						break
					}
				}
				if eq {
					states[i].unite(&states[j])
					n--
				}
			}
		}
	}
	for i, state := range states {
		(*res)[i] = state.findRoot().val
	}
	return
}

func (m *Mili) minimize() {
	m.canonize()
	n, res := m.split1()
	for {
		curn := m.split(&res)
		if curn == n {
			break
		}
		n = curn
	}

	count := make([]int, m.amStates)
	newStates := make([]int, 0)
	for _, x := range res {
		if count[x] == 0 {
			count[x] = 1
			newStates = append(newStates, x)
		}
	}
	mapping := make([]int, m.amStates)
	for i, x := range newStates {
		mapping[x] = i
	}

	var temp Mili
	temp.amStates = len(newStates)
	temp.inAlpSize = m.inAlpSize
	temp.startState = res[0]
	temp.outSignalMatrix = make([][]string, temp.amStates)
	temp.transMatrix = make([][]int, temp.amStates)
	for i := 0; i < temp.amStates; i++ {
		temp.outSignalMatrix[i] = m.outSignalMatrix[newStates[i]]
		temp.transMatrix[i] = make([]int, m.inAlpSize)

		for j := 0; j < m.inAlpSize; j++ {
			temp.transMatrix[i][j] = mapping[res[m.transMatrix[newStates[i]][j]]]
		}
	}

	temp.canonize()
	*m = temp
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var (
		automate1 Mili
		automate2 Mili
	)
	automate1.initiate(reader)
	reader.ReadLine()
	automate2.initiate(reader)
	automate2.minimize()
	automate1.minimize()
	if automate1.amStates == automate2.amStates &&
		automate1.inAlpSize == automate2.inAlpSize &&
		automate1.startState == automate2.startState {
		equal := true
	mainLoop:
		for i := 0; i < automate1.amStates; i++ {
			for j := 0; j < automate1.inAlpSize; j++ {
				if automate1.outSignalMatrix[i][j] != automate2.outSignalMatrix[i][j] ||
					automate1.transMatrix[i][j] != automate2.transMatrix[i][j] {
					equal = false
					break mainLoop
				}
			}
		}
		if equal {
			fmt.Println("EQUAL")
			return
		}
	}
	fmt.Println("NOT EQUAL")
}
