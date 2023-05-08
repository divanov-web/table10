package task_straregy

type CommonTaskProgression struct {
}

func (ctp *CommonTaskProgression) GetFirstStatusCode() (code string, error error) {
	return "in_progress", nil
}
