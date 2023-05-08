package task_straregy

type SoloTaskProgression struct {
}

func (ctp *SoloTaskProgression) GetFirstStatusCode() (code string, error error) {
	return "in_progress", nil
}
