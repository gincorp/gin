package taskmanager

// TaskManager ...
// A TaskManager contains functions for dealing with RabbitMQ methods
// and APIs. Both `MasterManager` and `JobManager` implement the methods in
// this interface
type TaskManager interface {
    Consume(string) (map[string]interface{}, error)
    ShouldRespond() bool
    StartAPI()
}

// ShouldRespond ...
// Whether a node should respond to messages on the queue
func (m MasterManager) ShouldRespond() bool {
    return false
}

// ShouldRespond ...
// Whether a node should respond to messages on the queue
func (m JobManager) ShouldRespond() bool {
    return true
}
