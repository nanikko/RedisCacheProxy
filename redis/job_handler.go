package redis

// JobHandler contains the actual work belong to a job
type JobHandler interface {
	Get(string) string
}
