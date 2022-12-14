from astar import AStar

import copy

def printGrid(grid):
    for line in grid:
        for pos in line:
            print(pos, end="")
        print()

def findLandmarks(grid):
    start_pos = (0, 0)
    end_pos = (0, 0)
    for i, line in enumerate(grid):
        for j, pos in enumerate(line):
            if pos == 'S':
                start_pos = (i, j)
            elif pos == 'E':
                end_pos = (i, j)
    return start_pos, end_pos

def checkTile(grid, currentTile, nextTile) -> bool:
    print(f"Checking traversal from tile {currentTile} to {nextTile}.")
    print(f"Values are: {grid[currentTile[0]][currentTile[1]]} and {grid[nextTile[0]][nextTile[1]]}")
    if grid[currentTile[0]][currentTile[1]] == 'S': # and grid[nextTile[0]][nextTile[1]] == 'a':
        print("- traversable (because from starting position)")
        return True
    if grid[currentTile[0]][currentTile[1]] == 'z' and grid[nextTile[0]][nextTile[1]] == 'E':
        print("- traversable (because to end position)")
        return True
    elif grid[nextTile[0]][nextTile[1]] == 'E':
        print("- not traversable (way to end only allowed from z)")
        return False
    # mark ground more than 1 higher than the current tile as not traversable
    if ord(grid[currentTile[0]][currentTile[1]]) - ord(grid[nextTile[0]][nextTile[1]]) < -1:
        print("- not traversable")
        return False
    print("- traversable")
    return True

def readFile(filename):
    puzzle_grid = []
    input_file = open(filename, "r")
    for line in input_file:
        puzzle_line = []
        for char in line.strip():
            puzzle_line.append(char)
        puzzle_grid.append(puzzle_line)
    return puzzle_grid

def printSolution(grid, path):
    localGrid = copy.deepcopy(grid)
    for step in path:
        localGrid[step[0]][step[1]] = localGrid[step[0]][step[1]].upper()
    for line in localGrid:
        for elem in line:
            print(elem, end="")
        print()

def solve1(puzzle_grid, startPos, endPos):
    res = MazeSolver(puzzle_grid).astar(startPos, endPos)
    print(res)
    printSolution(puzzle_grid, res)
    print(len(res) -1)

def solve2(puzzle_grid):
    pass

class MazeSolver(AStar):

    def __init__(self, maze):
        self.lines = maze
        self.width = len(self.lines[0])
        self.height = len(self.lines)

    def heuristic_cost_estimate(self, n1, n2):
        (x1, y1) = n1
        (x2, y2) = n2
        return ((x1 - x2) ** 2) + ((y1 - y2) ** 2)

    def distance_between(self, n1, n2):
        n1_elevation = ord(self.lines[n1[1]][n1[0]])
        n2_elevation = ord(self.lines[n2[1]][n2[0]])
        return n2_elevation - n1_elevation + 1

    def neighbors(self, node):
        x, y = node
        return[(nx, ny) for nx, ny in[(x, y - 1), (x, y + 1), (x - 1, y), (x + 1, y)]if 0 <= nx < self.width and 0 <= ny < self.height]

if __name__ == "__main__":
    puzzle_grid = readFile("test_input.txt")
    print(f"Puzzle has size {len(puzzle_grid[0])}x{len(puzzle_grid)}")
    startPos, endPos = findLandmarks(puzzle_grid)
    solve1(puzzle_grid, startPos, endPos)
    solve2(puzzle_grid)
