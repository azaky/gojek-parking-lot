package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func setupTest(n int) *ParkingLot {
	lot := NewParkingLot()
	lot.Init(n)
	return lot
}

func TestPark(t *testing.T) {
	lot := setupTest(3)

	num, err := lot.Park(&Car{"A", "Aqua"})
	require.NoError(t, err)
	require.Equal(t, 1, num)

	num, err = lot.Park(&Car{"B", "Black"})
	require.NoError(t, err)
	require.Equal(t, 2, num)

	num, err = lot.Park(&Car{"C", "Cyan"})
	require.NoError(t, err)
	require.Equal(t, 3, num)

	_, err = lot.Park(&Car{"C", "Cyan"})
	require.Equal(t, err, ErrFull)
}

func TestLeave(t *testing.T) {
	lot := setupTest(3)

	num, err := lot.Park(&Car{"A", "Aqua"})
	require.NoError(t, err)
	require.Equal(t, 1, num)

	car, err := lot.Leave(1)
	require.NoError(t, err)
	require.NotNil(t, car)
	require.Equal(t, "A", car.RegistrationNumber)
	require.Equal(t, "Aqua", car.Colour)

	num, err = lot.Park(&Car{"B", "Brown"})
	require.NoError(t, err)
	require.Equal(t, 1, num)
}

func TestColours(t *testing.T) {
	lot := setupTest(3)

	lot.Park(&Car{"A", "Aqua"})
	lot.Park(&Car{"B", "Black"})
	lot.Park(&Car{"C", "Black"})
	
	regs := lot.RegistrationNumbersForColour("Aqua")
	require.Equal(t, []string{"A"}, regs)

	regs = lot.RegistrationNumbersForColour("Black")
	require.Equal(t, []string{"B", "C"}, regs)

	regs = lot.RegistrationNumbersForColour("Cyan")
	require.Empty(t, regs)

	slots := lot.SlotNumbersForColour("Aqua")
	require.Equal(t, []int{1}, slots)

	slots = lot.SlotNumbersForColour("Black")
	require.Equal(t, []int{2, 3}, slots)

	slots = lot.SlotNumbersForColour("Cyan")
	require.Empty(t, slots)
}

func TestSlotNumberForRegistrationNumber(t *testing.T) {
	lot := setupTest(3)

	lot.Park(&Car{"A", "Aqua"})
	lot.Park(&Car{"B", "Black"})
	lot.Park(&Car{"C", "Black"})

	slot, err := lot.SlotNumberForRegistrationNumber("A")
	require.NoError(t, err)
	require.Equal(t, 1, slot)

	slot, err = lot.SlotNumberForRegistrationNumber("B")
	require.NoError(t, err)
	require.Equal(t, 2, slot)

	slot, err = lot.SlotNumberForRegistrationNumber("C")
	require.NoError(t, err)
	require.Equal(t, 3, slot)

	_, err = lot.SlotNumberForRegistrationNumber("D")
	require.Equal(t, ErrRegistrationNumberNotFound, err)
}