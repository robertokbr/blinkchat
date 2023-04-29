package main

import (
	"fmt"
	"time"
)

type Match struct {
	clientOne string
	clientTwo string
}

func pushJob(job chan string, number int) {
	job <- fmt.Sprintf("user - %v", number)
	time.Sleep(1 * time.Second)
	pushJob(job, number+1)
}

func main() {
	jobs := make(chan string, 0)

	go pushJob(jobs, 1)

	for {
		job := <-jobs
		job2 := <-jobs
		fmt.Println(job)
		fmt.Println(job2)
		fmt.Println("-----")

		select {
		case <-time.After(5 * time.Second):
			fmt.Println("timeout")
		}
	}
}
