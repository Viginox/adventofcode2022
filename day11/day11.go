package main

import (
    "fmt"
    "os"
    "bufio"
    "strings"
    "strconv"
    "encoding/json"
    "errors"
    "sort"
)

type MonkeyInfo struct {
    Name          string
    StartingItems []int64
    Operation     string
    OperationValue string
    Test          int64
    IfTrue        int64
    IfFalse       int64
    Inspections   int64
}

func parseMonkeyInfo (info_text []string) (MonkeyInfo) {
    monkeyInfo := MonkeyInfo{
        Inspections: 0,
    }
    fmt.Println("> Monkey info found")
    for _, line := range(info_text) {
        line_values := strings.Split(strings.TrimSpace(line), ":")
        par_name := strings.TrimSpace(line_values[0])
        par_value := strings.TrimSpace(line_values[1])
        //fmt.Println("Parsed " + par_name + " and " + par_value)
        switch (par_name) {
        case "Starting items":
            err := json.Unmarshal([]byte("[" + par_value + "]"), &monkeyInfo.StartingItems)
            if err != nil {
                panic(err)
            }
        case "Operation":
            // remove "new = " for better parsability later on
            operation_parts := strings.Split(par_value, " ")
            monkeyInfo.Operation = operation_parts[3]
            monkeyInfo.OperationValue = operation_parts[4]
        case "Test":
            // it is always "divisible"
            test_values := strings.Split(par_value, " ")
            test_val, _ := strconv.ParseInt(test_values[2], 10, 64)
            monkeyInfo.Test = int64(test_val)
        case "If true":
            // always starts with "throw to "
            true_values := strings.Split(par_value, " ")
            monkey_value, _ := strconv.ParseInt(true_values[3], 10, 64)
            monkeyInfo.IfTrue = int64(monkey_value)
        case "If false":
            // always starst with "throw to"
            false_values := strings.Split(par_value, " ")
            monkey_value, _ := strconv.ParseInt(false_values[3], 10, 64)
            monkeyInfo.IfFalse = int64(monkey_value)
        default:
            monkeyInfo.Name = par_name
        }
    }
    fmt.Println("> Monkey info parsed")
    return monkeyInfo
}

func MonkeyTurn(monkeyNr int, monkeys []MonkeyInfo, worry_reduction int64) {
    monkey := monkeys[monkeyNr]
    for _, item := range(monkey.StartingItems) {
        fmt.Printf("  Monkey inspects an item with a worry level of %d\n", item)
        monkeys[monkeyNr].Inspections += 1
        worry_level := item
        op_value := item
        monkeys[monkeyNr].StartingItems = monkeys[monkeyNr].StartingItems[1:]
        if (monkey.OperationValue != "old") {
            val,  err := strconv.ParseInt(monkey.OperationValue, 10, 64)
            if err != nil {
                panic(err)
            }
            op_value = int64(val)
        }
        switch (monkey.Operation) {
        case "+":
            worry_level += op_value
            fmt.Printf("    Worry level is increased by %d to %d\n", op_value, worry_level)
        case "-":
            worry_level -= op_value
            fmt.Printf("    Worry level is decreased by %d to %d\n", op_value, worry_level)
        case "*":
            worry_level *= op_value
            fmt.Printf("    Worry level is multiplied by %d to %d\n", op_value, worry_level)
        case "/":
            worry_level /= op_value
            fmt.Printf("    Worry level is divided by %d to %d\n", op_value, worry_level)
        default:
            panic(errors.New("Unknown operation " + monkey.Operation))
        }
        if worry_level < 0 {
            panic(errors.New("Worry level is smaller than 0!"))
        }
        worry_level = int64(worry_level / worry_reduction)
        fmt.Printf("    Monkey is now bored with item. Worry level is devided by %d to %d\n", worry_reduction, worry_level)
        // check value
        var divisor int64 = 1
        // we don't know the trace of the item
        for _, m := range(monkeys) {
            // gather all operations from all monkeys
            divisor *= m.Test
        }
        // reduce worry level to a "manageable" size
        worry_level = worry_level % divisor
        if (worry_level % monkey.Test == 0) {
            fmt.Printf("    Current worry level is divisible by %d\n", monkey.Test)
            monkeys[monkey.IfTrue].StartingItems = append(monkeys[monkey.IfTrue].StartingItems, worry_level)
            fmt.Printf("    Item with worry level %d is thrown to monkey %d\n", worry_level, monkey.IfTrue)
        } else {
            fmt.Printf("    Current worry level is not divisible by %d\n", monkey.Test)
            monkeys[monkey.IfFalse].StartingItems = append(monkeys[monkey.IfFalse].StartingItems, worry_level)
            fmt.Printf("    Item with worry level %d is thrown to monkey %d\n", worry_level, monkey.IfFalse)
        }
    }
}

func readFile(filename string) ([]MonkeyInfo) {
    fmt.Println("Reading file " + filename)
    monkeyInfo := make([]MonkeyInfo, 0)
    input_file, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    defer input_file.Close()
    
    fileScanner := bufio.NewScanner(input_file)
    fileScanner.Split(bufio.ScanLines)

    fmt.Println("Parsing monkey info from data")
    monkeyText := make([]string, 0)
    for fileScanner.Scan() {
        line_text := fileScanner.Text()
        if (line_text != "") {
            monkeyText = append(monkeyText, line_text)
        } else {
            monkeyInfo = append(monkeyInfo, parseMonkeyInfo(monkeyText))
            monkeyText = make([]string, 0)
        }
    }
    // parse last monkey
    monkeyInfo = append(monkeyInfo, parseMonkeyInfo(monkeyText))
    fmt.Println("Monkey info parsed successfully")
    return monkeyInfo
}

func solve1 (monkeyInfo []MonkeyInfo, Rounds int) {
    fmt.Println("Running computation for Solution 1: " + strconv.Itoa(Rounds) + " rounds")
    for round := 0; round < Rounds; round++ {
        fmt.Println("Round " + strconv.Itoa(round+1))
        for m, _ := range(monkeyInfo) {
            fmt.Println("Monkey " + strconv.Itoa(m) + ":")
            MonkeyTurn(m, monkeyInfo, 3)
        }
    }
    inspections := make([]int64, 0)
    for _, monkey := range(monkeyInfo) {
        fmt.Println(monkey.Name)
        fmt.Println(monkey.StartingItems)
        fmt.Println(monkey.Inspections)
        inspections = append(inspections, monkey.Inspections)
    }
    sort.Slice(inspections[:], func(i, j int) bool {
		return inspections[:][i] > inspections[:][j]
	})
    fmt.Println(inspections)
    result := inspections[0] * inspections[1]
    fmt.Printf("[Solve1] Top 2 monkeys did %d inspections\n", result)
}

func solve2(monkeyInfo []MonkeyInfo, Rounds int) {
    fmt.Println("Running computation for Solution 1: " + strconv.Itoa(Rounds) + " rounds")
    for round := 0; round < Rounds; round++ {
        fmt.Println("Round " + strconv.Itoa(round+1))
        for m, _ := range(monkeyInfo) {
            fmt.Println("Monkey " + strconv.Itoa(m) + ":")
            MonkeyTurn(m, monkeyInfo, 1)
        }
    }
    inspections := make([]int64, 0)
    for _, monkey := range(monkeyInfo) {
        fmt.Println(monkey.Name)
        fmt.Println(monkey.StartingItems)
        fmt.Println(monkey.Inspections)
        inspections = append(inspections, monkey.Inspections)
    }
    sort.Slice(inspections[:], func(i, j int) bool {
		return inspections[:][i] > inspections[:][j]
	})
    fmt.Println(inspections)
    result := inspections[0] * inspections[1]
    fmt.Printf("=== After %d Rounds ===\n", Rounds)
    fmt.Printf("[Solve2] Top 2 monkeys did %d inspections\n", result)
}

func main() {
    monkeyInfo := readFile("input.txt")
    monkeyInfo2 := readFile("input.txt")
    fmt.Println(monkeyInfo)
    solve1(monkeyInfo, 20)
    fmt.Println(monkeyInfo2)
    solve2(monkeyInfo2, 10000)
}
