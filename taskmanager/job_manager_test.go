package taskmanager

import (
	"fmt"
	"testing"
)

func TestNewJobManager(t *testing.T) {
	var jm interface{}

	t.Run("Initialises and returns a JobManager", func(t *testing.T) {
		jm = NewJobManager()

		switch jm.(type) {
		case JobManager:
		default:
			t.Errorf("NewJobManager() error = Received %T, expected JobManager", jm)
		}
	})

	list := jm.(JobManager).JobList
	t.Run("Initialises with a list of Jobs", func(t *testing.T) {
		if len(list) == 0 {
			t.Errorf("NewJobManager().JobList error = list is empty")
		}
	})

	// test default job functions
	for _, j := range []string{"post-to-web", "get-from-web", "log"} {
		t.Run(fmt.Sprintf("Initialises JobList and contains %s", j), func(t *testing.T) {
			_, ok := list[j]
			if !ok {
				t.Errorf("NewJobManager().JobList error = no such entry %s", j)
			}
		})
	}

}

func simpleJob(jn JobNotification) (output map[string]interface{}, err error) { return }
func TestJobManager_AddJob(t *testing.T) {
	jm := NewJobManager()
	jm.AddJob("simple-job", simpleJob)

	t.Run("adds jobs to JobList", func(t *testing.T) {
		_, ok := jm.JobList["simple-job"]
		if !ok {
			t.Errorf("NewJobManager().JobList error = no such entry 'simple-job'")
		}
	})
}

func TestJobManager_Consume(t *testing.T) {
	jm := NewJobManager()
	jm.AddJob("simple-job", simpleJob)

	o, err := jm.Consume(`
{
  "UUID": "abcd-efg-1234-foo",
  "Register": "name",
  "Type": "simple-job"
}
`)

	t.Run("consumes and runs a job", func(t *testing.T) {
		if err != nil {
			t.Errorf("Consume() error = %v", err)
		}
	})

	t.Run("returned job data has the correct UUID", func(t *testing.T) {
		if o["UUID"].(string) != "abcd-efg-1234-foo" {
			t.Errorf("Consume()[uuid] error = received %q, expected abcd-efg-1234-foo", o["UUID"].(string))
		}
	})

	t.Run("returned job data has the correct Register value", func(t *testing.T) {
		if o["Register"].(string) != "name" {
			t.Errorf("Consume()[register] error = received %q, expected name", o["Register"].(string))
		}
	})

	t.Run("returned job data has the correct Failed state", func(t *testing.T) {
		if o["Failed"].(bool) != false {
			t.Errorf("Consume()[failed] error = received %b, expected false", o["Failed"].(bool))
		}
	})

	t.Run("returned job data has the correct Data map", func(t *testing.T) {
		switch o["Data"].(type) {
		case map[string]interface{}:
		default:
			t.Errorf("Consume()[data] error = received %T, expected map[string]interface{}", o["Data"])
		}
	})

	t.Run("returned job data has a Start time", func(t *testing.T) {
		if o["Start"].(string) == "" {
			t.Errorf("Consume()[start] error = received %q, expected something useful", o["Start"].(string))
		}
	})

	t.Run("returned job data has an End time", func(t *testing.T) {
		if o["End"].(string) == "" {
			t.Errorf("Consume()[end] error = received %q, expected something useful", o["End"].(string))
		}
	})

	t.Run("returned job data has a Duration", func(t *testing.T) {
		if o["End"].(string) == " ms" {
			t.Errorf("Consume()[duration] error = received %q, expected something useful", o["Duraction"].(string))
		}
	})

}
