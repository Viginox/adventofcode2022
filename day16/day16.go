package main

import (
    "fmt"
    "bufio"
    "os"
    "strings"
    "strconv"
    "io"
    "regexp"
    "errors"
    "io/ioutil"
    "math"
    "sort"
)

type Valves []Valve

type Valve struct {
    Name        string
    FlowRate    int
    ConnectionNames []string
    Connections []*Valve
    IsOpen      bool
    OpenedAt    int
    Weight      int
}

func (v *Valve) Open(time int) {
    v.IsOpen = true
    v.OpenedAt = time
}

func (v *Valve) Close() {
    v.IsOpen = false
    v.OpenedAt = 0
}

func (v *Valve) GetTimeOpen(current_time int) (time_open int) {
    time_open = current_time - v.OpenedAt
    if time_open < 0 {
        return 0
    }
    return time_open
}

func (v *Valve) isOpen() (state bool) {
    return v.IsOpen
}

func (v *Valve) GetReleasedPressure(time int) (pressure int) {
    if !v.IsOpen {
        return 0
    }
    return v.GetTimeOpen(time) * v.FlowRate
}

func (v *Valve) GetDistanceTo(valve *Valve, checked_valves ...*Valve) (int) {
    checked_valves = append(checked_valves, v)
    if v == valve {
        return 0
    }
    minDist := math.MaxInt
    for val := 0; val < len(v.Connections); val++ {
        next_valve := v.Connections[val]
        already_checked := false
        for c := 0; c < len(checked_valves); c++ {
            if checked_valves[c] == next_valve {
                already_checked = true
            }
        }
        if already_checked {
            continue
        }
        dist := next_valve.GetDistanceTo(valve, checked_valves...) + 1
        if dist < minDist {
            minDist = dist
        }
    }
    if minDist == math.MaxInt {
        minDist = 0
    }
    return minDist
}

func readFile(r io.Reader) (input_lines []string) {
    scanner := bufio.NewScanner(r)
    scanner.Split(bufio.ScanLines)

    for scanner.Scan() {
        input_lines = append(input_lines, strings.TrimSpace(scanner.Text()))
    }

    return input_lines
}

func parseInput(input_lines []string) (valves []Valve) {
    fmt.Println("Reading input file from stdin")
    parse_re := regexp.MustCompile(`^Valve (?P<name>\w{2}) .*?=(?P<flow_rate>\d+); tunnel[s]? lead[s]? to valve[s]? (?P<valves>.*)$`)
    for _, line := range(input_lines) {
        line_match := parse_re.FindStringSubmatch(line)
        if line_match == nil {
            panic("Could not parse input line " + line)
        }
        flow_rate, err := strconv.ParseInt(line_match[2], 10, 0)
        if err != nil {
            panic(err)
        }
        valve := Valve{
            Name: line_match[1],
            FlowRate: int(flow_rate),
            ConnectionNames: strings.Split(line_match[3], ", "),
            IsOpen: false,
            OpenedAt: 0,
            Weight: -1,
        }
        valves = append(valves, valve)
    }
    return valves
}

func findConnections(valves []Valve) {
    fmt.Println("Looking for connections ...")
    for v := 0; v < len(valves); v++ {
        for _, targetValve := range(valves[v].ConnectionNames) {
            for s, searchValve := range(valves) {
                if searchValve.Name == targetValve {
                    valves[v].Connections = append(valves[v].Connections, &valves[s])
                    break
                }
            }
        }
        if len(valves[v].ConnectionNames) != len(valves[v].Connections) {
            panic(errors.New("Could not find all required connections for valve!"))
        }
    }
}

func calculateWeights(valve *Valve, weight int) {
    fmt.Printf("Calculating weight for valve %s with weight %d\n", valve.Name, weight)
    valve.Weight = weight
    for n := 0; n < len(valve.Connections); n++ {
        next_valve := valve.Connections[n]
        if next_valve.Weight != -1 && next_valve.Weight <= weight {
            continue
        }
        if valve.FlowRate == 0 {
            // +1 because no Flow Rate so no need to open
            calculateWeights(next_valve, weight+1)
        } else {
            // +2 because opening this valve takes 1 minute
            calculateWeights(next_valve, weight+2)
        }
    } 
}

func generateGraph(valves []Valve, filename string) {
    fmt.Println("Generating puml graph ...")
    puml_start := "@startuml\nhide empty description\n"
    puml_end := "@enduml\n"
    puml_text := puml_start
    for _, valve := range(valves) {
        puml_text += "\n" + valve.Name + " : " + "Flow Rate = " + strconv.Itoa(valve.FlowRate)
        puml_text += "\n" + valve.Name + " : " + "Path Cost: 1"
    }
    puml_text += "\n[*] --> AA\n"
    puml_text += "\n"
    for _, valve := range(valves) {
        for _, connection := range(valve.Connections) {
            puml_text += "\n" + valve.Name + " --> " + connection.Name
            puml_text += "\nnote on link\n  Path Cost: 1\nend note"
        }
    }
    puml_text += "\n" + puml_end
    // write to file
    ioutil.WriteFile(filename, []byte(puml_text), 0644)

    //fmt.Println(puml_text)
}

func calculateReleasedPressure(valves []Valve, time int) (int) {
    releasedPressure := 0
    for _, valve := range(valves) {
        releasedPressure += valve.GetReleasedPressure(time)
    }
    return releasedPressure
}

func doStep(valves []Valve, current_valve *Valve, time int, time_limit int) (int) {
    return 0
    fmt.Printf("=== Minute %d ===\n", time)
    fmt.Printf("At valve %v\n", current_valve)
    // recursion bottom
    if time > time_limit {
        // return complete result
        fmt.Printf("Reached time limit %d\n", time_limit)
        releasedPressure := calculateReleasedPressure(valves, time_limit)
        // released pressure was calculated, we can now close the valve
        current_valve.Close()
        return releasedPressure
    }
    // open valve
    fmt.Println(current_valve.isOpen())
    if !current_valve.isOpen() && current_valve.FlowRate != 0 {
        fmt.Printf("Opening valve %s\n", current_valve.Name)
        time += 1
        current_valve.Open(time)
    }
    // move to next valve
    maxReleasedPressure := 0
    for n := 0; n < len(current_valve.Connections); n++ {
        next_valve := current_valve.Connections[n]
        if next_valve.isOpen() {
            // skip an open valve
            //continue
        }
        fmt.Printf("Moving to valve %s\n", next_valve.Name)
        releasedPressure := doStep(valves, next_valve, time+1, time_limit)
        if releasedPressure > maxReleasedPressure {
            maxReleasedPressure = releasedPressure
        }
    }
    fmt.Println(maxReleasedPressure)
    // close current valve again before "leaving"
    current_valve.Close()
    fmt.Println("Reached the end")
    return maxReleasedPressure
}

func iterate(current_valve **Valve, current_time int) {
    if (*current_valve).isOpen() || (*current_valve).FlowRate == 0 {
        fmt.Println("Checking connections")
        for v, targetValve := range((*current_valve).Connections) {
            if !targetValve.isOpen() {
                fmt.Printf("Moving to valve %v\n", targetValve)
                (*current_valve) = (*current_valve).Connections[v]
            }
        }
    } else {
        fmt.Println("Opening Valve")
        (*current_valve).Open(current_time)
    }
}

func CalculateViaSteps(valves []Valve, time_to_spend int) {
    current_valve := &valves[0]
    maxReleasedPressure := doStep(valves, current_valve, 1, 30)
    fmt.Printf("Released Pressure: %d\n", maxReleasedPressure)
}

func SmartCalculation(valves []Valve, time_to_spend int) (int) {
    maxPossibleValues := make([][]int, len(valves))
    for v := 0; v < len(valves); v++ {
        valve := &valves[v]
        maxPossibleValues[v] = []int{(time_to_spend - valve.Weight) * valve.FlowRate, v}
    }
    sort.SliceStable(maxPossibleValues, func(i, j int) bool {
        return maxPossibleValues[i][0] > maxPossibleValues[j][0]
    })
    fmt.Println(maxPossibleValues)
    // do some magic
    start_valve := &valves[maxPossibleValues[0][1]]
    start_valve.Open(start_valve.Weight+1)
    for p := 1; p < len(maxPossibleValues); p++ {
        valve := &valves[maxPossibleValues[p][1]]
        fmt.Println("distance ... ")
        valve.Open(valve.Weight + start_valve.GetDistanceTo(valve) + 1)
        start_valve = valve
    }
    return calculateReleasedPressure(valves, time_to_spend)
}

func solve1(valves []Valve, time_to_spend int) {
    fmt.Println("Starting solving of part 1 ...")
    //CalculateViaSteps(valves, time_to_spend)
    releasedPressure := SmartCalculation(valves, time_to_spend)
    fmt.Println(releasedPressure)
    result := 0
    fmt.Printf("[Solve1] Result: %d\n", result)
}

func main() {
    input_lines := readFile(os.Stdin)
    valves := parseInput(input_lines)
    findConnections(valves)
    calculateWeights(&valves[0], 0)
    fmt.Println(valves)
    generateGraph(valves, "graph.puml")

    solve1(valves, 30)
}
