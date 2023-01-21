package error_group

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type errorStatusGroup struct {
	errors        []error
	errorsMutex   *sync.Mutex
	highestStatus int
	lowestStatus  int
	statuses      []int
	statusesMutex *sync.Mutex
}

//goland:noinspection GoExportedFuncWithUnexportedType
func NewErrorStatusGroup() *errorStatusGroup {
	errorMutex := sync.Mutex{}
	statusMutex := sync.Mutex{}

	return &errorStatusGroup{
		errorsMutex:   &errorMutex,
		statusesMutex: &statusMutex,
	}
}

func (esg *errorStatusGroup) AddError(err error) {
	if err == nil {
		return
	}

	esg.errorsMutex.Lock()
	defer esg.errorsMutex.Unlock()

	esg.errors = append(esg.errors, err)
}

func (esg *errorStatusGroup) AddStatus(status int) {
	esg.statusesMutex.Lock()
	defer esg.statusesMutex.Unlock()

	if status < esg.lowestStatus {
		esg.lowestStatus = status
	}

	if status > esg.highestStatus {
		esg.highestStatus = status
	}

	esg.statuses = append(esg.statuses, status)
}

func (esg *errorStatusGroup) AddStatusAndError(status int, err error) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		esg.AddStatus(status)
		wg.Done()
	}()

	go func() {
		esg.AddError(err)
		wg.Done()
	}()

	wg.Wait()
}

func (esg *errorStatusGroup) Error() string {
	esg.errorsMutex.Lock()
	esg.statusesMutex.Lock()
	defer esg.errorsMutex.Unlock()
	defer esg.statusesMutex.Unlock()

	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("lowest status: [%d]", esg.lowestStatus))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("highest status: [%d]", esg.highestStatus))
	sb.WriteString("\n")

	for _, currentError := range esg.errors {
		sb.WriteString(currentError.Error())
		sb.WriteString("\n")
	}

	return strings.TrimSuffix(sb.String(), "\n")
}

func (esg *errorStatusGroup) FirstError() error {
	esg.errorsMutex.Lock()
	defer esg.errorsMutex.Unlock()

	return esg.errors[0]
}

func (esg *errorStatusGroup) FirstStatus() int {
	esg.statusesMutex.Lock()
	defer esg.statusesMutex.Unlock()

	return esg.statuses[0]
}

func (esg *errorStatusGroup) HighestStatus() int {
	esg.statusesMutex.Lock()
	defer esg.statusesMutex.Unlock()

	return esg.highestStatus
}

func (esg *errorStatusGroup) LastError() error {
	esg.errorsMutex.Lock()
	defer esg.errorsMutex.Unlock()

	return esg.errors[len(esg.errors)-1]
}

func (esg *errorStatusGroup) LastStatus() int {
	esg.statusesMutex.Lock()
	defer esg.statusesMutex.Unlock()

	return esg.statuses[len(esg.statuses)-1]
}

func (esg *errorStatusGroup) LowestStatus() int {
	esg.statusesMutex.Lock()
	defer esg.statusesMutex.Unlock()

	return esg.lowestStatus
}

func (esg *errorStatusGroup) ToError() error {
	return errors.New(esg.Error())
}

func (esg *errorStatusGroup) StatusAndToError() (int, error) {
	var wg sync.WaitGroup
	var highestStatus int
	var err error
	wg.Add(2)

	go func() {
		highestStatus = esg.HighestStatus()
		wg.Done()
	}()

	go func() {
		err = esg.ToError()
		wg.Done()
	}()

	return highestStatus, err
}
