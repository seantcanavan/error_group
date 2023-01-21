package error_group

import (
	"strings"
	"sync"
)

type errorGroup struct {
	mutex  *sync.Mutex
	errors []error
}

func (eg *errorGroup) Add(err error) {
	if err == nil {
		return
	}

	eg.mutex.Lock()
	defer eg.mutex.Unlock()

	eg.errors = append(eg.errors, err)
	return
}

func (eg *errorGroup) All() []error {
	eg.mutex.Lock()
	defer eg.mutex.Unlock()

	return eg.errors
}

func (eg *errorGroup) Error() string {
	eg.mutex.Lock()
	defer eg.mutex.Unlock()

	sb := strings.Builder{}

	for _, currentError := range eg.errors {
		sb.WriteString(currentError.Error())
		sb.WriteString("\n")
	}

	return strings.TrimSuffix(sb.String(), "\n")
}

func (eg *errorGroup) First() error {
	eg.mutex.Lock()
	defer eg.mutex.Unlock()

	return eg.errors[0]
}

func (eg *errorGroup) Last() error {
	eg.mutex.Lock()
	defer eg.mutex.Unlock()

	return eg.errors[len(eg.errors)-1]
}

func (eg *errorGroup) Len() int {
	eg.mutex.Lock()
	defer eg.mutex.Unlock()

	return len(eg.errors)
}

//goland:noinspection GoExportedFuncWithUnexportedType
func NewErrorGroup() *errorGroup {
	errorMutex := sync.Mutex{}

	return &errorGroup{
		mutex: &errorMutex,
	}
}
