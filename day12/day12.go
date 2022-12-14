package main

import (
    "bufio"
    "fmt"
    "math"
    "os"

    "github.com/beefsack/go-astar"
)

var _ astar.Pather = Tile{}

type Grid struct {
    tiles       [][]*Tile
    start       [2]int
    end         [2]int
}

func (g Grid) TileAt(x int, y int) (*Tile) { return g.tiles[y][x] }
func (g Grid) ElevationAt(x, y int) int    { return g.TileAt(x, y).elevation }
func (g Grid) MaxX() int                   { return len(g.tiles[0]) - 1 }
func (g Grid) MaxY() int                   { return len(g.tiles) - 1 }

func (g Grid) print() {
    for _, tile_line := range(g.tiles) {
        for _, tile := range(tile_line) {
            fmt.Print(tile.elevationS)
        }
        fmt.Println()
    }
}

func (g Grid) ShortestPath(from *Tile) ([]*Tile, float64, bool) {
    path, dist, found := astar.Path(from, g.TileAt(g.end[0], g.end[1]))

    var tilePath []*Tile
    for _, step := range path {
        tilePath = append(tilePath, step.(*Tile))
    }
    return tilePath, dist, found
}

func (g Grid) ShortestPathFromElevation(startElevation int) (*Tile, int) {
    var shortest = math.MaxInt
    var startPos *Tile

    for x := 0; x < g.MaxX(); x++ {
        for y := 0; y < g.MaxY(); y++ {
            current_startPos := g.TileAt(x, y)
            if current_startPos.elevation != startElevation {
                continue
            }
            if path_len, wasFound := g.ShortestPathLen(current_startPos); wasFound && path_len < shortest {
                shortest = path_len
                startPos = current_startPos
            }
        }
    }

    return startPos, shortest
}

func (g Grid) ShortestPathLen(from *Tile) (int, bool) {
    path, _, found := g.ShortestPath(from)
    return len(path) -1 ,found
}

type Tile struct {
    elevation   int
    elevationS  string
    grid        *Grid
    x, y        int
}

func (t Tile) PathNeighbors() (out []astar.Pather) {
    if t.x > 0 && t.grid.TileAt(t.x-1, t.y).elevation <= t.elevation+1 {
		out = append(out, t.grid.TileAt(t.x-1, t.y))
	}
	if t.x < t.grid.MaxX() && t.grid.TileAt(t.x+1, t.y).elevation <= t.elevation+1 {
		out = append(out, t.grid.TileAt(t.x+1, t.y))
	}
	if t.y > 0 && t.grid.TileAt(t.x, t.y-1).elevation <= t.elevation+1 {
		out = append(out, t.grid.TileAt(t.x, t.y-1))
	}
	if t.y < t.grid.MaxY() && t.grid.TileAt(t.x, t.y+1).elevation <= t.elevation+1 {
		out = append(out, t.grid.TileAt(t.x, t.y+1))
	}

	return out
}

func (t Tile) PathNeighborCost(to astar.Pather) float64 {
    target := to.(*Tile)
    return float64(target.elevation) - float64(t.elevation) + 1
}

func (t Tile) PathEstimatedCost(to astar.Pather) float64 {
    target := to.(*Tile)
	return float64(t.manhattanDistance(target))
}

func (t Tile) manhattanDistance(target *Tile) int {
	return int(math.Abs(float64(target.x)-float64(t.x)) + math.Abs(float64(target.y)-float64(t.y)))
}

func solve1(grid *Grid) {
    fmt.Println("Staring solving for Part 1")
    grid.print()
    length, _ := grid.ShortestPathLen(grid.TileAt(grid.start[0], grid.start[1]))
    fmt.Printf("[Solve1] Found path of length %d\n", length)
}

func solve2(grid *Grid, level rune) {
    fmt.Println("Staring solving for Part 2")
    startPos, path_len := grid.ShortestPathFromElevation(int(level))
    fmt.Printf("[Solve2] Found shortest path of len %d from tile [%d|%d]\n", path_len, startPos.x, startPos.y)
}

func parseGrid(filename string) (*Grid) {
    grid := &Grid{}
    fmt.Printf("Reading file %s\n", filename)

    input_file, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    defer input_file.Close()

    fileScanner := bufio.NewScanner(input_file)
    fileScanner.Split(bufio.ScanLines)

    var y int = 0
    for fileScanner.Scan() {
        grid.tiles = append(grid.tiles, nil)
        line_text := fileScanner.Text()
        for x, c := range line_text {
            var val int
            switch (c) {
            case 'E':
                grid.end = [2]int{x, y}
                val = int('z')
            case 'S':
                grid.start = [2]int{x, y}
                val = int('a')
            default:
                val = int(c)
            }
            grid.tiles[y] = append(grid.tiles[y], &Tile{elevation: val, elevationS: string(c), grid: grid, x: x, y: y})
        }
        y++
    }
    return grid
}

func main() {
    grid := parseGrid("input.txt")
    solve1(grid)
    solve2(grid, 'a')
}
