package error_group

import (
	"errors"
	"strings"
	"sync"
)

type errorGroup struct {
	mutex  *sync.Mutex
	errors []error
}

//goland:noinspection GoExportedFuncWithUnexportedType
func NewErrorGroup() *errorGroup {
	errorMutex := sync.Mutex{}

	return &errorGroup{
		mutex: &errorMutex,
	}
}

// Add adds an error to this error group instance.
func (eg *errorGroup) Add(err error) {
	if err == nil {
		return
	}

	eg.mutex.Lock()
	defer eg.mutex.Unlock()

	eg.errors = append(eg.errors, err)
	return
}

// All returns a new slice containing every error in this error group instance.
func (eg *errorGroup) All() []error {
	eg.mutex.Lock()
	defer eg.mutex.Unlock()

	duplicate := make([]error, len(eg.errors))

	copy(duplicate, eg.errors)

	return duplicate
}

// Error fulfills the builtin.Error interface and returns a concatenated string of all the errors in this error group instance.
func (eg *errorGroup) Error() string {
	eg.mutex.Lock()
	defer eg.mutex.Unlock()

	if len(eg.errors) == 0 {
		return ""
	}

	sb := strings.Builder{}

	for _, currentError := range eg.errors {
		sb.WriteString(currentError.Error())
		sb.WriteString("\n")
	}

	return strings.TrimSuffix(sb.String(), "\n")
}

// First returns the first error saved to this error group instance. Since this
// library is thread safe - the first error saved is not deterministic if the
// library is used in a multithreaded environment.
func (eg *errorGroup) First() error {
	eg.mutex.Lock()
	defer eg.mutex.Unlock()

	return eg.errors[0]
}

// Last returns the (current) last error saved to this error group instance.
// Subsequent calls to Add can cause the value returned here to no longer be the last.
func (eg *errorGroup) Last() error {
	eg.mutex.Lock()
	defer eg.mutex.Unlock()

	return eg.errors[len(eg.errors)-1]
}

// Len returns the (current) length or number of errors saved to this error instance.
// Subsequent calls to Add can cause the value returned here to no longer be accurate.
func (eg *errorGroup) Len() int {
	eg.mutex.Lock()
	defer eg.mutex.Unlock()

	return len(eg.errors)
}

// ToError is a convenience function that converts the errors contained in this
// error group into one single error. This is useful for returning the ErrorGroup
// object instance as a single generic builtin.Error interface instance.
func (eg *errorGroup) ToError() error {
	return errors.New(eg.Error())
}
