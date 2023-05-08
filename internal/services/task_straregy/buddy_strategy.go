package task_straregy

type BuddyTaskProgression struct {
}

func (ctp *BuddyTaskProgression) GetFirstStatusCode() (code string, error error) {
	return "new", nil
}
