package main

import (
    "bufio"
    "errors"
    "fmt"
    "io"
    "os"
    "strings"

    "github.com/cheggaaa/pb"
)

var RockShapes [][][]rune = [][][]rune{
    {{'#', '#', '#', '#'}},
    {{ 0 , '#',  0 }, {'#', '#', '#'}, { 0 , '#',  0 }},
    {{'#', '#', '#'}, { 0 ,  0 , '#'}, { 0 ,  0 , '#'}},
    {{'#'}, {'#'}, {'#'}, {'#'}},
    {{'#', '#'}, {'#', '#'}},
}

type Position struct {
    X int
    Y int
}

type Rock struct {
    Position Position
    Shape [][]rune
    ShapeIndex int
    Width int
    Cave *Cave
    isMoving bool
}

func (r *Rock) Push(direction rune) {
    //fmt.Printf("Direction: %c\n", direction)
    switch (direction) {
    case '<':
        //fmt.Println("Move left")
        if r.Position.X > 0 && !r.Collides(-1) {
            r.Position.X -= 1
        } else {
            //fmt.Println("Nothing happened")
        }
    case '>':
        //fmt.Println("Pushed right")
        if r.Position.X + r.Width < r.Cave.Width && !r.Collides(1){
            r.Position.X += 1
        } else {
            //fmt.Println("Nothing happened")
        }
    default:
        panic(errors.New("Unknown direction detected! -> " + string(direction)))
    }
}

func (r *Rock) Collides(x_offset int) (collision bool) {
    for row := 0; row < len(r.Shape); row++ {
        for pos := 0; pos < len(r.Shape[row]); pos++ {
            if r.Shape[row][pos] != '#' {
                continue
            }
            if r.Cave.Grid[r.Position.Y+row][r.Position.X+pos+x_offset] == '#' {
                return true
            }
        }
    }
    return false
}

func (r *Rock) hasHitGround() (hitGround bool) {
    for row := 0; row < len(r.Shape); row++ {
        for pos := 0; pos < len(r.Shape[row]); pos++ {
            if r.Shape[row][pos] == '#' && (r.Position.Y+row == 0 || r.Cave.Grid[r.Position.Y+row-1][r.Position.X+pos] == '#') {
                return true
            }
        }
    }
    return false
}

func (r *Rock) Drop() (hitGround bool) {
    //fmt.Println("Dropping ...")
    // push the rock
    r.Push(r.Cave.GetJetDirection())
    // check, if we hit the ground
    onGround := r.hasHitGround()
    if onGround {
        newTopRow := r.Position.Y + len(r.Shape)
        if newTopRow > r.Cave.CurrentTopRow {
            r.Cave.CurrentTopRow = r.Position.Y + len(r.Shape)
        }
        return onGround
    }
    // drop rock by 1
    r.Position.Y -= 1
    return false
}

func (r *Rock) ShapeEquals(shape [][]rune) (equals bool) {
    if len(r.Shape) != len(shape) || len(r.Shape[0]) != len(shape[0]) {
        return false
    }
    for row := 0; row < len(r.Shape); row++ {
        for pos := 0; pos < len(r.Shape[row]); pos++ {
            if r.Shape[row][pos] != shape[row][pos] {
                return false
            }
        }
    }
    return true
}

type Cave struct {
    Rocks []Rock
    Width int
    Grid [][]rune
    JetDirections []rune
    JetIndex int
    CurrentTopRow int
    ReductionCount int
}

func (c *Cave) Init(width int) {
    c.Width = width
    c.Grid = append(c.Grid, make([]rune, width))
}

func (c *Cave) GetWidth() (int) {
    return c.Width
}

func (c *Cave) GetHeight() (int) {
    return len(c.Grid)
}

func (c *Cave) PurgeOld() {
    if len(c.Grid) > 1000000 {
        c.Grid = c.Grid[500000:]
        c.ReductionCount += 1
        c.CurrentTopRow -= 500000
    }
}

func (c *Cave) GetJetDirection() (rune) {
    jet_direction := c.JetDirections[c.JetIndex]
    c.JetIndex = (c.JetIndex+1) % len(c.JetDirections)
    return jet_direction
}

func (c *Cave) ExtendTo(height int) {
    //fmt.Printf("Extending to height %d\n", height)
    for {
        if len(c.Grid) >= height {
            break
        }
        c.Grid = append(c.Grid, make([]rune, c.Width))
    }
}

func (c *Cave) GetCaveRowAt(height int) (row []rune) {
    return c.Grid[height]
}

func (c *Cave) Draw(rock *Rock, symbol ...rune) {
    rock_symbol := '#'
    if len(symbol) > 0 {
        rock_symbol = symbol[0]
    }
    for row := 0; row < len(rock.Shape); row++ {
        for pos := 0; pos < len(rock.Shape[row]); pos++ {
            if rock.Shape[row][pos] == '#' {
                c.Grid[rock.Position.Y+row][rock.Position.X+pos] = rock_symbol
            }   
        }
    }
}

func (c *Cave) Print() {
    for r := 0; r < len(c.Rocks); r++ {
        rock := &c.Rocks[r]
        if rock.isMoving {
            //c.Draw(rock, '@')
        }
    }
    for row := len(c.Grid)-1; row >= 0; row-- {
        fmt.Print("|")
        for pos := 0; pos < len(c.Grid[row]); pos++ {
            if c.Grid[row][pos] == 0 {
                fmt.Print(".")
            } else {
                fmt.Printf("%c", c.Grid[row][pos])
            }
        }
        fmt.Println("|")
    }
    fmt.Print("+")
    fmt.Print(strings.Repeat("-", len(c.Grid[0])))
    fmt.Println("+")
}

func (c *Cave) SpawnNewRock() (*Rock) {
    spawn_x := 2
    spawn_y := c.CurrentTopRow + 3
    rockShape := RockShapes[len(c.Rocks) % len(RockShapes)]
    c.ExtendTo(spawn_y + len(rockShape))
    newRock := Rock{
        Shape: rockShape,
        ShapeIndex: len(c.Rocks) % len(RockShapes),
        Width: len(rockShape[0]),
        Position: Position{X: spawn_x, Y: spawn_y},
        Cave: c,
        isMoving: true,
    }
    //fmt.Printf("Spawned new rock at %d | %d\n", spawn_x, spawn_y)
    c.Rocks = append(c.Rocks, newRock)
    return &c.Rocks[len(c.Rocks)-1]
}

func readFile(r io.Reader) (input_values []rune) {
    scanner := bufio.NewScanner(r)
    scanner.Split(bufio.ScanRunes)

    for scanner.Scan() {
        text_runes := []rune(scanner.Text())
        if text_runes[0] == 10 {
            break
        }
        input_values = append(input_values, text_runes...)
    }

    return input_values
}

func solve1(jet_directions []rune, rocks_to_drop int) (int) {
    cave := Cave{}
    cave.Init(7)
    cave.JetDirections = jet_directions
    fmt.Println("Starting solving of part 1")
    fmt.Printf("Dropping rocks...\n")
    bar := pb.StartNew(rocks_to_drop)
    for i := 0; i < rocks_to_drop; i++ {
        //fmt.Printf("Rock %d of %d\r", i, rocks_to_drop)
        bar.Increment()
        rock := cave.SpawnNewRock()
        for {
            hitGround := rock.Drop()
            if hitGround {
                rock.isMoving = false
                cave.Draw(rock)
                break
            }
        }
    }
    bar.Finish()
    return cave.CurrentTopRow
}

func solve2(jet_directions []rune, rocks_to_drop int) (int) {
    cave := Cave{}
    cave.Init(7)
    cave.JetDirections = jet_directions
    fmt.Println("Starting solving of part 2")
    count := 0
    var tracked [][]int 
    for {
        rock := cave.SpawnNewRock()

        for _, tracked_rock := range(tracked) {
            prev_shape_index := tracked_rock[0]
            prev_jet_index := tracked_rock[1]
            prev_count := tracked_rock[2]
            prev_top_row := tracked_rock[3]
            if prev_shape_index != rock.ShapeIndex || prev_jet_index != rock.Cave.JetIndex {
                continue
            } else {
                period := count - prev_count
                if count % period == rocks_to_drop % period {
                    fmt.Println("Found a cycle. Wohoo!")

                    cycle_height := rock.Cave.CurrentTopRow - prev_top_row
                    rocks_remaining := rocks_to_drop - count
                    cycles_remaining := (rocks_remaining / period) + 1
                    return prev_top_row + (cycle_height * cycles_remaining)
                }
            }
        }
        tracked = append(tracked, []int{rock.ShapeIndex, rock.Cave.JetIndex, count, rock.Cave.CurrentTopRow})

        for {
            hitGround := rock.Drop()
            if hitGround {
                rock.isMoving = false
                cave.Draw(rock)
                break
            }
        }
        fmt.Printf("Dropped rocks: %d\r", count)
        if cave.JetIndex == 0 {
            fmt.Println(rock.Shape)
        }
        if rock.ShapeEquals(RockShapes[0]) && cave.JetIndex == 0 && count != 0 {
            break
        }
        count++
    }
    fmt.Println()
    return count
}

func main() {
    jet_directions := readFile(os.Stdin)
    fmt.Printf("Found %d jet directions\n", len(jet_directions))
    height := solve1(jet_directions, 2022)
    fmt.Printf("Stones stack up to a height of %d\n", height)

    rocks_to_simulate := 1000000000000
    // assuming a repetetive pattern
    height = solve2(jet_directions, rocks_to_simulate)
    fmt.Printf("Stones stack up to aheight of %d\n", height)
}
