package gotmx

type Task struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	NextID int    `json:"next_id"`
	PrevID int    `json:"prev_id"`
	Target string `json:"target"`
}
