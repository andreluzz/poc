package controllers

import (
	"fmt"
	"sync"
	"time"

	"github.com/andreluzz/poc/amqp"
)

//Job represents an running instance of the job definition
type Job struct {
	ID          string    `json:"id" sql:"id"`
	ServiceID   string    `json:"service_id" sql:"service_id"`
	Status      string    `json:"status" sql:"status"`
	Start       time.Time `json:"start_at" sql:"start_at"`
	Finish      time.Time `json:"finish_at" sql:"finish_at"`
	Tasks       []Task    `json:"tasks"`
	Execution   chan *Task
	Responses   chan *Task
	Instance    int
	Concurrency int
	WG          sync.WaitGroup
}

func (j *Job) run() {
	fmt.Printf("JOB: %s | Instance: %00d | processing > JOB Instance ID: %s\n", j.ServiceID, j.Instance, j.ID)
	//TODO: Check and wait until JOB Instance is in created status
	j.Status = statusProcessing
	//TODO: Update job instance on DB (status, service_id)
	//TODO: Read tasks from database
	mockTasks(&j.Tasks)

	fmt.Printf("Total tasks %d\n", len(j.Tasks))
	//TODO: Read Tasks params from database

	j.WG.Add(len(j.Tasks))

	j.defineTasksToExecute("", "", 0)
	go func() {
		for tsk := range j.Responses {
			fmt.Printf("    Task: %s | Status: %s\n", tsk.ID, tsk.Status)
			//TODO: check task status to deal with errors
			j.WG.Done()
			j.defineTasksToExecute(tsk.ID, tsk.ParentID, tsk.Sequence)
		}
	}()

	j.WG.Wait()

	fmt.Printf("JOB: %s | Instance: %00d | Completed  > JOB Instance ID: %s\n", j.ServiceID, j.Instance, j.ID)
}

func (j *Job) work() {
	fmt.Printf("JOB %d - Start Task Process \n", j.Instance)
	for tsk := range j.Execution {
		tsk.Run(j.Responses)
	}
}

//Process keep checkin channel to process job messages
func (j *Job) Process(jobs <-chan *amqp.Message) {
	fmt.Println("Start JOB Process")
	for i := 0; i < j.Concurrency; i++ {
		go j.work()
	}

	for msg := range jobs {
		fmt.Println("Get message from jobs channel")
		j.ID = msg.ID
		j.run()
	}
}

func (j *Job) defineTasksToExecute(id, parentID string, sequence int) {
	//check if sequence is completed
	sequenceCompleted := true
	for _, t := range j.Tasks {
		if t.ParentID == parentID && t.Sequence == sequence && (t.Status == statusProcessing || t.Status == statusCreated) {
			sequenceCompleted = false
		}
	}

	if sequenceCompleted {
		fmt.Println("Sequence completed")
		sequence++
	}

	for i, t := range j.Tasks {
		if t.ParentID == parentID && t.Sequence == sequence && t.Status == statusCreated {
			fmt.Printf("Push to channel execution -> Task %s\n", t.ID)
			j.Execution <- &j.Tasks[i]
		}
	}

	if id != "" {
		//Check if has childs to start executing
		for i, t := range j.Tasks {
			if t.ParentID == id && t.Sequence == 0 && t.Status == statusCreated {
				fmt.Printf("Push to channel execution -> Task %s\n", t.ID)
				j.Execution <- &j.Tasks[i]
			}
		}
	}

}

func (j *Job) end() {
	close(j.Responses)
}

func mockTasks(tasks *[]Task) {
	t := Task{
		Status:      statusCreated,
		Sequence:    0,
		ID:          "AAA",
		ExecTimeout: 2,
	}
	*tasks = append(*tasks, t)
	t = Task{
		Status:      statusCreated,
		Sequence:    0,
		ID:          "BBB",
		ExecTimeout: 5,
	}
	*tasks = append(*tasks, t)
	t = Task{
		Status:      statusCreated,
		Sequence:    0,
		ID:          "BB1",
		ParentID:    "BBB",
		ExecTimeout: 3,
	}
	*tasks = append(*tasks, t)
	t = Task{
		Status:      statusCreated,
		Sequence:    0,
		ID:          "BB2",
		ParentID:    "BBB",
		ExecTimeout: 2,
	}
	*tasks = append(*tasks, t)
	t = Task{
		Status:      statusCreated,
		Sequence:    1,
		ID:          "BB3",
		ParentID:    "BBB",
		ExecTimeout: 10,
	}
	*tasks = append(*tasks, t)
	t = Task{
		Status:      statusCreated,
		Sequence:    0,
		ID:          "B31",
		ParentID:    "BB3",
		ExecTimeout: 1,
	}
	*tasks = append(*tasks, t)
	t = Task{
		Status:      statusCreated,
		Sequence:    1,
		ID:          "B32",
		ParentID:    "BB3",
		ExecTimeout: 1,
	}
	*tasks = append(*tasks, t)
	t = Task{
		Status:      statusCreated,
		Sequence:    1,
		ID:          "CCC",
		ExecTimeout: 5,
	}
	*tasks = append(*tasks, t)
	t = Task{
		Status:      statusCreated,
		Sequence:    2,
		ID:          "DDD",
		ExecTimeout: 2,
	}
	*tasks = append(*tasks, t)
	t = Task{
		Status:      statusCreated,
		Sequence:    0,
		ID:          "DD1",
		ParentID:    "DDD",
		ExecTimeout: 5,
	}
	*tasks = append(*tasks, t)
	t = Task{
		Status:      statusCreated,
		Sequence:    1,
		ID:          "DD2",
		ParentID:    "DDD",
		ExecTimeout: 2,
	}
	*tasks = append(*tasks, t)
}
