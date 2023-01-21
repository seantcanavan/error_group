package error_group

import (
	"errors"
	"strings"
	"sync"
)

type errorStatusGroup struct {
	errorMutex    *sync.Mutex
	firstError    error
	firstStatus   int
	highestStatus int
	lastError     error
	lastStatus    int
	lowestStatus  int
	multiErrors   []error
	multiStatuses []int
	statusMutex   *sync.Mutex
}

//goland:noinspection GoExportedFuncWithUnexportedType
func NewErrorStatusGroup() *errorStatusGroup {
	errorMutex := sync.Mutex{}
	statusMutex := sync.Mutex{}
	var multiErrors []error
	var multiStatuses []int

	return &errorStatusGroup{
		errorMutex:    &errorMutex,
		multiErrors:   multiErrors,
		multiStatuses: multiStatuses,
		statusMutex:   &statusMutex,
	}
}

func (m *errorStatusGroup) AddError(err error) {
	if err == nil {
		return
	}

	m.errorMutex.Lock()
	defer m.errorMutex.Unlock()

	m.multiErrors = append(m.multiErrors, err)
	return
}

func (m *errorStatusGroup) AddStatus(httpStatus int) {
	// don't take up the lock for a 0 value or 200 value - they're both default / okay
	if httpStatus == 0 || httpStatus == 200 {
		return
	}

	m.statusMutex.Lock()
	defer m.statusMutex.Unlock()

	m.multiStatuses = append(m.multiStatuses, httpStatus)
	return
}

func (m *errorStatusGroup) AddStatusAndError(httpStatus int, err error) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		m.AddStatus(httpStatus)
		wg.Done()
	}()

	go func() {
		m.AddError(err)
		wg.Done()
	}()

	wg.Wait()
	return
}

func (m *errorStatusGroup) Error() error {
	m.errorMutex.Lock()
	defer m.errorMutex.Unlock()

	if len(m.multiErrors) == 0 {
		return nil
	}

	sb := strings.Builder{}

	for _, currentError := range m.multiErrors {
		sb.WriteString(currentError.Error())
		sb.WriteString("\n")
	}

	return errors.New(sb.String())
}

func (m *errorStatusGroup) Status() int {
	m.statusMutex.Lock()
	defer m.statusMutex.Unlock()

	high := 200 // if no errors happened, then 200 happened!

	for _, currentStatus := range m.multiStatuses {
		if currentStatus > high {
			high = currentStatus
		}
	}

	return high
}

func (m *errorStatusGroup) StatusAndError() (int, error) {
	return m.Status(), m.Error()
}
