package workerpool

type TaskResult chan error
type TaskRegister chan bool
type Task func(tskRslt TaskResult)
type ActivePool struct {
    taskList []Task
    poolSize int
    tRlt TaskResult
    tRgt TaskRegister
}

func InitializeActivePool(pS int) *ActivePool {
    return &ActivePool{poolSize: pS, tRlt: make(TaskResult, pS), tRgt: make(TaskRegister, pS) }
}

func (tp *ActivePool) AddWork(t Task) {
    tp.taskList = append(tp.taskList, t)
}

func (tp *ActivePool) DoWork() (int, int) {
    totalTaskCount := len(tp.taskList)
    i, resultCount, successCount := 0, 0, 0
    for ;; {
        select {
        // any result that comes from go routine
        // run functions - should increment certain
        // counters and then pop from the registration
        // pool so that another function can register and
        // thus be run in a go-routine
        case err := <- tp.tRlt:
            resultCount++
            if err == nil {
                successCount++
            }
            // If the total number of results seen till
            // now are equal to the number of tasks
            // then we can stop checking for work
            if resultCount == totalTaskCount {
                return resultCount, successCount
            }
            <- tp.tRgt
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