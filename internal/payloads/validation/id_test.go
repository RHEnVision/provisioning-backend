package validation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDigitsOnlySucceeds(t *testing.T) {
	id1 := "1111"
	err := DigitsOnly(id1)
	require.NoError(t, err)

	id2 := "1"
	err = DigitsOnly(id2)
	require.NoError(t, err)

	id3 := "0"
	err = DigitsOnly(id3)
	require.NoError(t, err)

	id4 := "-1111"
	err = DigitsOnly(id4)
	require.NoError(t, err)
}

func TestDigitsOnlyWithSymbolsFails(t *testing.T) {
	id1 := "۹"
	err := DigitsOnly(id1)
	if err == nil {
		t.Errorf("no error for invalid input: %s", id1)
	}

	id2 := "aa"
	err = DigitsOnly(id2)
	if err == nil {
		t.Errorf("no error for invalid input: %s", id2)
	}

	id3 := "7a"
	err = DigitsOnly(id3)
	if err == nil {
		t.Errorf("no error for invalid input: %s", id3)
	}

	id4 := "a7"
	err = DigitsOnly(id4)
	if err == nil {
		t.Errorf("no error for invalid input: %s", id4)
	}

	id5 := "7a7"
	err = DigitsOnly(id5)
	if err == nil {
		t.Errorf("no error for invalid input: %s", id5)
	}

	id6 := "zj{{=9243*9806}}zj"
	err = DigitsOnly(id6)
	if err == nil {
		t.Errorf("no error for invalid input: %s", id6)
	}

	id7 := "123.456"
	err = DigitsOnly(id7)
	if err == nil {
		t.Errorf("no error for invalid input: %s", id7)
	}

	id8 := "123⁰456"
	err = DigitsOnly(id8)
	if err == nil {
		t.Errorf("no error for invalid input: %s", id8)
	}

	id9 := "10^10"
	err = DigitsOnly(id9)
	if err == nil {
		t.Errorf("no error for invalid input: %s", id9)
	}

	id10 := "1_1"
	err = DigitsOnly(id10)
	if err == nil {
		t.Errorf("no error for invalid input: %s", id10)
	}

	id11 := "170141183460469231731687303715884105727"
	err = DigitsOnly(id11)
	if err == nil {
		t.Errorf("no error for invalid input: %s", id11)
	}
}
