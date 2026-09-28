package transfer

import (
	"context"
	"encoding/json"
	"net/http"
)

type Handler interface {
	CreatePlayer(ctx context.Context, name string, balance int64) error
}

type Handlercontr struct {
	uc Handler
}

func NewHandlerContr(uc Handler) *Handlercontr {
	return &Handlercontr{uc: uc}
}

func (h *Handlercontr) HandlerCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("ошибка метода"))
		return
	}
	var reg struct {
		Name    string `json:"name"`
		Balance int64  `json:"balance"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ошибка"))
		return
	}

	if err := h.uc.CreatePlayer(r.Context(), reg.Name, reg.Balance); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ошибка регистарции"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("создание аккаунта прошло успешно"))
}
