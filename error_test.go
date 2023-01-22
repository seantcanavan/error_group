package error_group

import (
	"errors"
	"github.com/jgroeneveld/trial/assert"
	"math/rand"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestErrorGroup_Add(t *testing.T) {
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

	t.Run("verify Add() operations were successful and Len() returns correct value", func(t *testing.T) {
		assert.Equal(t, eg.Len(), numToAdd)
	})
}

func TestErrorGroup_All(t *testing.T) {
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

	t.Run("verify All() returns the correct first error message", func(t *testing.T) {
		assert.Equal(t, firstMessage, allErrors[0].Error())

	})
	t.Run("verify All() returns the correct middle messages", func(t *testing.T) {
		assert.Equal(t, lastMessage, allErrors[len(allErrors)-1].Error())
	})
	t.Run("verify All() returns the correct last message", func(t *testing.T) {
		for i := 1; i < len(allErrors)-1; i++ {
			assert.Equal(t, middleMessage, allErrors[i].Error())
		}
	})
	t.Run("verify All() returns a new slice that is not affected by Add()", func(t *testing.T) {
		eg.Add(errors.New(generateRandomString(10)))
		assert.Equal(t, len(allErrors), 12)
	})
}

func TestErrorGroup_Error(t *testing.T) {
	eg := NewErrorGroup()
	first := "first message"
	last := "last message"

	eg.Add(errors.New(first))
	eg.Add(errors.New(last))

	t.Run("verify Error() returns a correctly formatted error string", func(t *testing.T) {
		errString := eg.Error()
		assert.Equal(t, strings.Join([]string{first, last}, "\n"), errString)
	})
	t.Run("verify Error() returns the empty string when there are no errors", func(t *testing.T) {
		other := NewErrorGroup()
		assert.Equal(t, "", other.Error())
	})
}

func TestErrorGroup_First(t *testing.T) {
	eg := NewErrorGroup()
	first := "first message"
	last := "last message"

	eg.Add(errors.New(first))
	eg.Add(errors.New(last))

	t.Run("verify First() returns the correct error string", func(t *testing.T) {
		assert.Equal(t, first, eg.First().Error())
	})
}

func TestErrorGroup_Last(t *testing.T) {
	eg := NewErrorGroup()
	first := "first message"
	last := "last message"

	eg.Add(errors.New(first))
	eg.Add(errors.New(last))

	t.Run("verify Last() returns the correct error string", func(t *testing.T) {
		assert.Equal(t, last, eg.Last().Error())
	})
}

func TestErrorGroup_Len(t *testing.T) {
	eg := NewErrorGroup()
	toAdd := 10
	for i := 0; i < toAdd; i++ {
		eg.Add(errors.New(generateRandomString(10)))
	}

	t.Run("verify Len() returns the correct number of errors", func(t *testing.T) {
		assert.Equal(t, toAdd, eg.Len())
	})
}

func TestErrorGroup_ToError(t *testing.T) {
	eg := NewErrorGroup()

	first := "first message"
	last := "last message"

	eg.Add(errors.New(first))
	eg.Add(errors.New(last))

	t.Run("verify ToError() returns the correctly formatted error message", func(t *testing.T) {
		assert.True(t, reflect.DeepEqual(errors.New(eg.Error()), eg.ToError()))
	})
	t.Run("verify ToError() returns nil when there are no errors", func(t *testing.T) {
		other := NewErrorGroup()
		assert.Nil(t, other.ToError())
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
