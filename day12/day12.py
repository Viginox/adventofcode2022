import sys
import copy

from pathfinding.core.diagonal_movement import DiagonalMovement
from pathfinding.core.grid import Grid
from pathfinding.finder.a_star import AStarFinder

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
                start_pos = [i, j]
            elif pos == 'E':
                end_pos = [i, j]
    return start_pos, end_pos


def checkTile(grid, currentTile, nextTile) -> bool:
    print(f"Checking traversal from tile {currentTile} to {nextTile}.")
    print(f"Values are: {grid[currentTile[0]][currentTile[1]]} and {grid[nextTile[0]][nextTile[1]]}")
    current_elevation = 0
    target_elevation = 0
    if grid[currentTile[0]][currentTile[1]] == 'S':
        current_elevation = ord("a")
    else:
        current_elevation = ord(grid[currentTile[0]][currentTile[1]])
    if grid[nextTile[0]][nextTile[1]] == 'E':
        target_elevation = ord("z")
    else:
        target_elevation = ord(grid[nextTile[0]][nextTile[1]])
    if (target_elevation - current_elevation) > 1:
        print("- not traversable")
        return False
    print("- traversable")
    return True

def getPathLength(path_grid):
    count = 0
    for line in path_grid:
        count += line.count('.')
    # number of possible steps in grid
    path_values = (len(path_grid) * len(path_grid[0]))
    # minus 1 to subtract the starting point
    path_values -= 1
    return path_values - count

def findPath(grid, path_grid, startPos, endPos):
    printGrid(path_grid)
    #print(f"Investigating Tile {startPos}")
    # Directions: left, up, right, down
    # recursion bottom
    if startPos == endPos:
        print(f"This is the top. Hooray!")
        path_grid[startPos[0]][startPos[1]] = 'E'
        printGrid(path_grid)
        path_length = getPathLength(path_grid)
        print(f"Path length is {path_length}")
        return path_length
    paths = [["?"] * len(grid) for i in range (len(grid[0]))] * 4
    paths_lengths = [0] * 4
    # go down
    if startPos[0] == len(grid)-1:
        # do nothing
        pass
        #print("downmost tile")
    else:
        path_d = copy.deepcopy(path_grid)
        nextPos = [startPos[0]+1, startPos[1]]
        if path_d[nextPos[0]][nextPos[1]] == '.' and checkTile(grid, startPos, nextPos):
            path_d[startPos[0]][startPos[1]] = 'v'
            paths_lengths[3] = findPath(grid, path_d, nextPos, endPos)
            paths[3] = path_d
    # go right
    if startPos[1] == len(grid[0])-1:
        # do nothing
        pass
        #print("rightmost tile")
    else:
        path_r = copy.deepcopy(path_grid)
        nextPos = [startPos[0], startPos[1]+1]
        if path_r[nextPos[0]][nextPos[1]] == '.' and checkTile(grid, startPos, nextPos):
            path_r[startPos[0]][startPos[1]] = '>'
            paths_lengths[2] = findPath(grid, path_r, nextPos, endPos)
            paths[2] = path_r
    # go up
    if startPos[0] == 0:
        # do nothing
        pass
        #print("upmost tile")
    else:
        path_u = copy.deepcopy(path_grid)
        nextPos = [startPos[0]-1, startPos[1]]
        if path_u[nextPos[0]][nextPos[1]] == '.' and checkTile(grid, startPos, nextPos):
            path_u[startPos[0]][startPos[1]] = '^'
            paths_lengths[1] = findPath(grid, path_u, nextPos, endPos)
            paths[1] = path_u
    # go left
    if startPos[1] == 0:
        # do nothing
        pass
        #print("leftmost tile")
    else:
        path_l = copy.deepcopy(path_grid)
        nextPos = [startPos[0], startPos[1]-1]
        if path_l[nextPos[0]][nextPos[1]] == '.' and checkTile(grid, startPos, nextPos):
            path_l[startPos[0]][startPos[1]] = '<'
            paths_lengths[0] = findPath(grid, path_l, nextPos, endPos)
            paths[0] = path_l
    #print(paths_lengths)
    smallest_idx = 0
    smallest_size = sys.maxsize
    for idx, size in enumerate(paths_lengths):
        #print(f"IDX {idx} has size {size}")
        if size != 0 and size < smallest_size:
            smallest_size = size
            smallest_idx = idx
    #print(f"Smallest path length is {smallest_size}")
    path_grid = paths[smallest_idx]
    #printGrid(path_grid)
    return smallest_size

def readFile(filename):
    puzzle_grid = []
    input_file = open(filename, "r")
    for line in input_file:
        puzzle_line = []
        for char in line.strip():
            puzzle_line.append(char)
        puzzle_grid.append(puzzle_line)
    return puzzle_grid

def solve1(puzzle_input):
    print("[Solve1] Finding the shortest path to the top")
    startPos, endPos = findLandmarks(puzzle_input)
    print(f"Starting Point at [{startPos[0]}|{startPos[1]}] and End Point at [{endPos[0]}|{endPos[1]}]")
    path_grid = [['.'] * len(puzzle_input[0]) for i in range(len(puzzle_input))]
    length = findPath(puzzle_input, path_grid, startPos, endPos)
    print(length)
    #path = astar(puzzle_input, (startPos[0], startPos[1]), (endPos[0], endPos[1]))
    #print(path)
    #print(f"Path has length {length}")
    #grid = Grid(matrix=puzzle_input)
    #start = grid.node(startPos[0], startPos[1])
    #end = grid.node(endPos[0], endPos[1])

    #finder = AStarFinder(diagonal_movement=DiagonalMovement.never)
    #path, runs = finder.find_path(start, end, grid)
    #print(path)
    print(f"[Solve1] Shortest path has length {length}")

def solve2():
    print(f"[Solve2] Nothing done yet")

class Node():
    """A node class for A* Pathfinding"""

    def __init__(self, parent=None, position=None):
        self.parent = parent
        self.position = position

        self.g = 0
        self.h = 0
        self.f = 0

    def __eq__(self, other):
        print(f"NODE COMP: {self.position} == {other.position} --> {self.position == other.position}")
        return self.position == other.position


def astar(maze, start, end):
    """Returns a list of tuples as a path from the given start to the given end in the given maze"""

    # Create start and end node
    start_node = Node(None, start)
    start_node.g = start_node.h = start_node.f = 0
    end_node = Node(None, end)
    end_node.g = end_node.h = end_node.f = 0

    # Initialize both open and closed list
    open_list = []
    closed_list = []

    # Add the start node
    open_list.append(start_node)

    # Loop until you find the end
    while len(open_list) > 0:

        # Get the current node
        current_node = open_list[0]
        current_index = 0
        for index, item in enumerate(open_list):
            if item.f < current_node.f:
                current_node = item
                current_index = index

        # Pop current off open list, add to closed list
        open_list.pop(current_index)
        closed_list.append(current_node)

        # Found the goal
        if current_node == end_node:
            print(f"Reached the goal!")
            path = []
            current = current_node
            while current is not None:
                path.append(current.position)
                current = current.parent
            return path[::-1] # Return reversed path

        # Generate children
        children = []
        for new_position in [(0, -1), (0, 1), (-1, 0), (1, 0)]: # Adjacent squares

            # Get node position
            node_position = (current_node.position[0] + new_position[0], current_node.position[1] + new_position[1])

            # Make sure within range
            if node_position[0] >= (len(maze)) or node_position[0] < 0 or node_position[1] >= (len(maze[0])) or node_position[1] < 0:
                continue

            # Make sure walkable terrain
            if not checkTile(maze, [current_node.position[0], current_node.position[1]], [node_position[0], node_position[1]]):
                print("Terrain not walkable")
                continue
            #if maze[node_position[0]][node_position[1]] != 0:
            #    continue

            # Create new node
            new_node = Node(current_node, node_position)

            # Append
            children.append(new_node)

        # Loop through children
        for child in children:

            # Child is on the closed list
            for closed_child in closed_list:
                if child == closed_child:
                    continue

            # Create the f, g, and h values
            child.g = current_node.g + 1
            child.h = ((child.position[0] - end_node.position[0]) ** 2) + ((child.position[1] - end_node.position[1]) ** 2)
            child.f = child.g + child.h

            # Child is already in the open list
            for open_node in open_list:
                if child == open_node and child.g >= open_node.g:
                    continue

            # Add the child to the open list
            open_list.append(child)
            print(open_list)
            print(closed_list)


if __name__ == "__main__":
    puzzle_input = readFile("input.txt")
    puzzle_len = len(puzzle_input)
    puzzle_width = len(puzzle_input[0])
    print(f"Read puzzle of size {puzzle_len}x{puzzle_width}")
    printGrid(puzzle_input)

    solve1(puzzle_input)
    solve2()
