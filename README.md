# Error Group / Error Status Group
Thread-safe store for error messages thrown by a set of go routines. Optionally stores HTTP status codes as well. Return the AND of all the errors or status codes when done.

When working with HTTP status codes you can request the highest status value in the store in the case where you need to return only a single status value to the caller.

## How to Build locally
1. `make build`

## How to Test locally
1. `make test`

## How to Use
1. `go get github.com/seantcanavan/error_group@latest`
2. `import github.com/seantcanavan/error_group`
3. Initialize a new ErrorGroup `eg := error_group.NewErrorGroup` or ErrorAndStatusGroup `esg := error_group.NewErrorAndStatusGroup`
4. Create a wait group to start working in parallel `var wg sync.WaitGroup`
5. Perform work in parallel `go func() { eg.AddError(parallelWork())}()`
6. Return a combined error when done `return eg.Error()`

## Sample ErrorStatusGroup example
Perform O(N) operations in O(1) time via go routines.
``` go
func GlobalSearch(ctx context.Context, gReq *GlobalSearchReq) (*GlobalSearchRes, int, error) {
	var admins []*admin.Admin
	var users []*user.User
	var learners []*learner.Learner
	var teachers []*teacher.Teacher

	esg := error_group.NewErrorStatusGroup()
	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		getAdmins, getAdminsHttpStatus, getAdminsErr := admin.Search(ctx, &admin.SearchReq{...})
		esg.AddStatusAndError(getAdminsHttpStatus, getAdminsErr)
		admins = getAdmins
		wg.Done()
	}()

	go func() {
		getUsers, getUsersHttpStatus, getUsersErr := user.Search(ctx, &user.SearchReq{...})
		esg.AddStatusAndError(getUsersHttpStatus, getUsersErr)
		users = getUsers
		wg.Done()
	}()

	go func() {
		getLearners, getLearnersHttpStatus, getLearnersErr := learner.Search(ctx, &learner.SearchReq{...})
		esg.AddStatusAndError(getLearnersHttpStatus, getLearnersErr)
		learners = getLearners
		wg.Done()
	}()

	go func() {
		getTeachers, getTeachersHttpStatus, getTeachersErr := teacher.Search(ctx, &teacher.SearchReq{...})
		esg.AddStatusAndError(getTeachersHttpStatus, getTeachersErr)
		teachers = getTeachers
		wg.Done()
	}()

	wg.Wait()

    // If no errors are encountered, esg.ToError() returns nil as expected to indiciate to the caller all is good
    // esg.HighestStatus will return 200 if no errors occurred and previous functions returned 200 OK
	return &GlobalSearchRes{
		Admins:      admins,
		Users:       users,
		Learners:    learners,
		Teachers:    teachers,
	}, esg.HighestStatus(), esg.ToError()
}
```

## All tests are passing
```
Sat Jan 21 11:47 PM error_group: make test
go test -v
=== RUN   TestErrorStatusGroupMultipleThreads
=== RUN   TestErrorStatusGroupMultipleThreads/verify_the_number_of_status_values_is_correct
=== RUN   TestErrorStatusGroupMultipleThreads/verify_the_number_of_error_values_is_correct
--- PASS: TestErrorStatusGroupMultipleThreads (1.79s)
    --- PASS: TestErrorStatusGroupMultipleThreads/verify_the_number_of_status_values_is_correct (0.00s)
    --- PASS: TestErrorStatusGroupMultipleThreads/verify_the_number_of_error_values_is_correct (0.00s)
=== RUN   TestErrorStatusGroup_AddError
=== RUN   TestErrorStatusGroup_AddError/verify_AddError()_correctly_added_all_errors
--- PASS: TestErrorStatusGroup_AddError (0.86s)
    --- PASS: TestErrorStatusGroup_AddError/verify_AddError()_correctly_added_all_errors (0.00s)
=== RUN   TestErrorStatusGroup_AddStatus
=== RUN   TestErrorStatusGroup_AddStatus/verify_all_status_values_were_added_via_AddStatus()_and_check_with_LenStatuses()
--- PASS: TestErrorStatusGroup_AddStatus (0.04s)
    --- PASS: TestErrorStatusGroup_AddStatus/verify_all_status_values_were_added_via_AddStatus()_and_check_with_LenStatuses() (0.00s)
=== RUN   TestErrorStatusGroup_AddStatusAndError
=== RUN   TestErrorStatusGroup_AddStatusAndError/verify_all_values_were_added_via_AddStatusAndError()_and_check_with_LenErrors()
=== RUN   TestErrorStatusGroup_AddStatusAndError/verify_all_values_were_added_via_AddStatusAndError()_and_check_with_LenStatuses()
--- PASS: TestErrorStatusGroup_AddStatusAndError (0.91s)
    --- PASS: TestErrorStatusGroup_AddStatusAndError/verify_all_values_were_added_via_AddStatusAndError()_and_check_with_LenErrors() (0.00s)
    --- PASS: TestErrorStatusGroup_AddStatusAndError/verify_all_values_were_added_via_AddStatusAndError()_and_check_with_LenStatuses() (0.00s)
=== RUN   TestErrorStatusGroup_All
=== RUN   TestErrorStatusGroup_All/verify_All()_returns_the_correct_number_of_errors
=== RUN   TestErrorStatusGroup_All/verify_All()_returns_the_correct_number_of_statuses
=== RUN   TestErrorStatusGroup_All/verify_slice_returned_by_All()_is_not_affected_by_more_calls_to_AddError()
=== RUN   TestErrorStatusGroup_All/verify_slice_returned_by_All()_is_not_affected_by_more_calls_to_AddStatus()
--- PASS: TestErrorStatusGroup_All (0.00s)
    --- PASS: TestErrorStatusGroup_All/verify_All()_returns_the_correct_number_of_errors (0.00s)
    --- PASS: TestErrorStatusGroup_All/verify_All()_returns_the_correct_number_of_statuses (0.00s)
    --- PASS: TestErrorStatusGroup_All/verify_slice_returned_by_All()_is_not_affected_by_more_calls_to_AddError() (0.00s)
    --- PASS: TestErrorStatusGroup_All/verify_slice_returned_by_All()_is_not_affected_by_more_calls_to_AddStatus() (0.00s)
=== RUN   TestErrorStatusGroup_Error
=== RUN   TestErrorStatusGroup_Error/verify_output_of_Error()_is_correct
=== RUN   TestErrorStatusGroup_Error/verify_Error()_returns_the_empty_string_when_there_are_no_errors
--- PASS: TestErrorStatusGroup_Error (0.00s)
    --- PASS: TestErrorStatusGroup_Error/verify_output_of_Error()_is_correct (0.00s)
    --- PASS: TestErrorStatusGroup_Error/verify_Error()_returns_the_empty_string_when_there_are_no_errors (0.00s)
=== RUN   TestErrorStatusGroup_FirstError
=== RUN   TestErrorStatusGroup_FirstError/verify_FirstError()_returns_the_correct_error_value
--- PASS: TestErrorStatusGroup_FirstError (0.00s)
    --- PASS: TestErrorStatusGroup_FirstError/verify_FirstError()_returns_the_correct_error_value (0.00s)
=== RUN   TestErrorStatusGroup_FirstStatus
=== RUN   TestErrorStatusGroup_FirstStatus/verify_FirstStatus()_returns_the_correct_status_value
--- PASS: TestErrorStatusGroup_FirstStatus (0.00s)
    --- PASS: TestErrorStatusGroup_FirstStatus/verify_FirstStatus()_returns_the_correct_status_value (0.00s)
=== RUN   TestErrorStatusGroup_HighestStatus
=== RUN   TestErrorStatusGroup_HighestStatus/verify_HighestStatus()_returns_the_correct_status_value
--- PASS: TestErrorStatusGroup_HighestStatus (0.00s)
    --- PASS: TestErrorStatusGroup_HighestStatus/verify_HighestStatus()_returns_the_correct_status_value (0.00s)
=== RUN   TestErrorStatusGroup_LastError
=== RUN   TestErrorStatusGroup_LastError/verify_LastError()_returns_the_correct_error_value
--- PASS: TestErrorStatusGroup_LastError (0.00s)
    --- PASS: TestErrorStatusGroup_LastError/verify_LastError()_returns_the_correct_error_value (0.00s)
=== RUN   TestErrorStatusGroup_LastStatus
=== RUN   TestErrorStatusGroup_LastStatus/verify_LastStatus()_returns_the_correct_status_value
--- PASS: TestErrorStatusGroup_LastStatus (0.00s)
    --- PASS: TestErrorStatusGroup_LastStatus/verify_LastStatus()_returns_the_correct_status_value (0.00s)
=== RUN   TestErrorStatusGroup_LowestStatus
=== RUN   TestErrorStatusGroup_LowestStatus/verify_LowestStatus()_returns_the_correct_status_value
--- PASS: TestErrorStatusGroup_LowestStatus (0.00s)
    --- PASS: TestErrorStatusGroup_LowestStatus/verify_LowestStatus()_returns_the_correct_status_value (0.00s)
=== RUN   TestErrorStatusGroup_ToStatusAndError
=== RUN   TestErrorStatusGroup_ToStatusAndError/verify_output_of_ToStatusAndError()_is_correct
--- PASS: TestErrorStatusGroup_ToStatusAndError (0.00s)
    --- PASS: TestErrorStatusGroup_ToStatusAndError/verify_output_of_ToStatusAndError()_is_correct (0.00s)
=== RUN   TestErrorStatusGroup_ToError
=== RUN   TestErrorStatusGroup_ToError/verify_output_of_ToError()_is_correct
=== RUN   TestErrorStatusGroup_ToError/verify_ToError()_returns_nil_when_there_are_no_errors
--- PASS: TestErrorStatusGroup_ToError (0.00s)
    --- PASS: TestErrorStatusGroup_ToError/verify_output_of_ToError()_is_correct (0.00s)
    --- PASS: TestErrorStatusGroup_ToError/verify_ToError()_returns_nil_when_there_are_no_errors (0.00s)
=== RUN   TestErrorGroup_Add
=== RUN   TestErrorGroup_Add/verify_Add()_operations_were_successful_and_Len()_returns_correct_value
--- PASS: TestErrorGroup_Add (0.86s)
    --- PASS: TestErrorGroup_Add/verify_Add()_operations_were_successful_and_Len()_returns_correct_value (0.00s)
=== RUN   TestErrorGroup_All
=== RUN   TestErrorGroup_All/verify_All()_returns_the_correct_first_error_message
=== RUN   TestErrorGroup_All/verify_All()_returns_the_correct_middle_messages
=== RUN   TestErrorGroup_All/verify_All()_returns_the_correct_last_message
=== RUN   TestErrorGroup_All/verify_All()_returns_a_new_slice_that_is_not_affected_by_Add()
--- PASS: TestErrorGroup_All (0.00s)
    --- PASS: TestErrorGroup_All/verify_All()_returns_the_correct_first_error_message (0.00s)
    --- PASS: TestErrorGroup_All/verify_All()_returns_the_correct_middle_messages (0.00s)
    --- PASS: TestErrorGroup_All/verify_All()_returns_the_correct_last_message (0.00s)
    --- PASS: TestErrorGroup_All/verify_All()_returns_a_new_slice_that_is_not_affected_by_Add() (0.00s)
=== RUN   TestErrorGroup_Error
=== RUN   TestErrorGroup_Error/verify_Error()_returns_a_correctly_formatted_error_string
=== RUN   TestErrorGroup_Error/verify_Error()_returns_the_empty_string_when_there_are_no_errors
--- PASS: TestErrorGroup_Error (0.00s)
    --- PASS: TestErrorGroup_Error/verify_Error()_returns_a_correctly_formatted_error_string (0.00s)
    --- PASS: TestErrorGroup_Error/verify_Error()_returns_the_empty_string_when_there_are_no_errors (0.00s)
=== RUN   TestErrorGroup_First
=== RUN   TestErrorGroup_First/verify_First()_returns_the_correct_error_string
--- PASS: TestErrorGroup_First (0.00s)
    --- PASS: TestErrorGroup_First/verify_First()_returns_the_correct_error_string (0.00s)
=== RUN   TestErrorGroup_Last
=== RUN   TestErrorGroup_Last/verify_Last()_returns_the_correct_error_string
--- PASS: TestErrorGroup_Last (0.00s)
    --- PASS: TestErrorGroup_Last/verify_Last()_returns_the_correct_error_string (0.00s)
=== RUN   TestErrorGroup_Len
=== RUN   TestErrorGroup_Len/verify_Len()_returns_the_correct_number_of_errors
--- PASS: TestErrorGroup_Len (0.00s)
    --- PASS: TestErrorGroup_Len/verify_Len()_returns_the_correct_number_of_errors (0.00s)
=== RUN   TestErrorGroup_ToError
=== RUN   TestErrorGroup_ToError/verify_ToError()_returns_the_correctly_formatted_error_message
=== RUN   TestErrorGroup_ToError/verify_ToError()_returns_nil_when_there_are_no_errors
--- PASS: TestErrorGroup_ToError (0.00s)
    --- PASS: TestErrorGroup_ToError/verify_ToError()_returns_the_correctly_formatted_error_message (0.00s)
    --- PASS: TestErrorGroup_ToError/verify_ToError()_returns_nil_when_there_are_no_errors (0.00s)
PASS
ok  	github.com/seantcanavan/error_group	4.488s
```
