package main

import (
	"sort"
	"strconv"
	"sync"
)

func startJob(in, out chan interface{}, curjob job, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		if out != nil {
			close(out)
		}
	}()
	curjob(in, out)
}

func ExecutePipeline(jobs ...job) {
	var in chan interface{}
	wg := &sync.WaitGroup{}
	for _, curjob := range jobs[:len(jobs)-1] {
		wg.Add(1)
		out := make(chan interface{})
		go startJob(in, out, curjob, wg)
		in = out
	}
	wg.Add(1)
	go startJob(in, nil, jobs[len(jobs)-1], wg)
	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for num := range in {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			innnerWg := &sync.WaitGroup{}
			var first string
			var second string
			innnerWg.Add(2)
			go func() {
				defer innnerWg.Done()
				first = DataSignerCrc32(s)
			}()
			go func() {
				defer innnerWg.Done()
				second = DataSignerCrc32(DataSignerMd5help(s))
			}()
			innnerWg.Wait()
			out <- first + "~" + second
		}(strconv.Itoa(num.(int)))
	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for inpt := range in {
		hashes := make([]string, 6)
		wg.Add(1)
		go func(inpt interface{}) {
			defer wg.Done()
			innerWg := &sync.WaitGroup{}
			for x := 0; x < 6; x++ {
				innerWg.Add(1)
				go func(x int) {
					defer innerWg.Done()
					hashes[x] = DataSignerCrc32(strconv.Itoa(x) + inpt.(string))
				}(x)
			}
			innerWg.Wait()
			var res string
			for _, hash := range hashes {
				res += hash
			}
			out <- res
		}(inpt)
	}
	wg.Wait()
}

var crc32blocker = make(chan struct{}, 1)

func DataSignerMd5help(s string) string {
	crc32blocker <- struct{}{}
	out := DataSignerMd5(s)
	<-crc32blocker
	return out
}

func CombineResults(in, out chan interface{}) {
	hashes := make([]string, 0)
	for x := range in {
		hashes = append(hashes, x.(string))
	}
	sort.Strings(hashes)
	s := hashes[0]
	for _, x := range hashes[1:] {
		s += "_" + x
	}
	out <- s
}
