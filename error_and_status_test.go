package error_group

import (
	"crypto/rand"
	"errors"
	"github.com/jgroeneveld/trial/assert"
	"math/big"
	"strconv"
	"sync"
	"testing"
)

func TestErrorStatusGroup_AddError(t *testing.T) {
	esg := NewErrorStatusGroup()

	var wg sync.WaitGroup
	numToAdd := 100000
	maxRoutines := 1000
	guard := make(chan struct{}, maxRoutines)

	for i := 0; i < numToAdd; i++ {
		guard <- struct{}{}
		wg.Add(1)
		go func() {
			<-guard
			esg.AddError(errors.New(generateRandomString(20)))
			wg.Done()
		}()
	}

	wg.Wait()

	t.Run("verify all errors were added successfully", func(t *testing.T) {
		assert.Equal(t, esg.LenErrors(), numToAdd)
	})
}

func TestErrorStatusGroup_AddStatus(t *testing.T) {
	esg := NewErrorStatusGroup()

	var wg sync.WaitGroup
	numToAdd := 100000
	maxRoutines := 1000
	guard := make(chan struct{}, maxRoutines)

	for i := 0; i < numToAdd; i++ {
		guard <- struct{}{}
		wg.Add(1)
		go func() {
			<-guard
			esg.AddStatus(GenerateRandomNumber())
			wg.Done()
		}()
	}

	wg.Wait()

	t.Run("verify all statuses were added successfully", func(t *testing.T) {
		assert.Equal(t, esg.LenStatuses(), numToAdd)
	})
}

func TestErrorStatusGroup_AddStatusAndError(t *testing.T) {
	esg := NewErrorStatusGroup()

	var wg sync.WaitGroup
	numToAdd := 100000
	maxRoutines := 1000
	guard := make(chan struct{}, maxRoutines)

	for i := 0; i < numToAdd; i++ {
		guard <- struct{}{}
		wg.Add(1)
		go func() {
			<-guard
			esg.AddStatusAndError(GenerateRandomNumber(), errors.New(generateRandomString(10)))
			wg.Done()
		}()
	}

	wg.Wait()

	t.Run("verify all errors were added successfully", func(t *testing.T) {
		assert.Equal(t, esg.LenErrors(), numToAdd)
	})

	t.Run("verify all statuses were added successfully", func(t *testing.T) {
		assert.Equal(t, esg.LenStatuses(), numToAdd)
	})
}

func TestErrorStatusGroup_All(t *testing.T) {
	firstMessage := "first message"
	lastMessage := "last message"
	middleMessage := "middle message"

	esg := NewErrorStatusGroup()
	esg.AddStatusAndError(1, errors.New(firstMessage))

	numToAdd := 10
	for i := 0; i < numToAdd; i++ {
		esg.AddStatusAndError(2, errors.New(middleMessage))
	}

	esg.AddStatusAndError(3, errors.New(lastMessage))

	allStatuses, allErrors := esg.All()

}

func TestErrorStatusGroup_Error(t *testing.T) {
	assert.True(t, false)
}

func TestErrorStatusGroup_FirstError(t *testing.T) {
	assert.True(t, false)
}

func TestErrorStatusGroup_FirstStatus(t *testing.T) {
	assert.True(t, false)
}

func TestErrorStatusGroup_HighestStatus(t *testing.T) {
	assert.True(t, false)
}

func TestErrorStatusGroup_LastError(t *testing.T) {
	assert.True(t, false)
}

func TestErrorStatusGroup_LastStatus(t *testing.T) {
	assert.True(t, false)
}

func TestErrorStatusGroup_LowestStatus(t *testing.T) {
	assert.True(t, false)
}

func TestErrorStatusGroup_ToStatusAndError(t *testing.T) {
	assert.True(t, false)
}

func TestErrorStatusGroup_ToError(t *testing.T) {
	assert.True(t, false)
}

func GenerateRandomNumber() int {
	const letters = "123456789"
	numLength := 3
	ret := make([]byte, numLength)
	for i := 0; i < numLength; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return 0
		}
		ret[i] = letters[num.Int64()]
	}

	byteToInt, _ := strconv.Atoi(string(ret))
	return byteToInt
}
