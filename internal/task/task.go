package task

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type TaskRequest struct {
	Title string `json:"title"`
}
type TaskPatchRequest struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
}
