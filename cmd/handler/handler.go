package handler

type Handler struct{
	DB *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return Handler{DB: db}
}
