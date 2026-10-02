package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kirillat6/go-basis/internal/errs"
	"github.com/kirillat6/go-basis/internal/handler"
)

func TestGetId(t *testing.T) {
	tests := []struct{
		name string 
		path string
		id string
		wanted int
		wantedErr error
	}{
		{
			name: "Успешное получение ID",
			path: "/tasks/42",
			id: "42",
			wanted: 42,
			wantedErr: nil,
		}, 
		{
			name: "Некорректный ID",
			path: "/tasks/abc",
			id: "abc",
			wanted: 0,
			wantedErr: errs.ErrBadRequest,
		},
		{
			name: "ID отсутствует",
			path: "/tasks/",
			id: "",
			wanted: 0,
			wantedErr: errs.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T){
			r := httptest.NewRequest(http.MethodGet, tt.path, nil)
			r.SetPathValue("id", tt.id)
			id, err := handler.GetId(r)
			if tt.wantedErr == nil {
				if err != nil {
					t.Errorf("Ожидалась ошибка nil, получено: %v", err)
				}
			} else {
				if !errors.Is(err, tt.wantedErr) {
					t.Errorf("Ожидалась ошибка: %v, получено: %v", tt.wantedErr, err)
				}
			}
			if id != tt.wanted {
				t.Errorf("Получен не правильный id: %v, ожидался: %v", id, tt.wanted)
			}
		})
	}
}