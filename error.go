package error_group

import (
	"strings"
	"sync"
)

type errorGroup struct {
	mutex  *sync.Mutex
	errors []error
}

func (eg *errorGroup) AddError(err error) {
	if err == nil {
		return
	}

	eg.mutex.Lock()
	defer eg.mutex.Unlock()

	eg.errors = append(eg.errors, err)
	return
}

func (eg *errorGroup) GetLast() error {
	eg.mutex.Lock()
	defer eg.mutex.Unlock()

	return eg.errors[len(eg.errors)-1]
}

func (eg *errorGroup) GetFirst() error {
	eg.mutex.Lock()
	defer eg.mutex.Unlock()

	return eg.errors[0]
}

func (eg *errorGroup) GetAll() []error {
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

	return sb.String()
}

//goland:noinspection GoExportedFuncWithUnexportedType
func NewErrorGroup() *errorGroup {
	errorMutex := sync.Mutex{}
	var multiErrors []error // TODO(Canavan): check if this is necessary

	return &errorGroup{
		mutex:  &errorMutex,
		errors: multiErrors,
	}
}
