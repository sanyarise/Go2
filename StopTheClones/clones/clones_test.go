package clones

import (
	"bytes"
	"testing"
)

type MockCloneHandler struct {
	count int
}

func (m *MockCloneHandler) call(d Clone) {
	m.count++
}

var testclon = Clone{"TEST 1", "TEST 2"}

func TestWriter(t *testing.T) {
	buffer := bytes.Buffer{}
	write := GetWriter(&buffer)

	write(testclon)

	expected := "TEST 2\n"
	received := buffer.String()

	if received != expected {
		t.Errorf("expected %q, received %q", expected, received)
	}
}

func TestCSVWriter(t *testing.T) {
	buffer := bytes.Buffer{}
	write := GetCSVWriter(&buffer)

	write(testclon)

	expected := "\"TEST 1\",\"TEST 2\"\n"
	received := buffer.String()

	if received != expected {
		t.Errorf("expected %q, received %q", expected, received)
	}
}

func TestApplyFuncToChan(t *testing.T) {
	t.Run("Empty channel works", func(t *testing.T) {
		mockHandler := MockCloneHandler{0}
		emptyChannel := make(chan Clone)
		close(emptyChannel)
		ApplyFuncToChan(emptyChannel, mockHandler.call)

		if mockHandler.count != 0 {
			t.Error("Handler was called")
		}
	})

	t.Run("Handler is called for each items in channel", func(t *testing.T) {
		expect := 3
		mockHandler := MockCloneHandler{0}
		channel := make(chan Clone, expect)
		for i := 0; i < expect; i++ {
			channel <- Clone{"test", "test"}
		}
		close(channel)

		ApplyFuncToChan(channel, mockHandler.call)

		if mockHandler.count != expect {
			t.Error("Handler was called")
		}
	})

	t.Run("Handler received items in channel", func(t *testing.T) {
		expect := Clone{"test", "test"}
		mockHandler := func(d Clone) {
			if d != expect {
				t.Error("Wrong item received")
			}
		}
		channel := make(chan Clone, 1)
		channel <- expect
		close(channel)

		ApplyFuncToChan(channel, mockHandler)
	})
}
