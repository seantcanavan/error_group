package error_group

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/jgroeneveld/trial/assert"
	"math/big"
	"strconv"
	"strings"
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

	t.Run("verify number of errors is correct", func(t *testing.T) {
		assert.Equal(t, len(allErrors), 12)
	})
	t.Run("verify number of statuses is correct", func(t *testing.T) {
		assert.Equal(t, len(allStatuses), 12)
	})
	t.Run("verify errors returned is a new slice", func(t *testing.T) {
		esg.AddError(errors.New(generateRandomString(10)))
		assert.Equal(t, len(allErrors), 12)
	})
	t.Run("verify statuses returned is a new slice", func(t *testing.T) {
		esg.AddStatus(GenerateRandomNumber())
		assert.Equal(t, len(allStatuses), 12)
	})
}

func TestErrorStatusGroup_Error(t *testing.T) {
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

	t.Run("verify output of Error is correct", func(t *testing.T) {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("lowest status: [%d]", 1))
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("highest status: [%d]", 3))
		sb.WriteString("\n")

		var stringsToConcat []string

		stringsToConcat = append(stringsToConcat, firstMessage)

		for i := 0; i < numToAdd; i++ {
			stringsToConcat = append(stringsToConcat, middleMessage)
		}

		stringsToConcat = append(stringsToConcat, lastMessage)

		sb.WriteString(strings.Join(stringsToConcat, "\n"))

		assert.Equal(t, sb.String(), esg.Error())
	})

}

func TestErrorStatusGroup_FirstError(t *testing.T) {
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

	t.Run("verify first error returns the correct value", func(t *testing.T) {
		assert.Equal(t, esg.FirstError().Error(), firstMessage)
	})
}

func TestErrorStatusGroup_FirstStatus(t *testing.T) {
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

	t.Run("verify first status returns the correct value", func(t *testing.T) {
		assert.Equal(t, esg.FirstStatus(), 1)
	})
}

func TestErrorStatusGroup_HighestStatus(t *testing.T) {
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

	t.Run("verify highest status returns the correct value", func(t *testing.T) {
		assert.Equal(t, 3, esg.HighestStatus())
	})
}

func TestErrorStatusGroup_LastError(t *testing.T) {
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

	t.Run("verify last error returns the correct value", func(t *testing.T) {
		assert.Equal(t, lastMessage, esg.LastError().Error())
	})
}

func TestErrorStatusGroup_LastStatus(t *testing.T) {
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

	t.Run("verify last status returns the correct value", func(t *testing.T) {
		assert.Equal(t, 3, esg.LastStatus())
	})
}

func TestErrorStatusGroup_LowestStatus(t *testing.T) {
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

	t.Run("verify lowest status returns the correct value", func(t *testing.T) {
		assert.Equal(t, 1, esg.LowestStatus())
	})
}

func TestErrorStatusGroup_ToStatusAndError(t *testing.T) {
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

	t.Run("verify output of to status and error is correct", func(t *testing.T) {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("lowest status: [%d]", 1))
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("highest status: [%d]", 3))
		sb.WriteString("\n")

		var stringsToConcat []string

		stringsToConcat = append(stringsToConcat, firstMessage)

		for i := 0; i < numToAdd; i++ {
			stringsToConcat = append(stringsToConcat, middleMessage)
		}

		stringsToConcat = append(stringsToConcat, lastMessage)

		sb.WriteString(strings.Join(stringsToConcat, "\n"))

		statusCode, errVal := esg.ToStatusAndError()

		assert.Equal(t, 3, statusCode)
		assert.Equal(t, sb.String(), errVal.Error())
	})
}

func TestErrorStatusGroup_ToError(t *testing.T) {
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

	t.Run("verify output of to error is correct", func(t *testing.T) {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("lowest status: [%d]", 1))
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("highest status: [%d]", 3))
		sb.WriteString("\n")

		var stringsToConcat []string

		stringsToConcat = append(stringsToConcat, firstMessage)

		for i := 0; i < numToAdd; i++ {
			stringsToConcat = append(stringsToConcat, middleMessage)
		}

		stringsToConcat = append(stringsToConcat, lastMessage)

		sb.WriteString(strings.Join(stringsToConcat, "\n"))

		errString := esg.ToError().Error()

		assert.Equal(t, sb.String(), errString)
	})
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
