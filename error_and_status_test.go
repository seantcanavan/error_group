package error_group

import (
	"errors"
	"github.com/jgroeneveld/trial/assert"
	"net/http"
	"strings"
	"testing"
)

func TestErrorAndAddError(t *testing.T) {
	multiError := NewErrorStatusGroup()

	firstError := "first error"
	multiError.AddError(errors.New(firstError))
	assert.True(t, strings.Contains(multiError.Error().Error(), firstError))

	secondError := "second error"
	multiError.AddError(errors.New(secondError))
	assert.True(t, strings.Contains(multiError.Error().Error(), secondError))
}

func TestStatusAndAddStatus(t *testing.T) {
	multiError := NewErrorStatusGroup()

	multiError.AddStatus(http.StatusOK)
	assert.Equal(t, multiError.Status(), http.StatusOK)
	multiError.AddStatus(http.StatusConflict)
	assert.Equal(t, multiError.Status(), http.StatusConflict)
	multiError.AddStatus(http.StatusCreated)
	assert.Equal(t, multiError.Status(), http.StatusConflict)
}

func TestAddStatusAndError(t *testing.T) {
	multiError := NewErrorStatusGroup()

	firstError := "first error"
	multiError.AddStatusAndError(http.StatusAlreadyReported, errors.New(firstError))

	assert.True(t, strings.Contains(multiError.Error().Error(), firstError))
	assert.Equal(t, multiError.Status(), http.StatusAlreadyReported)

	secondError := "second error"
	multiError.AddStatusAndError(http.StatusOK, errors.New(secondError))

	assert.True(t, strings.Contains(multiError.Error().Error(), secondError))
	assert.Equal(t, multiError.Status(), http.StatusAlreadyReported)

	multiError.AddStatusAndError(http.StatusInternalServerError, nil)

	assert.Equal(t, multiError.Status(), http.StatusInternalServerError)
}
