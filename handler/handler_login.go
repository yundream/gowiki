package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Login ")
}
