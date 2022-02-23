package where

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// errorCallback always returns an error
func errorCallback(single, low, high *string) (*string, *string, *string, error) {
	return nil, nil, nil, errors.New("callback returns error")
}

func TestWalk_Empty(t *testing.T) {
	// Walk on an empty ValueSet should not call the callback function
	set := NewValueSet()
	newset, err := set.Walk(errorCallback)
	assert.NoError(t, err) // callback should not be called
	assert.Equal(t, newset, set)
}

func TestWalk_SingleError(t *testing.T) {
	// a callback error during singles processing should return an error from Walk
	set := NewValueSet()
	set.AddSingle("foo")
	newset, err := set.Walk(errorCallback)
	assert.Error(t, err)
	assert.Nil(t, newset)

}

func TestWalk_RangeError(t *testing.T) {
	// a callback error during range processing should return an error from Walk
	set := NewValueSet()
	set.AddRange("low", "high")
	newset, err := set.Walk(errorCallback)
	assert.Error(t, err)
	assert.Nil(t, newset)
}

// changeCallback is a callback function that returns certain changes depending on the values
// of the single or range given.
func changeCallback(single, low, high *string) (*string, *string, *string, error) {
	changed := "changed value"
	if single != nil {
		switch *single {
		case "delete me":
			single = nil
		case "change me":
			single = &changed
		}
	}
	if low != nil {
		switch *low {
		case "delete me":
			low = nil
			high = nil
		case "change me":
			low = &changed
			high = &changed
		}
	}
	return single, low, high, nil
}

func TestWalk_OneSingleNoChange(t *testing.T) {
	// test that Walk does not change a non-matching single value
	set := NewValueSet()
	set.AddSingle("a value")
	newset, err := set.Walk(changeCallback)
	assert.NoError(t, err)
	assert.Equal(t, set, newset)
}

func TestWalk_OneSingleChange(t *testing.T) {
	// test that Walk changes a matching single value
	set := NewValueSet()
	set.AddSingle("change me")
	newset, err := set.Walk(changeCallback)
	assert.NoError(t, err)
	want := NewValueSet()
	want.AddSingle("changed value")
	assert.Equal(t, newset, want)
}

func TestWalk_OneSingleDelete(t *testing.T) {
	// test that Walk deletes a matching single value
	set := NewValueSet()
	set.AddSingle("delete me")
	newset, err := set.Walk(changeCallback)
	assert.NoError(t, err)
	want := NewValueSet()
	assert.Equal(t, newset, want)
}

func TestWalk_TwoSinglesNoChange(t *testing.T) {
	// test that Walk does not change non-matching single values
	set := NewValueSet()
	set.AddSingle("a value")
	set.AddSingle("another value")
	newset, err := set.Walk(changeCallback)
	assert.NoError(t, err)
	assert.Equal(t, set, newset)
}

func TestWalk_TwoSinglesChange(t *testing.T) {
	// test that Walk changes only one matching single value
	set := NewValueSet()
	set.AddSingle("change me")
	set.AddSingle("another value")
	newset, err := set.Walk(changeCallback)
	assert.NoError(t, err)
	want := NewValueSet()
	want.AddSingle("changed value")
	want.AddSingle("another value")
	assert.Equal(t, newset, want)
}

func TestWalk_TwoSinglesDelete(t *testing.T) {
	// test that Walk deletes only one matching single value
	set := NewValueSet()
	set.AddSingle("delete me")
	set.AddSingle("another value")
	newset, err := set.Walk(changeCallback)
	assert.NoError(t, err)
	want := NewValueSet()
	want.AddSingle("another value")
	assert.Equal(t, newset, want)
}

func TestWalk_OneRangeNoChange(t *testing.T) {
	// test that Walk does not change a non-matching range
	set := NewValueSet()
	set.AddRange("low", "high")
	newset, err := set.Walk(changeCallback)
	assert.NoError(t, err)
	assert.Equal(t, set, newset)
}

func TestWalk_OneRangeChange(t *testing.T) {
	// test that Walk changes a single non-matching range
	set := NewValueSet()
	set.AddRange("change me", "and me")
	newset, err := set.Walk(changeCallback)
	assert.NoError(t, err)
	want := NewValueSet()
	want.AddRange("changed value", "changed value")
	assert.Equal(t, newset, want)
}

func TestWalk_OneRangeDelete(t *testing.T) {
	// test that Walk deletes a single matching range
	set := NewValueSet()
	set.AddRange("delete me", "and me")
	newset, err := set.Walk(changeCallback)
	assert.NoError(t, err)
	want := NewValueSet()
	assert.Equal(t, newset, want)
}

func TestWalk_TwoRangesNoChange(t *testing.T) {
	// test that Walk does not delete any non-matching ranges
	set := NewValueSet()
	set.AddRange("low1", "high1")
	set.AddRange("low2", "high2")
	newset, err := set.Walk(changeCallback)
	assert.NoError(t, err)
	assert.Equal(t, set, newset)
}

func TestWalk_TwoRangesChange(t *testing.T) {
	// test that Walk changes a single matching range
	set := NewValueSet()
	set.AddRange("change me", "and me")
	set.AddRange("low2", "high2")
	newset, err := set.Walk(changeCallback)
	assert.NoError(t, err)
	want := NewValueSet()
	want.AddRange("changed value", "changed value")
	want.AddRange("low2", "high2")
	assert.Equal(t, newset, want)
}

func TestWalk_TwoRangesDelete(t *testing.T) {
	// test that Walk deletes a single matching range
	set := NewValueSet()
	set.AddRange("delete me", "and me")
	set.AddRange("low2", "high2")
	newset, err := set.Walk(changeCallback)
	assert.NoError(t, err)
	want := NewValueSet()
	want.AddRange("low2", "high2")
	assert.Equal(t, newset, want)
}
