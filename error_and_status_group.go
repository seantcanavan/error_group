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
		highestStatus: 200,
		lowestStatus:  200,
		statusesMutex: &statusMutex,
	}
}

// AddError adds an error to this error status group instance.
func (esg *errorStatusGroup) AddError(err error) {
	if err == nil {
		return
	}

	esg.errorsMutex.Lock()
	defer esg.errorsMutex.Unlock()

	esg.errors = append(esg.errors, err)
}

// AddStatus adds a status to this error status group instance. Status values should be
// 0 or greater. Negative status values will be ignored.
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

// AddStatusAndError adds an error and a status value to this error status group instance.
// Status values should be 0 or greater. Negative status values will be ignored.
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

// All returns two new slices - one containing every error value in this error status group instance.
// The other containing every status value in this error status group instance.
func (esg *errorStatusGroup) All() ([]int, []error) {
	esg.errorsMutex.Lock()
	esg.statusesMutex.Lock()
	defer esg.errorsMutex.Unlock()
	defer esg.statusesMutex.Unlock()

	dupErrors := make([]error, len(esg.errors))
	dupStatuses := make([]int, len(esg.statuses))

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		copy(dupErrors, esg.errors)
		wg.Done()
	}()

	go func() {
		copy(dupStatuses, esg.statuses)
		wg.Done()
	}()

	wg.Wait()

	return dupStatuses, dupErrors
}

// Error fulfills the builtin.Error interface and returns a concatenated string of all the errors in this
// error status group instance. It will also contain the highest and lowest status values encountered.
func (esg *errorStatusGroup) Error() string {
	esg.errorsMutex.Lock()
	esg.statusesMutex.Lock()
	defer esg.errorsMutex.Unlock()
	defer esg.statusesMutex.Unlock()

	if len(esg.errors) < 1 {
		return ""
	}

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

// FirstError returns the first error value saved to this error status group instance.
// Since this library is thread safe - the first error value saved is not deterministic
// if the library is used in a multithreaded environment.
func (esg *errorStatusGroup) FirstError() error {
	esg.errorsMutex.Lock()
	defer esg.errorsMutex.Unlock()

	return esg.errors[0]
}

// FirstStatus returns the first status value saved to this error status group instance.
// Since this library is thread safe - the first status value saved is not deterministic
// if the library is used in a multithreaded environment.
func (esg *errorStatusGroup) FirstStatus() int {
	esg.statusesMutex.Lock()
	defer esg.statusesMutex.Unlock()

	return esg.statuses[0]
}

// HighestStatus returns the current highest status value saved to this error status group instance. Subsequent
// calls to AddStatus or AddStatusAndError can cause the value returned here to no longer be accurate.
func (esg *errorStatusGroup) HighestStatus() int {
	esg.statusesMutex.Lock()
	defer esg.statusesMutex.Unlock()

	return esg.highestStatus
}

// LastError returns the last error value saved to this error status group instance. Subsequent calls
// to AddError or AddStatusAndError can cause the value returned here to no longer be the last.
func (esg *errorStatusGroup) LastError() error {
	esg.errorsMutex.Lock()
	defer esg.errorsMutex.Unlock()

	return esg.errors[len(esg.errors)-1]
}

// LastStatus returns the last status value saved to this error status group instance. Subsequent calls
// to AddStatus or AddStatusAndError can cause the value returned here to no longer be the last.
func (esg *errorStatusGroup) LastStatus() int {
	esg.statusesMutex.Lock()
	defer esg.statusesMutex.Unlock()

	return esg.statuses[len(esg.statuses)-1]
}

// LenErrors returns the (current) number of error values saved to this error status group instance.
// Subsequent calls to AddError or AddStatusAndError can cause the value returned here to no longer be accurate.
func (esg *errorStatusGroup) LenErrors() int {
	esg.errorsMutex.Lock()
	defer esg.errorsMutex.Unlock()

	return len(esg.errors)
}

// LenStatuses returns the (current) number of status values saved to this error status group instance.
// Subsequent calls to AddStatus or AddStatusAndError can cause the value returned here to no longer be accurate.
func (esg *errorStatusGroup) LenStatuses() int {
	esg.statusesMutex.Lock()
	defer esg.statusesMutex.Unlock()

	return len(esg.statuses)
}

// LowestStatus returns the current lowest status value saved to this error status group instance. Subsequent
// calls to AddStatus or AddStatusAndError can cause the value returned here to no longer be accurate.
func (esg *errorStatusGroup) LowestStatus() int {
	esg.statusesMutex.Lock()
	defer esg.statusesMutex.Unlock()

	return esg.lowestStatus
}

// ToStatusAndError returns the current highest status value in conjunction with a combined error value representing
// all the errors currently saved to this error status group. This should be used when execution is finished and a
// summary result is ready to be returned to the caller for processing.
func (esg *errorStatusGroup) ToStatusAndError() (int, error) {
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

	wg.Wait()

	return highestStatus, err
}

// ToError is a convenience function that converts the errors and statuses contained
// in this error status group into one single error. This is useful for returning the ErrorStatusGroup
// object instance as a single generic builtin.Error interface instance.
func (esg *errorStatusGroup) ToError() error {
	errMessage := esg.Error()
	if errMessage == "" {
		return nil
	}
	return errors.New(errMessage)
}
