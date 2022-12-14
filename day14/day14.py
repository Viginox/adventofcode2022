import sys
import copy

def printGrid(grid):
    for line in grid:
        for x, pos in enumerate(line):
            if x > 0:
                print(pos, end="")
        print()

def readFile(filename):
    puzzle_input = []
    with open(filename, "r") as input_file:
        for line in input_file:
            puzzle_input.append(line.strip())
    return puzzle_input

def adaptGrid(grid, width, height, infinite=False):
    for y, line in enumerate(grid):
        fill_char = '.'
        if infinite and y == len(grid)-1:
            fill_char = '#'
        while len(line) < width:
            line.append(fill_char)
    while len(grid) < height:
        grid.append(['.'] * width)

def drawRockLine(grid, start, end):
    #print(f"Drawing rock from {start} to {end}")
    adaptGrid(grid, start[0]+1, start[1]+1)
    adaptGrid(grid, end[0]+1  , end[1]+1  )
    if start[0] == end[0]:
        if start[1] < end[1]:
            y_start = start[1]
            y_end = end[1]
        else:
            y_start = end[1]
            y_end = start[1]
        for y in range(y_start, y_end+1):
            grid[y][start[0]] = "#"
    else:
        if start[0] < end[0]:
            x_start = start[0]
            x_end = end[0]
        else:
            x_start = end[0]
            x_end = start[0]
        for x in range(x_start, x_end+1):
            grid[start[1]][x] = "#"

def parsePuzzleLine(puzzle_grid, puzzle_line):
    #print(f"Parsing puzzle line {puzzle_line}")
    # draw rock for all selected places
    line_sections = puzzle_line.split(" -> ")
    section_coords = []
    for section in line_sections:
        section_text = section.split(",")
        section_coord = [int(section_text[0]), int(section_text[1])]
        section_coords.append(section_coord)
    for i, coord in enumerate(section_coords):
        if i == len(section_coords)-1:
            break
        drawRockLine(puzzle_grid, coord, section_coords[i+1])

def isSandInGrid(grid, sand):
    if sand[0] >= len(grid[0]) or sand[1] >= len(grid)-1 or sand[0] < 0:
        return False
    return True

def simulateSandUnit(grid, sand_unit, infinite=False):
    #print(f"Sand at {sand_unit}")
    if infinite:
        adaptGrid(grid, sand_unit[0]+2, len(grid), infinite)
        #printGrid(grid)
    if not isSandInGrid(grid, sand_unit):
        return False, sand_unit
    # move down
    if grid[sand_unit[1]+1][sand_unit[0]] == ".":
        return simulateSandUnit(grid, [sand_unit[0], sand_unit[1]+1], infinite)
    # down left
    elif grid[sand_unit[1]+1][sand_unit[0]-1] == ".":
        return simulateSandUnit(grid, [sand_unit[0]-1, sand_unit[1]+1], infinite)
    # down right
    elif grid[sand_unit[1]+1][sand_unit[0]+1] == ".":
        return simulateSandUnit(grid, [sand_unit[0]+1, sand_unit[1]+1], infinite)
    # done
    else:
        grid[sand_unit[1]][sand_unit[0]] = "o"
        return True, sand_unit

def simulateSandDrop(grid, sand_souce, infinite=False):
    count = 0
    # spawn sand unit
    while True:
        sand_unit = sand_source
        success, rest_pos = simulateSandUnit(grid, sand_unit, infinite)
        #printGrid(grid)
        if not infinite and not success:
            break
        count += 1
        print(f"Current count: {count}", end="\r")
        if infinite and rest_pos == sand_source:
            break
    return count

def parsePuzzleInput(puzzle_input):
    puzzle_grid = [['.'] * 10 for i in range(10)]
    for line in puzzle_input:
        parsePuzzleLine(puzzle_grid, line)
    return puzzle_grid

def solve1(puzzle_grid, sand_source):
    print("Start solving part 1")
    sand_count = simulateSandDrop(puzzle_grid, sand_source)
    #printGrid(puzzle_grid)
    print(f"[Solve1] Result: {sand_count} sand units were counted")

def solve2(puzzle_grid, sand_source):
    print("Start solving part 2")
    # adapt grid
    drawRockLine(puzzle_grid, [0, len(puzzle_grid)+1], [len(puzzle_grid[0]), len(puzzle_grid)+1])
    #printGrid(puzzle_grid)
    # simulate sand
    sand_count = simulateSandDrop(puzzle_grid, sand_source, True)
    #printGrid(puzzle_grid)
    print(f"[Solve2] Result: {sand_count} sand units were counted")

if __name__ == "__main__":
    puzzle_input = readFile("input.txt")
    sand_source = [500, 0]

    puzzle_grid = parsePuzzleInput(puzzle_input)
    puzzle_grid[sand_source[1]][sand_source[0]] = "+"
    puzzle_grid1 = copy.deepcopy(puzzle_grid)
    puzzle_grid2 = copy.deepcopy(puzzle_grid)
    #printGrid(puzzle_grid)

    solve1(puzzle_grid1, sand_source)
    solve2(puzzle_grid2, sand_source)

