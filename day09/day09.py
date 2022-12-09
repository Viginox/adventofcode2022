

def readInput(filename):
    input_file = open(filename,"r")
    input_lines = input_file.readlines()
    return input_lines

def moveHead(head_pos, direction):
    if (direction == "U"):
        # up equals x--
        #print("Move up")
        head_pos[0] -= 1
    elif (direction == "D"):
        # down equals x++
        #print("Move down")
        head_pos[0] += 1
    elif (direction == "L"):
        # left equals y--
        #print("Move left")
        head_pos[1] -= 1
    elif (direction == "R"):
        # right equals y++
        #print("Move right")
        head_pos[1] += 1
    else:
        print(f"ERROR: Unknown direction: {direction}")

def moveKnot(head_pos, tail_pos):
    # move tail to new position
    # check head pos in relation to tail pos
    x_rel = tail_pos[0] - head_pos[0]
    y_rel = tail_pos[1] - head_pos[1]
    #print(f"Relative Knot position: [{int(x_rel)}|{int(y_rel)}]")
    # realtive position is only in 1 axis
    if (abs(x_rel) > 1 and y_rel == 0):
        tail_pos[0] -= int(x_rel / abs(x_rel))
    elif (abs(y_rel) > 1 and x_rel == 0):
        tail_pos[1] -= int(y_rel / abs(y_rel))
    elif (abs(x_rel) > 1 or abs(y_rel) > 1):
        tail_pos[0] -= int(x_rel / abs(x_rel))
        tail_pos[1] -= int(y_rel / abs(y_rel))
    #print(f"Move Knot to [{tail_pos[0]}|{tail_pos[1]}]")

def doCommand(head_pos, tail_pos, visit_array, command):
    direction, step_size = command.strip().split(" ")
    for i in range(int(step_size)):
        moveHead(head_pos, direction)
        #print(f"Head moved to [{head_pos[0]}|{head_pos[1]}]")
        moveKnot(head_pos, tail_pos)
        visit_array[tail_pos[0]][tail_pos[1]] = '#'

def doLongerCommand(knots_pos, visit_array, command):
    direction, step_size = command.strip().split(" ")
    for i in range(int(step_size)):
        # move head
        moveHead(knots_pos[0], direction)
        for i in range(1, len(knots_pos)-1):
            #print(f"Knot {i}")
            moveKnot(knots_pos[i-1], knots_pos[i])
        # move tail
        moveKnot(knots_pos[-2], knots_pos[-1])
        visit_array[knots_pos[-1][0]][knots_pos[-1][1]] = '#'
        #printKnots(knots_pos, len(visit_array))

def printKnots(knots_pos, array_size):
    visual_array = [['.'] * array_size for i in range(array_size)]
    for i, knot in enumerate(knots_pos):
        if i == 0:
            visual_array[knot[0]][knot[1]] = 'H'
        elif i == len(knots_pos)-1:
            visual_array[knot[0]][knot[1]] = 'T'
        else:
            visual_array[knot[0]][knot[1]] = i
    visual_array[knots_pos[0][0]][knots_pos[0][1]] = 'H'
    print('-'*array_size)
    for line in visual_array:
        for dot in line:
            print(dot, end="")
        print


def countVisitedLocations(visit_array):
    locs_visited = 0
    for line in visit_array:
        locs_visited += line.count('#')
    return locs_visited

def solve1(command_lines, array_size):
    head_pos = [int(array_size/2), int(array_size/2)]
    tail_pos = [int(array_size/2), int(array_size/2)]
    visit_array = [['.']*array_size for i in range(array_size)]
    for command in command_lines:
        doCommand(head_pos, tail_pos, visit_array, command)
    visit_count = countVisitedLocations(visit_array)
    print(f"Counted {visit_count} visited locations")
    print(f"[Solve1] Result: {visit_count}")

def solve2(command_lines, knots, array_size):
    knots_pos = [[int(array_size/2), int(array_size/2)] for i in range(knots)]
    visit_array = [['.'] * array_size for i in range(array_size)]
    for command in command_lines:
        doLongerCommand(knots_pos, visit_array, command)
    visit_count = countVisitedLocations(visit_array)
    print(f"Counted {visit_count} visited locations")
    print(f"[Solve2] Result: {visit_count}")

if __name__ == "__main__":
    command_lines = readInput("input.txt")

    solve1(command_lines, 500)
    solve2(command_lines, 10, 500)
