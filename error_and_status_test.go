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

func TestErrorStatusGroupMultipleThreads(t *testing.T) {
	esg := NewErrorStatusGroup()

	var sAndEWG sync.WaitGroup
	var eWG sync.WaitGroup
	var sWG sync.WaitGroup
	toAdd := 100000
	maxRoutines := 1000
	sandEGuard := make(chan struct{}, maxRoutines)
	eGuard := make(chan struct{}, maxRoutines)
	sGuard := make(chan struct{}, maxRoutines)

	for i := 0; i < toAdd; i++ {
		sandEGuard <- struct{}{}
		sAndEWG.Add(1)
		go func() {
			<-sandEGuard
			esg.AddStatusAndError(GenerateRandomNumber(), errors.New(generateRandomString(20)))
			sAndEWG.Done()
		}()
	}

	for i := 0; i < toAdd; i++ {
		eGuard <- struct{}{}
		eWG.Add(1)
		go func() {
			<-eGuard
			esg.AddError(errors.New(generateRandomString(20)))
			eWG.Done()
		}()
	}

	for i := 0; i < toAdd; i++ {
		sGuard <- struct{}{}
		sWG.Add(1)
		go func() {
			<-sGuard
			esg.AddStatus(GenerateRandomNumber())
			sWG.Done()
		}()
	}

	sWG.Wait()
	eWG.Wait()
	sAndEWG.Wait()

	t.Run("verify the number of status values is correct", func(t *testing.T) {
		assert.Equal(t, 2*toAdd, esg.LenErrors())
	})
	t.Run("verify the number of error values is correct", func(t *testing.T) {
		assert.Equal(t, 2*toAdd, esg.LenStatuses())
	})
}

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

	t.Run("verify AddError() correctly added all errors", func(t *testing.T) {
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

	t.Run("verify all status values were added via AddStatus() and check with LenStatuses()", func(t *testing.T) {
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

	t.Run("verify all values were added via AddStatusAndError() and check with LenErrors()", func(t *testing.T) {
		assert.Equal(t, esg.LenErrors(), numToAdd)
	})

	t.Run("verify all values were added via AddStatusAndError() and check with LenStatuses()", func(t *testing.T) {
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

	t.Run("verify All() returns the correct number of errors", func(t *testing.T) {
		assert.Equal(t, len(allErrors), 12)
	})
	t.Run("verify All() returns the correct number of statuses", func(t *testing.T) {
		assert.Equal(t, len(allStatuses), 12)
	})
	t.Run("verify slice returned by All() is not affected by more calls to AddError()", func(t *testing.T) {
		esg.AddError(errors.New(generateRandomString(10)))
		assert.Equal(t, len(allErrors), 12)
	})
	t.Run("verify slice returned by All() is not affected by more calls to AddStatus()", func(t *testing.T) {
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

	t.Run("verify output of Error() is correct", func(t *testing.T) {
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
	t.Run("verify Error() returns the empty string when there are no errors", func(t *testing.T) {
		other := NewErrorStatusGroup()
		assert.Equal(t, "", other.Error())
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

	t.Run("verify FirstError() returns the correct error value", func(t *testing.T) {
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

	t.Run("verify FirstStatus() returns the correct status value", func(t *testing.T) {
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

	t.Run("verify HighestStatus() returns the correct status value", func(t *testing.T) {
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

	t.Run("verify LastError() returns the correct error value", func(t *testing.T) {
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

	t.Run("verify LastStatus() returns the correct status value", func(t *testing.T) {
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

	t.Run("verify LowestStatus() returns the correct status value", func(t *testing.T) {
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

	t.Run("verify output of ToStatusAndError() is correct", func(t *testing.T) {
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

	t.Run("verify output of ToError() is correct", func(t *testing.T) {
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
	t.Run("verify ToError() returns nil when there are no errors", func(t *testing.T) {
		other := NewErrorStatusGroup()
		assert.Nil(t, other.ToError())
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
