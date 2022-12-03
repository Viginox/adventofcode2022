
input_filename = "./input.txt"

rock    = 1  # A X
paper   = 2  # B Y
scissor = 3  # C Z

shape_scores = {
    "A": 1,
    "B": 2,
    "C": 3,
}

relations = {
    "X": "A",
    "Y": "B",
    "Z": "C"
}

result_relations = {
    "X": "lose",
    "Y": "draw",
    "Z": "win"
}

results = {
    "A": {"A": "draw", "B": "lose", "C": "win"},
    "B": {"A": "win", "B": "draw", "C": "lose"},
    "C": {"A": "lose", "B": "win", "C": "draw"},
}

shape_for_results = {
    "A": {"draw": "A", "win": "B", "lose": "C"},
    "B": {"lose": "A", "draw": "B", "win": "C"},
    "C": {"win": "A", "lose": "B", "draw": "C"}
}

result_scores = {
    "win": 6,
    "lose": 0,
    "draw": 3
}


def readInput(filename):
    file = open(filename, 'r')
    lines = file.readlines()
    print(f"Found {len(lines)} lines in file {filename}")
    return lines

def solve1(file_lines):
    total_score = 0
    for line in file_lines:
        parts = line.rstrip("\n").split(" ")
        O = parts[0]
        P = relations[parts[1]]
        shape_score = shape_scores[P]
        result = results[P][O]
        result_score = result_scores[result]
        game_score = shape_score + result_score
        total_score += game_score
    return total_score


def solve2(file_lines):
    total_score = 0
    for line in file_lines:
        parts = line.rstrip("\n").split(" ")
        # shape the opponent will choose
        O = parts[0]
        # result to end the round with
        R = parts[1]
        result = result_relations[R]
        # calculate the shape to choose for the dedicated result
        P = shape_for_results[O][result]
        # get the score for the chosen shape
        shape_score = shape_scores[P]
        # get the score for the set result
        result_score = result_scores[result]
        game_score = shape_score + result_score
        total_score += game_score
    return total_score


if __name__ == "__main__":
    file_lines = readInput(input_filename)
    res1 = solve1(file_lines)
    print(f"Result of Part1: {res1}")
    res2 = solve2(file_lines)
    print(f"Result of Part2: {res2}")
