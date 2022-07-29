package main

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrAlreadyInitialized         = errors.New("parking lot is already initialized")
	ErrNotInitialized             = errors.New("parking lot has not been initialized")
	ErrFull                       = errors.New("parking lot is full")
	ErrSlotIsEmpty                = errors.New("slot is empty")
	ErrSlotIsNotEmpty             = errors.New("slot is not empty")
	ErrSlotDoesNotExist           = errors.New("slot does not exist")
	ErrRegistrationNumberNotFound = errors.New("registration number not found")
)

type Car struct {
	RegistrationNumber string
	Colour             string
}

type ParkingLot struct {
	initialized bool
	slots       []*Car
}

func NewParkingLot() *ParkingLot {
	return &ParkingLot{
		initialized: false,
		slots:       make([]*Car, 0),
	}
}

func (p *ParkingLot) Init(n int) error {
	if p.initialized {
		return ErrAlreadyInitialized
	}
	p.initialized = true
	p.slots = make([]*Car, n)
	return nil
}

func (p *ParkingLot) Park(car *Car) (int, error) {
	if !p.initialized {
		return -1, ErrNotInitialized
	}
	for i := range p.slots {
		if p.slots[i] == nil {
			return p.assignCarToSlot(car, i+1)
		}
	}
	return -1, ErrFull
}

func (p *ParkingLot) assignCarToSlot(car *Car, slot int) (int, error) {
	if !p.initialized {
		return -1, ErrNotInitialized
	}
	if slot <= 0 || slot > len(p.slots) {
		return -1, ErrSlotDoesNotExist
	}
	if p.slots[slot-1] != nil {
		return -1, ErrSlotIsNotEmpty
	}
	p.slots[slot-1] = car
	return slot, nil
}

func (p *ParkingLot) Leave(slot int) (*Car, error) {
	if !p.initialized {
		return nil, ErrNotInitialized
	}
	if slot <= 0 || slot > len(p.slots) {
		return nil, ErrSlotDoesNotExist
	}
	if p.slots[slot-1] == nil {
		return nil, ErrSlotIsEmpty
	}
	car := p.slots[slot-1]
	p.slots[slot-1] = nil
	return car, nil
}

func (p *ParkingLot) RegistrationNumbersForColour(colour string) []string {
	if !p.initialized {
		return nil
	}
	res := make([]string, 0)
	for _, car := range p.slots {
		if car != nil && car.Colour == colour {
			res = append(res, car.RegistrationNumber)
		}
	}
	return res
}

func (p *ParkingLot) SlotNumbersForColour(colour string) []int {
	if !p.initialized {
		return nil
	}
	res := make([]int, 0)
	for i, car := range p.slots {
		if car != nil && car.Colour == colour {
			res = append(res, i+1)
		}
	}
	return res
}

func (p *ParkingLot) SlotNumberForRegistrationNumber(reg string) (int, error) {
	if !p.initialized {
		return -1, ErrNotInitialized
	}
	for i, car := range p.slots {
		if car != nil && car.RegistrationNumber == reg {
			return i + 1, nil
		}
	}
	return -1, ErrRegistrationNumberNotFound
}

func (p *ParkingLot) Status() string {
	if !p.initialized {
		return ErrNotInitialized.Error()
	}

	entries := make([][]string, 0)
	indents := make([]int, 3)

	addEntries := func(entry []string) {
		for i, e := range entry {
			if len(e) > indents[i] {
				indents[i] = len(e)
			}
		}
		entries = append(entries, entry)
	}

	addEntries([]string{"Slot No. ", "Registration No ", "Colour"})
	for i, car := range p.slots {
		if car != nil {
			addEntries([]string{strconv.Itoa(i + 1), car.RegistrationNumber, car.Colour})
		}
	}

	lines := make([]string, 0)
	for _, entry := range entries {
		line := ""
		for i, e := range entry {
			line += e
			if i < len(entry)-1 {
				line += strings.Repeat(" ", indents[i]-len(e)+1)
			}
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}
