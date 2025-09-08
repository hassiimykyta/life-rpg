package resp

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/hassiimykyta/life-rpg/apps/gateway/internal/dto"
)

func OK(w http.ResponseWriter, r *http.Request, data any, message string, codes ...int) {
	code := http.StatusOK
	if len(codes) > 0 {
		code = codes[0]
	}
	JSON(w, r, data, code, message)
}

func ERROR(w http.ResponseWriter, r *http.Request, message string, codes ...int) {
	code := http.StatusInternalServerError
	if len(codes) > 0 {
		code = codes[0]
	}
	JSON(w, r, nil, code, message)
}

func JSON(w http.ResponseWriter, r *http.Request, data any, code int, message string) {
	render.JSON(w, r, dto.BasicResponse{
		Code:    code,
		Data:    data,
		Message: message,
	})
}
