import re
import copy
import sys

crate_stacks = []
steps = []

def read_input(filename):
    stack_width = 3
    stack_space = 1
    input_file = open(filename, 'r')
    lines = input_file.readlines()

    section = 0
    for line in lines:
        if line.rstrip() == "":
            section = 1
            continue
        if section == 0:
            crate_stacks.append(line.rstrip("\n"))
        else:
            steps.append(line.strip())

    noStacks = round(len(crate_stacks[0])/(stack_width+stack_space))
    #print(f"Number of Stacks: {noStacks}")
    stacks = []
    for i in range (noStacks):
        stacks.append([])
    # split lines into sections of crates
    for stack in crate_stacks[:-1]:
        for c, i in enumerate(range (0, len(stack), stack_width+stack_space)):
            current_crate = stack[i:i+stack_width].strip()
            if current_crate != "":
                stacks[c].insert(0, current_crate)
    #print(f"Crate Stack:\n{stacks}")
    #print(f"\nSteps:\n{steps}")
    return stacks, steps

def parseStep(step):
    step_re = re.compile("^move\s(?P<amount>\d+)\sfrom\s(?P<source>\d+)\sto\s(?P<dest>\d+)$")
    step_data = step_re.match(step)
    #print(step_data)
    # use -1 to correct for starting an array with 0 instead of 1
    return int(step_data[1]), int(step_data[2])-1, int(step_data[3])-1

def moveCrate(stacks, source, dest):
    crate = stacks[source].pop()
    stacks[dest].append(crate)

def performStep9000(stacks, step):
    # print(f"Do {step}")
    #print(f"On Stacks: {stacks}")
    amount, source, dest = parseStep(step)
    #print(f"Amount: {amount}\nSource: {source}\nDestination: {dest}")
    for i in range (amount):
        moveCrate(stacks, source, dest)

def moveCrates(stacks, source, amount, dest):
    crates = []
    for i in range(amount):
        crates.append(stacks[source].pop())
    for i in range(amount):
        stacks[dest].append(crates.pop())

def performStep9001(stacks, step):
    #print(f"Do {step}")
    #print(f"On Stacks: {stacks}")
    amount, source, dest = parseStep(step)
    moveCrates(stacks, source, amount, dest)

def displayStack(stacks):
    for stack in stacks:
        for crate in stack:
            print(crate, end=" ")
        print()

def printSolution(stacks):
    print("Solution:", end=" ")
    for stack in stacks:
        print(stack[-1].lstrip("[").rstrip("]"), end="")
    print()

def solve1(stacks, steps):
    for s, step in enumerate(steps):
        sys.stdout.write(f"\r{s} of {len(steps)-1}")
        performStep9000(stacks, step)
    sys.stdout.flush()
    print("\nDone!")
    displayStack(stacks)
    printSolution(stacks)
    print("Solved 1!")

def solve2(stacks, steps):
    print("Nothing done yet")
    for s, step in enumerate(steps):
        sys.stdout.write(f"\r{s} of {len(steps)-1}")
        performStep9001(stacks, step)
    sys.stdout.flush()
    print("\nDone!")
    displayStack(stacks)
    printSolution(stacks)
    print("Solved 2!")


if __name__ == "__main__":
    stacks, steps = read_input("./input.txt")
    stacks1 = copy.deepcopy(stacks)
    stacks2 = copy.deepcopy(stacks)

    displayStack(stacks1)
    print("Input Stacks:")
    solve1(stacks1, steps)
    print("Input Stacks:")
    displayStack(stacks2)
    solve2(stacks2, steps)
