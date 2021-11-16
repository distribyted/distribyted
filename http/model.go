package http

type RouteAdd struct {
	Magnet string `json:"magnet" binding:"required"`
}

type Error struct {
	Error string `json:"error"`
}
