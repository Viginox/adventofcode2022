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
                grid[i][j] = 'a'
                start_pos = [i, j]
            elif pos == 'E':
                grid[i][j] = 'z'
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

def readFile(filename):
    puzzle_grid = []
    input_file = open(filename, "r")
    for line in input_file:
        puzzle_line = []
        for char in line.strip():
            puzzle_line.append(char)
        puzzle_grid.append(puzzle_line)
    return puzzle_grid

class GridNode:
    def __init__(self, parent=None, position=None):
        self.parent = parent
        self.position = position

        self.f = 0
        self.g = 0
        self.h = 0

    def __eq__(self, other):
        return self.position == other.position

def AStar(startPos, endPos, grid, heuristic=None):
    startNode = GridNode(None, startPos)

    endNode = GridNode(None, endPos)

    openList = []
    closedList = []

    openList.append(startNode)

    while openList:
        print("-- Next Item")

        # get node with lowest f value
        current_node = openList[0]
        current_index = 0
        for index, item in enumerate(openList):
            if item.f < current_node.f:
                current_node = item
                current_index = index

        # Remove from open list and add to closed list
        openList.pop(current_index)
        closedList.append(current_node)

        print(f"Current position: {current_node.position}")
        print(f"End position: {endNode.position}")
        # condition for goal reached
        if current_node == endNode:
            print(f"Reached the goal!")
            path = []
            current = current_node
            while current is not None:
                path.append(current.position)
                current = current.parent
            return path[::-1] # Return reversed path

        children = []
        for nextPosition in [(1, 0), (0, 1), (-1, 0),(0, -1)]:

            nodePosition = (current_node.position[0] + nextPosition[0], current_node.position[1] + nextPosition[1])
            print(f"Target node is {nodePosition}")

            # check range
            if nodePosition[0] >= len(grid) or nodePosition[0] < 0 \
                    or nodePosition[1] >= len(grid[0]) or nodePosition[1] < 0:
                print("--> out of range")
                continue

            # check condition
            if not checkTile(grid, [current_node.position[0], current_node.position[1]], [nodePosition[0], nodePosition[1]]):
                print("--> not walkable")
                continue

            # create new node from child
            newNode = GridNode(current_node, nodePosition)

            # Append to child list
            children.append(newNode)

        # compute children
        for child in children:
            print(f"--- Next child {child.position}")

            # check if already closed
            if child in closedList:
                print("Already computed")
                continue

            # calculate the f, g and h values
            child.g = current_node.g + 1
            child.h = abs((child.position[0] - endNode.position[0]) ** 2) + abs((child.position[1] - endNode.position[1]) ** 2)
            child.f = child.g + child.h
            print(f"--- has g value {child.g}")

            # check if already listed to be computed
            if child in openList:
                open_node = 0
                for o_node in openList:
                    if o_node == child:
                        open_node = o_node
                # hope that it was found - but should have been
                if child.g > open_node.g:
                    print("Alread in list - skipping")
                    continue

            # add to open list
            print(f"--- - adding to list")
            openList.append(child)

def printSolution(grid, path):
    localGrid = copy.deepcopy(grid)
    for step in path:
        localGrid[step[0]][step[1]] = '.' #localGrid[step[0]][step[1]].upper()
    for line in localGrid:
        for elem in line:
            print(elem, end="")
        print()

def solve1(puzzle_grid, startPos, endPos):
    res = AStar((startPos[0], startPos[1]), (endPos[0], endPos[1]), puzzle_grid)
    print(res)
    printSolution(puzzle_grid, res)
    print(len(res) -1)

def solve2(puzzle_grid):
    pass

if __name__ == "__main__":
    puzzle_grid = readFile("input.txt")
    startPos, endPos = findLandmarks(puzzle_grid)
    solve1(puzzle_grid, startPos, endPos)
    solve2(puzzle_grid)
