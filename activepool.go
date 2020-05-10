// Package to assist in careful use of resources like
// go-routines.
package workerpool

// TaskResult is the channel that holds the result of
// the task
type TaskResult chan error

// TaskRegister is the channel that is used to register
// functions before they are run in go-routines
type TaskRegister chan bool

// Task is the structure that defines the signature of
// functions that is required for ActivePool to work
type Task func(tskRslt TaskResult)

// ActivePool is called such as it keeps the
// number of active go-routines in check
type ActivePool struct {
	taskList []Task
	poolSize int
	tRlt     TaskResult
	tRgt     TaskRegister
}

// InitializeActivePool initializes an ActivePoolStruct based on
// the poolSize required and returns a pointer to it.
func InitializeActivePool(pS int) *ActivePool {
	return &ActivePool{poolSize: pS, tRlt: make(TaskResult, pS), tRgt: make(TaskRegister, pS)}
}

// AddWork adds Tasks to the task list that need to be completed
func (tp *ActivePool) AddWork(t Task) {
	tp.taskList = append(tp.taskList, t)
}

// Cleanup closes the channels and resets the taskList
// this ActivePool is no-longer to be used
func (tp *ActivePool) Cleanup() {
	close(tp.tRlt)
	close(tp.tRgt)
	tp.taskList = tp.taskList[:0]
}

// Reset resets the components for re-use
func (tp *ActivePool) Reset() {
	for len(tp.tRgt) > 0 {
		<-tp.tRgt
	}
	for len(tp.tRlt) > 0 {
		<-tp.tRlt
	}
	tp.taskList = tp.taskList[:0]
}

// DoWork does the actual heavy lifting and runs the functions
// with simple pooling logic. It returns the number of
// total results as well as the number of successes.
func (tp *ActivePool) DoWork() (int, int) {
	totalTaskCount := len(tp.taskList)
	i, resultCount, successCount := 0, 0, 0
	for {
		select {
		// any result that comes from go routine
		// run functions - should increment certain
		// counters and then pop from the registration
		// pool so that another function can register and
		// thus be run in a go-routine
		case err := <-tp.tRlt:
			resultCount++
			if err == nil {
				successCount++
			}
			// If the total number of results seen till
			// now are equal to the number of tasks
			// then we can stop checking for work
			if resultCount == totalTaskCount {
				tp.Reset()
				return resultCount, successCount
			}
			<-tp.tRgt
		// If the pool has capacity - then run
		// more go routines - till the pool
		// is full
		case tp.tRgt <- true:
			if i < totalTaskCount {
				go tp.taskList[i](tp.tRlt)
				i++
			}
		}
	}
}
