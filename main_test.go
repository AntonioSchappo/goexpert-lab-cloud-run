package main

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Handle(t *testing.T) {
	t.Run("should return 422 when cep is invalid", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		r.SetPathValue("cep", "123456")
		Handle(w, r)
		if w.Code != 422 {
			t.Errorf("expected 422, got %d", w.Code)
		}
	})

	t.Run("should return 404 when cep is not found", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", nil)
		r.SetPathValue("cep", "12345678")
		w := httptest.NewRecorder()

		Handle(w, r)
		assert.Equal(t, w.Code, 404)
	})

	t.Run("should return 200 when everything is ok", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", nil)
		r.SetPathValue("cep", "01153000")
		w := httptest.NewRecorder()

		Handle(w, r)
		assert.Equal(t, w.Code, 200)
		assert.NotNil(t, w.Body)
	})
}

func Test_validCep(t *testing.T) {
	t.Run("should return true when cep is valid", func(t *testing.T) {
		assert.True(t, validCep("12345678"))
	})

	t.Run("should return false when cep is invalid", func(t *testing.T) {
		assert.False(t, validCep("1234567"))
	})
}
