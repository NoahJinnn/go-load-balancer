/*
 Context before go into details:
	- There're 3 main subjects: Requester, Load Balancer, Worker
	- Requester send Requests to Load Balancer via workDeliver channel
	- Load Balancer knows Workers's task is done via doneNotifier channel
	- Request's results is sent back via their internal channel
	- Step:
	-> Requester generates and sends Requests to Load Balancer
	-> Load Balancer get Worker with lowest load and delegates the Request to it (via dispatch)
	-> Workers process Request, send back result for Request and notify Load Balancer it's done
	-> Load Balancer gets noti from the Worker and declines load for the Worker
*/
package main

import (
	"time"
)

func main() {
	myLB := &Balancer{
		pool: []*Worker{
			{
				make(chan Request),
				0,
				0,
			},
			{
				make(chan Request),
				0,
				0,
			},
		},
		doneNotifier: make(chan *Worker),
	}
	// Deliver req from requester to balancer
	workDeliver := make(chan Request)
	go SimulateRequester(workDeliver)
	go myLB.StartLB(workDeliver)
	time.Sleep(time.Duration(10000) * time.Second)
}
