package error_group

import (
	"errors"
	"github.com/jgroeneveld/trial/assert"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	eg := NewErrorGroup()

	var wg sync.WaitGroup
	numToAdd := 100000
	maxRoutines := 1000
	guard := make(chan struct{}, maxRoutines)

	for i := 0; i < numToAdd; i++ {
		guard <- struct{}{}
		wg.Add(1)
		go func() {
			<-guard
			eg.Add(errors.New(generateRandomString(20)))
			wg.Done()
		}()
	}

	wg.Wait()

	t.Run("verify all values were added successfully", func(t *testing.T) {
		assert.Equal(t, eg.Len(), numToAdd)
	})
}

func TestAll(t *testing.T) {
	firstMessage := "first message"
	lastMessage := "last message"
	middleMessage := "middle message"

	eg := NewErrorGroup()
	eg.Add(errors.New(firstMessage))

	numToAdd := 10
	for i := 0; i < numToAdd; i++ {
		eg.Add(errors.New(middleMessage))
	}

	eg.Add(errors.New(lastMessage))

	allErrors := eg.All()

	t.Run("verify first message is correct", func(t *testing.T) {
		assert.Equal(t, firstMessage, allErrors[0].Error())

	})
	t.Run("verify middle messages are correct", func(t *testing.T) {
		assert.Equal(t, lastMessage, allErrors[len(allErrors)-1].Error())
	})
	t.Run("verify last message is correct", func(t *testing.T) {
		for i := 1; i < len(allErrors)-1; i++ {
			assert.Equal(t, middleMessage, allErrors[i].Error())
		}
	})
}

func TestError(t *testing.T) {
	eg := NewErrorGroup()
	first := "first message"
	last := "last message"

	eg.Add(errors.New(first))
	eg.Add(errors.New(last))

	errString := eg.Error()
	assert.Equal(t, strings.Join([]string{first, last}, "\n"), errString)
}

func TestFirst(t *testing.T) {
	eg := NewErrorGroup()
	first := "first message"
	last := "last message"

	eg.Add(errors.New(first))
	eg.Add(errors.New(last))

	t.Run("verify first", func(t *testing.T) {
		assert.Equal(t, first, eg.First().Error())
	})
}

func TestLast(t *testing.T) {
	eg := NewErrorGroup()
	first := "first message"
	last := "last message"

	eg.Add(errors.New(first))
	eg.Add(errors.New(last))

	t.Run("verify last", func(t *testing.T) {
		assert.Equal(t, last, eg.Last().Error())
	})
}

func TestLen(t *testing.T) {
	eg := NewErrorGroup()
	toAdd := 10
	for i := 0; i < toAdd; i++ {
		eg.Add(errors.New(generateRandomString(10)))
	}

	t.Run("verify Len returns correct value", func(t *testing.T) {
		assert.Equal(t, toAdd, eg.Len())
	})
}

func generateRandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
