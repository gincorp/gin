package taskmanager

// TaskManager ...
// A TaskManager contains functions for dealing with RabbitMQ methods
// and APIs. Both `MasterManager` and `JobManager` implement the methods in
// this interface
type TaskManager interface {
	Consume(string) (map[string]interface{}, error)
}
