package main

type TaskManager interface {
    Consume(string)(map[string]interface{}, error)
    ShouldRespond()(bool)
}

func(m MasterManager)ShouldRespond() bool {
    return false
}

func(m JobManager)ShouldRespond() bool {
    return true
}
