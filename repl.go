package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var whitespace = regexp.MustCompile("\\s+")

type REPL struct {
	lot *ParkingLot
	w   *bufio.Writer
}

func NewREPL(lot *ParkingLot) *REPL {
	return &REPL{
		lot: lot,
	}
}

func (r *REPL) write(message string, args ...interface{}) {
	line := fmt.Sprintf(message, args...)
	r.w.WriteString(line + "\n")
	r.w.Flush()
}

func (r *REPL) parseInput(input string) (string, []string) {
	tokens := whitespace.Split(input, -1)
	nonemptyTokens := make([]string, 0)
	for _, token := range tokens {
		if token != "" {
			nonemptyTokens = append(nonemptyTokens, token)
		}
	}
	if len(nonemptyTokens) == 0 {
		return "", nil
	}
	cmd := nonemptyTokens[0]
	args := nonemptyTokens[1:]
	return cmd, args
}

func (r *REPL) Run(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	r.w = bufio.NewWriter(out)

	for scanner.Scan() {
		input := scanner.Text()
		if input == "" {
			continue
		}
		cmd, args := r.parseInput(input)
		if cmd == "" {
			continue
		}

		switch cmd {
		case "create_parking_lot":
			r.handleCreateParkingLot(args)
		case "park":
			r.handlePark(args)
		case "leave":
			r.handleLeave(args)
		case "status":
			r.handleStatus(args)
		case "registration_numbers_for_cars_with_colour":
			r.handleRegistrationNumbersForColour(args)
		case "slot_numbers_for_cars_with_colour":
			r.handleSlotNumbersForColour(args)
		case "slot_number_for_registration_number":
			r.handleSlotNumberForRegistrationNumber(args)
		case "exit":
			return
		default:
			r.write("Invalid command %s", cmd)
		}
	}
}

func (r *REPL) handleCreateParkingLot(args []string) {
	if len(args) != 1 {
		r.write("Invalid number of args to create_parking_lot <num slots>")
		return
	}
	num, err := strconv.Atoi(args[0])
	if err != nil {
		r.write("Invalid args to create_parking_lot: %s: %s", args[0], err.Error())
		return
	}
	if num <= 0 {
		r.write("Invalid args to create_parking_lot: %s: number must be greater than 0", args[0])
		return
	}
	if err := r.lot.Init(num); err != nil {
		r.write("Error creating parking lot: %s", err.Error())
		return
	}
	r.write("Created a parking lot with %d slots", num)
}

func (r *REPL) handlePark(args []string) {
	if len(args) != 2 {
		r.write("Invalid number of args to park <registration number> <colour>")
		return
	}
	car := &Car{
		RegistrationNumber: args[0],
		Colour:             args[1],
	}
	slot, err := r.lot.Park(car)
	if err != nil {
		if errors.Is(err, ErrFull) {
			r.write("Sorry, parking lot is full")
		} else {
			r.write("Error parking the car: %s", err.Error())
		}
		return
	}
	r.write("Allocated slot number: %d", slot)
}

func (r *REPL) handleLeave(args []string) {
	if len(args) != 1 {
		r.write("Invalid number of args to leave <slot number>")
		return
	}
	num, err := strconv.Atoi(args[0])
	if err != nil {
		r.write("Invalid args to leave: %s: %s", args[0], err.Error())
		return
	}
	if _, err := r.lot.Leave(num); err != nil {
		if errors.Is(err, ErrSlotIsEmpty) {
			r.write("Cannot leave, slot %d is empty", num)
		} else if errors.Is(err, ErrSlotDoesNotExist) {
			r.write("Invalid slot number %d", num)
		} else {
			r.write("Error leaving slot %d: %s", num, err.Error())
		}
		return
	}
	r.write("Slot number %d is free", num)
}

func (r *REPL) handleStatus(args []string) {
	// ignore args for now.
	r.write(r.lot.Status())
}

func (r *REPL) handleRegistrationNumbersForColour(args []string) {
	if len(args) != 1 {
		r.write("Invalid number of args for registration_numbers_for_cars_with_colour <colour>")
		return
	}
	regs := r.lot.RegistrationNumbersForColour(args[0])
	r.write(strings.Join(regs, ", "))
}

func (r *REPL) handleSlotNumbersForColour(args []string) {
	if len(args) != 1 {
		r.write("Invalid number of args for registration_numbers_for_cars_with_colour <colour>")
		return
	}
	slots := r.lot.SlotNumbersForColour(args[0])
	slotsStr := make([]string, len(slots))
	for i, num := range slots {
		slotsStr[i] = strconv.Itoa(num)
	}
	r.write(strings.Join(slotsStr, ", "))
}

func (r *REPL) handleSlotNumberForRegistrationNumber(args []string) {
	if len(args) != 1 {
		r.write("Invalid number of args for slot_number_for_registration_number <registration number>")
		return
	}
	slot, err := r.lot.SlotNumberForRegistrationNumber(args[0])
	if err != nil {
		if errors.Is(err, ErrRegistrationNumberNotFound) {
			r.write("Not found")
		} else {
			r.write("Error finding registration number %s: %s", args[0], err.Error())
		}
		return
	}
	r.write("%d", slot)
}
