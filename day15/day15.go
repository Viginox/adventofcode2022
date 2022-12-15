package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Position struct {
	X int
	Y int
}

type Sensor struct {
	Position      Position
	ClosestBeacon Position
	Distance      int
}

func (s *Sensor) manhattanDistance() (distance int) {
	return int(math.Abs(float64(s.ClosestBeacon.X)-float64(s.Position.X)) + math.Abs(float64(s.ClosestBeacon.Y)-float64(s.Position.Y)))
}

func (s *Sensor) Init(sensorPos Position, beaconPos Position) {
	s.Position = sensorPos
	s.ClosestBeacon = beaconPos
	s.Distance = s.manhattanDistance()

}

type Tile struct {
	value rune
}

type Grid struct {
	Tiles      [][]Tile
	XOffset    int
	YOffset    int
	XReduction int
	YReduction int
}

func (g *Grid) Init(width int, height int) {
	g.XReduction = math.MaxInt
	g.YReduction = math.MaxInt
	g.Tiles = make([][]Tile, height)
	for i := 0; i < height; i++ {
		g.Tiles[i] = make([]Tile, width)
	}
}

func (g *Grid) getTileAt(x int, y int) *Tile {
	return &g.Tiles[y+g.YOffset-g.YReduction][x+g.XOffset-g.XReduction]
}

func (g *Grid) GetWidth() int {
	if g.Tiles != nil {
		return len(g.Tiles[0])
	} else {
		return 0
	}
}

func (g *Grid) GetHeight() int {
	return len(g.Tiles)
}

func (g *Grid) Print() {
	for _, line := range g.Tiles {
		for _, tile := range line {
			if tile.value == 0 {
				fmt.Printf(".")
			} else {
				fmt.Printf("%c", tile.value)
			}
		}
		fmt.Printf("\n")
	}
}

func (g *Grid) ExtendAtEnd(xAxis int, yAxis int) {
	fmt.Printf("Extend at end x:%d, y:%d\n", xAxis, yAxis)
	if xAxis > 0 {
		for y, line := range g.Tiles {
			g.Tiles[y] = append(line, make([]Tile, xAxis)...)
		}
	}
	if yAxis > 0 {
		for i := 0; i < yAxis; i++ {
			g.Tiles = append(g.Tiles, make([]Tile, len(g.Tiles[0])))
		}
	}
}

func (g *Grid) InsertAtFront(xAxis int, yAxis int) {
	fmt.Printf("Insert at front for x:%d y:%d\n", xAxis, yAxis)
	if xAxis > 0 {
		for y, line := range g.Tiles {
			g.Tiles[y] = append(line, make([]Tile, xAxis)...)
			for i := len(g.Tiles[y]) - 1; i >= xAxis; i-- {
				g.Tiles[y][i] = g.Tiles[y][i-xAxis]
			}
			// delete first tiles
			for i := 0; i < xAxis; i++ {
				g.Tiles[y][i] = Tile{}
			}
		}
		if g.XOffset+xAxis > g.XReduction {
			g.XOffset += -(g.XReduction - xAxis)
            g.XReduction = 0
		} else if xAxis <= g.XReduction {
			g.XReduction -= xAxis
		}
	}
	if yAxis > 0 {
		// create new elements
		newTiles := make([][]Tile, yAxis)
		for i := 0; i < yAxis; i++ {
			newTiles[i] = make([]Tile, g.GetWidth())
		}
		g.Tiles = append(g.Tiles, newTiles...)
		// copy values
		for i := len(g.Tiles) - 1; i >= yAxis; i-- {
			g.Tiles[i] = g.Tiles[i-yAxis]
		}
		// delete first rows
		for i := 0; i < yAxis; i++ {
			g.Tiles[i] = make([]Tile, g.GetWidth())
		}
		// adapt offset
		if g.YOffset+yAxis > g.YReduction {
			g.YOffset += -(g.YReduction - yAxis)
            g.YReduction = 0
		} else if yAxis <= g.YReduction {
			g.YReduction -= yAxis
		}
	}
}

func (g *Grid) GrowToFit(sensor Sensor) {
	// let the first sensor be at x=0 and y=0
	// SensorPosition = RawPosition + Offset - Reduction
	// Offset to compensate negative positions
	// Reduction to remove too high values
	//
	// Check the X position
    /*
	if g.XReduction == math.MaxInt {
		g.XReduction = sensor.Position.X
	}
	if g.YReduction == math.MaxInt {
		g.YReduction = sensor.Position.Y
	}
    */
	sensor_x := sensor.Position.X + g.XOffset - g.XReduction
	max_x := sensor_x + sensor.Distance
	if max_x >= g.GetWidth() {
		g.ExtendAtEnd(max_x-g.GetWidth()+1, 0)
	}
	min_x := sensor_x - sensor.Distance
	if min_x < 0 {
		g.InsertAtFront(-(min_x), 0)
	}
    // ignore
    return
	sensor_y := sensor.Position.Y + g.YOffset - g.YReduction
	if sensor_y+sensor.Distance >= g.GetHeight() {
		g.ExtendAtEnd(0, (sensor_y+sensor.Distance)-g.GetHeight()+1)
	}
	if sensor_y-sensor.Distance < 0 {
		g.InsertAtFront(0, -(sensor_y - sensor.Distance))
	}
}

func (g *Grid) DrawSensor(sensor Sensor, search_row int) {
    //fmt.Printf("Drawing sensor %v\n", sensor)
	//g.GrowToFit(sensor)
	// mark sensor position
    if sensor.Position.Y == search_row {
        g.getTileAt(sensor.Position.X, 0).value = 'S'
    }
	// mark beacon position
    if sensor.ClosestBeacon.Y == search_row {
	    g.getTileAt(sensor.ClosestBeacon.X, 0).value = 'B'
    }
	// draw detection range
	//for y := sensor.Position.Y - sensor.Distance; y <= sensor.Position.Y+sensor.Distance; y++ {
    y := search_row
		dist := sensor.Position.Y - y
		if dist < 0 {
			dist *= -1
		}
		for x := sensor.Position.X - (sensor.Distance - dist); x <= sensor.Position.X+(sensor.Distance-dist); x++ {
            if y == search_row {
			    tile := g.getTileAt(x, 0)
			    if tile.value == 0 {
			    	tile.value = '#'
			    }
            }   
		}
	//}
}

func (g *Grid) CountBeaconNotPresentForRow(row int) (count int) {
	//line := g.Tiles[row+g.YOffset-g.YReduction]
    line := g.Tiles[0]
	count = 0
    for _, tile := range line {
		if tile.value == '#' {
			count++
		}
	}
	return count
}

func (g *Grid) GetPossibleBeaconPosition(search_start int, search_end int) (pos Position) {
    for y := 0; y < g.GetHeight(); y++ {
        for x := search_start; x < search_end; x++ {
            tile := g.Tiles[y][x]
            if tile.value == 0 {
                return Position{
                    X: x-g.XOffset,
                    Y: y,
                }
            }
        }
    }
    return Position{X: -1, Y: -1}
} 

func parseSensor(line string) (sensor Sensor) {
	parts := strings.Split(line, ":")
	if len(parts) >= 2 {
		sensor_string := parts[0]
		beacon_string := parts[1]
		pos_re := regexp.MustCompile(`^(?:.*?)x=(-?\d+), y=(-?\d+)`)
		sensor_match := pos_re.FindStringSubmatch(sensor_string)
		if sensor_match == nil {
			panic(errors.New("Sensor position could not be found in input string!"))
		}
		sensorX, _ := strconv.ParseInt(sensor_match[1], 10, 0)
		sensorY, _ := strconv.ParseInt(sensor_match[2], 10, 0)
		sensorPos := Position{
			X: int(sensorX),
			Y: int(sensorY),
		}
		beacon_match := pos_re.FindStringSubmatch(beacon_string)
		if beacon_match == nil {
			panic(errors.New("Beacon position could not be found in input string!"))
		}
		beaconX, _ := strconv.ParseInt(beacon_match[1], 10, 0)
		beaconY, _ := strconv.ParseInt(beacon_match[2], 10, 0)
		beaconPos := Position{
			X: int(beaconX),
			Y: int(beaconY),
		}
		sensor.Init(sensorPos, beaconPos)
	}
	return sensor
}

func getGridSize(sensors []Sensor) (oLimits []int) {
	checkPos := func(pos Position, dist int, limits []int) {
		if (pos.X-dist) < limits[0] {
			limits[0] = pos.X-dist
		} else if (pos.X+dist) > limits[1] {
			limits[1] = pos.X+dist
		}
		if (pos.Y-dist) < limits[2] {
			limits[2] = pos.Y-dist
		} else if (pos.Y+dist) > limits[3] {
			limits[3] = pos.Y+dist
		}
	}
	var limits = []int{0, 0, 0, 0}
	for _, sensor := range sensors {
		checkPos(sensor.Position, sensor.Distance, limits)
		checkPos(sensor.ClosestBeacon, 0, limits)
	}
	return limits
}

func parseInput(input_lines []string) (sensors []Sensor) {
	for _, line := range input_lines {
		sensors = append(sensors, parseSensor(line))
	}
	return sensors
}

func readFile(filename string) (lines []string) {
	fmt.Println("Reading file " + filename)

	input_file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer input_file.Close()

	fileScanner := bufio.NewScanner(input_file)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		lines = append(lines, strings.TrimSpace(fileScanner.Text()))
	}

	return lines
}

func solve1(grid Grid, limits []int, sensors []Sensor, search_row int) {
	// skip limits and create a default array. it will grow with demand
	//grid.Init(limits)
	grid.Init((-limits[0])+limits[1]+1, 1)
    grid.XOffset = -limits[0]
    grid.XReduction = 0
    grid.YReduction = 0
	for _, sensor := range sensors {
		if sensor.Position.Y-sensor.Distance <= search_row && search_row <= sensor.Position.Y+sensor.Distance {
			grid.DrawSensor(sensor, search_row)
			//break
		}
	}
	//grid.Print()
	row := search_row
	count := grid.CountBeaconNotPresentForRow(row)
	fmt.Printf("Counted %d positions in line %d where no beacon can be\n", count, row)
    fmt.Printf("[Solve1] Result: %d\n", count)
}

func checkBeacon(sensors []Sensor, search_start int, search_end int) Position {
    for search_row := search_start; search_row < search_end; search_row++ {
        //fmt.Printf("Searching in row %d\n", search_row)
        for search_point := search_start; search_point < search_end; search_point++ {
            //fmt.Printf("Search Point: %d\n", search_point)
            inSensorRange := false
            for _, sensor := range(sensors) {
                if !(sensor.Position.Y - sensor.Distance <= search_row && search_row <= sensor.Position.Y + sensor.Distance) {
                    continue
                }
                //fmt.Printf("Checking sensor %v\n", sensor)
                dist := sensor.Position.Y - search_row
                if dist < 0 {
                    dist *= -1
                }
                start_x := sensor.Position.X - (sensor.Distance - dist)
                end_x := sensor.Position.X + (sensor.Distance - dist)
                if start_x <= search_point && search_point < end_x {
                    search_point = end_x
                    inSensorRange = true
                    break
                }
            }
            if inSensorRange {
                continue
            }
            if search_point < search_end {
                return Position{X: search_point, Y: search_row}
            }
        }
    }
    return Position{X: -1, Y: -1}
}

func solve2(limits []int, sensors []Sensor, search_start int, search_end int) {
    fmt.Printf("Searching from %d to %d\n", search_start, search_end)
    var beaconPos Position
    beaconPos = checkBeacon(sensors, search_start, search_end) 
    fmt.Printf("Found possible beacon pos at %v\n", beaconPos)
    result := beaconPos.X * 4000000 + beaconPos.Y
    fmt.Printf("[Solve2] Result: %d\n", result)
}

func main() {
	input_lines := readFile("input.txt")
	sensors := parseInput(input_lines)
	limits := getGridSize(sensors)
	fmt.Printf("Limits are: %v\n", limits)
	grid := Grid{}

	solve1(grid, limits, sensors, 2000000)
    solve2(limits, sensors, 0, 4000000)
}
