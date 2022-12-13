package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
)

type SignalPair struct {
	s1 []interface{}
	s2 []interface{}
}

type Signals []Signal

type Signal struct {
	value []interface{}
}

func (s Signal) equals(signal Signal) bool {
	_, s_smaller := compare(s.value, signal.value)
	_, signal_smaller := compare(signal.value, s.value)
	if !s_smaller && !signal_smaller {
		return true
	}
	return false
}

func (s Signals) Len() int {
	return len(s)
}

func (s Signals) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Signals) Less(i, j int) bool {
	return smallerThan(s[i].value, s[j].value)
}

func (p *SignalPair) ParseSignal(signal string) []interface{} {
	var parsed_signal []interface{}
	err := json.Unmarshal([]byte(signal), &parsed_signal)
	if err != nil {
		panic(err)
	}
	return parsed_signal
}

func (p *SignalPair) Init(signal1 string, signal2 string) {
	p.s1 = p.ParseSignal(signal1)
	p.s2 = p.ParseSignal(signal2)
}

func (p *SignalPair) check() bool {
	_, isValid := compare(p.s1, p.s2)
	return isValid
}

func smallerThan(s1 []interface{}, s2 []interface{}) bool {
	_, isValid := compare(s1, s2)
	return isValid
}

func compare(left interface{}, right interface{}, indent_size ...int) (isDone bool, isSmaller bool) {
	indent := "  "
	if len(indent_size) > 0 {
		indent = strings.Repeat("  ", indent_size[0])
	}
	fmt.Printf("%s- Compare %v vs %v\n", indent, left, right)
	left_type := reflect.TypeOf(left)
	right_type := reflect.TypeOf(right)
	//fmt.Println(left_type.Kind())
	//fmt.Println(right_type.Kind())
	if right_type.Kind().String() == "float64" && left_type.Kind().String() == "float64" {
		if left.(float64) == right.(float64) {
			return false, false
		} else if left.(float64) <= right.(float64) {
			fmt.Printf("%s  - Left side is smaller, so inputs are in the right order\n", indent)
			return true, true
		} else {
			fmt.Printf("%s  - Right side is smaller, so inputs are not in the right order\n", indent)
			return true, false
		}
	} else if right_type.Kind().String() == "slice" && left_type.Kind().String() == "slice" {
		status := true
		for l, left_value := range left.([]interface{}) {
			right_values := right.([]interface{})
			if l < len(right_values) {
				done, valid := compare(left_value, right_values[l], (len(indent)/2)+1)
				if done {
					return done, valid
				} else {
					status = status && valid
				}
			} else {
				fmt.Printf("%s  - Right side ran out of items, so inputs are not in the right order\n", indent)
				return true, false
			}
		}
		if len(right.([]interface{})) > len(left.([]interface{})) {
			fmt.Printf("%s  - Left side ran out of items, so inputs are in the right order\n", indent)
			return true, true
		}
		return false, status
	} else if right_type.Kind().String() == "float64" {
		var new_right []interface{}
		new_right = append(new_right, right)
		fmt.Printf("%s- Mixed types; convert right to %v and retry comparison\n", indent, new_right)
		return compare(left, new_right, (len(indent)/2)+1)
	} else if left_type.Kind().String() == "float64" {
		var new_left []interface{}
		new_left = append(new_left, left)
		fmt.Printf("%s- Mixed types; convert left to %v and retry comparison\n", indent, new_left)
		return compare(new_left, right, (len(indent)/2)+1)
	}
	return false, false
}

func findSignal(signals []Signal, searchSignal Signal) (index int) {
	for i, signal := range signals {
		fmt.Printf("Checking signal %v\n", signal)
		if signal.equals(searchSignal) {
			return i + 1 // to compensate different starting values
		}
	}
	return -1
}

func readFile(filename string) []SignalPair {
	fmt.Println("Reading file " + filename)

	input_file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer input_file.Close()

	fileScanner := bufio.NewScanner(input_file)
	fileScanner.Split(bufio.ScanLines)

	var signalPairs []SignalPair
	var signals []string
	for fileScanner.Scan() {
		line_text := fileScanner.Text()
		if line_text != "" {
			signals = append(signals, strings.TrimSpace(line_text))
		} else {
			if len(signals) != 2 {
				panic(errors.New("ERROR while parsing input. 2 signal lines expected!"))
			}
			signalPair := SignalPair{}
			signalPair.Init(signals[0], signals[1])
			signals = nil
			signalPairs = append(signalPairs, signalPair)
		}
	}
	// compute last signal
	signalPair := SignalPair{}
	signalPair.Init(signals[0], signals[1])
	signalPairs = append(signalPairs, signalPair)
	return signalPairs
}

func readFileAll(filename string) Signals {
	fmt.Println("Reading file " + filename)

	input_file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer input_file.Close()

	fileScanner := bufio.NewScanner(input_file)
	fileScanner.Split(bufio.ScanLines)

	var signals []Signal
	for fileScanner.Scan() {
		line_text := fileScanner.Text()
		if line_text != "" {
			var parsed_signal []interface{}
			err := json.Unmarshal([]byte(line_text), &parsed_signal)
			if err != nil {
				panic(err)
			}
			signal := Signal{value: parsed_signal}
			signals = append(signals, signal)
		}
	}

	return signals
}

func solve1(signals []SignalPair) {
	fmt.Println("[Solve1] Starting")
	indices_sum := 0
	for n, signal := range signals {
		idx := n + 1
		fmt.Printf("\n== Pair %d ==\n", idx)
		check_result := signal.check()
		fmt.Printf("Check Result: %t\n", check_result)
		if check_result {
			indices_sum += idx
		}
	}
	result := indices_sum
	fmt.Printf("[Solve1] Result: %d\n", result)
}

func solve2(signals []Signal, divider_packages []Signal) {
	fmt.Println("[Solve2] Starting")
	sort.Sort(Signals(signals))
	fmt.Println("[Solve2] Sorting done:")
	for _, signal := range signals {
		fmt.Println(signal)
	}
	fmt.Println(signals)
	fmt.Println("[Solve2] Begin search for dividers")
	indices := make([]int, 0)
	for _, divider := range divider_packages {
		fmt.Printf("Searching for divider %v\n", divider)
		indices = append(indices, findSignal(signals, divider))
	}
	fmt.Println(indices)
	result := indices[0] * indices[1]
	fmt.Printf("[Solve2] Result: %d\n", result)
}

func main() {
	input_signals := readFile("input.txt")
	fmt.Println(input_signals)

	solve1(input_signals)

	all_input_signals := readFileAll("input.txt")
	var divider_packages []Signal
	var Sig2 []interface{}
	var Sig6 []interface{}
	var val2 float64 = 2
	var val6 float64 = 6
	var val2i []interface{}
	var val6i []interface{}
	val2i = append(val2i, val2)
	val6i = append(val6i, val6)
	Sig2 = append(Sig2, val2i)
	Sig6 = append(Sig6, val6i)
	divider_packages = append(divider_packages, Signal{value: Sig2})
	divider_packages = append(divider_packages, Signal{value: Sig6})
	all_input_signals = append(all_input_signals, divider_packages...)
	fmt.Println(all_input_signals)

	solve2(all_input_signals, divider_packages)
}
