package util

import (
	"errors"
	"fmt"
)

type processor struct {
	f        func() (interface{}, error)
	respChan chan queueResult
}

type queueResult struct {
	resp interface{}
	err  error
}

type processingQueue struct {
	trigger chan processor
}

func NewProcessingQueue(size uint) *processingQueue {
	queue := processingQueue{
		trigger: make(chan processor),
	}
	for i := 0; i < int(size); i++ {
		go func() {
			for {
				func() {
					p := <-queue.trigger

					defer func() {
						if r := recover(); r != nil {
							fmt.Println(r)
							p.respChan <- queueResult{
								resp: nil,
								err:  errors.New("processing queue panic"),
							}
						}
					}()

					result, err := p.f()

					p.respChan <- queueResult{
						resp: result,
						err:  err,
					}
				}()
			}
		}()
	}
	return &queue
}

func (q *processingQueue) Trigger(f func() (interface{}, error)) (interface{}, error) {
	respC := make(chan queueResult)
	q.trigger <- processor{
		f:        f,
		respChan: respC,
	}

	resp := <-respC
	return resp.resp, resp.err
}
