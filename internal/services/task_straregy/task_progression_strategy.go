package task_straregy

type TaskProgressionStrategy interface {
	GetFirstStatusCode() (code string, error error)
}
