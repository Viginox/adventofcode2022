package main

import (
    "fmt"
    "os"
    "bufio"
    "strconv"
)

func read_input(input_filename string) ([][]int) {
    arr := [][]int{}
    fmt.Println("Loading file " + input_filename)
    readFile, err := os.Open(input_filename)

    if err != nil {
        panic(err)
    }

    fileScanner := bufio.NewScanner(readFile)
    fileScanner.Split(bufio.ScanLines)

    var row []int
    for fileScanner.Scan() {
        row = make([]int, 0)
        line_text := fileScanner.Text()
        for _, c := range line_text {
            height, _ := strconv.ParseInt(string(c), 10, 64)
            row = append(row, int(height))
        }
        arr = append(arr, row)
    }
    fmt.Println("Read " + strconv.Itoa(len(arr)) + " lines with length " + strconv.Itoa(len(arr[0])))
    return arr
}

func checkVisibility(input_array [][]int) ([][]bool) {
    visibility := make([][]bool, 0)
    var v_row []bool
    for i, row := range input_array {
        v_row = make([]bool, 0)
        for j, pos := range row {
            v_left := true
            v_right := true
            v_top := true
            v_bottom := true
            //fmt.Println("Checking pos [" + strconv.Itoa(i) + "|" + strconv.Itoa(j) + "]")
            // it's a border
            if (i == 0 || j == 0 || i == len(input_array)-1 || j == len(row)-1 ) {
                // the border is always visible
                v_row = append(v_row, true)
                //fmt.Println("Found a border!")
                continue
            }
            // check top
            for p_i := 0; p_i < i; p_i++ {
                if input_array[p_i][j] >= pos {
                    v_top = false
                    break
                }
            }
            // check bottom
            for p_i := i+1; p_i < len(input_array); p_i++ {
                if input_array[p_i][j] >= pos {
                    v_bottom = false
                    break
                }
            }
            // check left
            for p_j := 0; p_j < j; p_j++ {
                if input_array[i][p_j] >= pos {
                    v_left = false
                    break
                }
            }
            // check right
            for p_j := j+1; p_j < len(row); p_j++ {
                if input_array[i][p_j] >= pos {
                    v_right = false
                    break
                }
            }
            v_row = append(v_row, v_top || v_bottom || v_left || v_right)
        }
        visibility = append(visibility, v_row)
    }
    return visibility
}

func countValues(input_array [][]bool, value bool) (count int) {
    count = 0
    fmt.Println("Checking values for array of size " + strconv.Itoa(len(input_array)) + "x" + strconv.Itoa(len(input_array[0])))
    for _, row := range input_array {
        // fmt.Println("Line " + strconv.Itoa(i) + " has length " + strconv.Itoa(len(row)))a
        for _, pos := range row {
            if (pos == value) {
                // fmt.Println("Found value at pos [" + strconv.Itoa(i) + "|" + strconv.Itoa(j) + "]")
                count += 1
            }
        }
    }
    return count
}

func scenicScore(input_array [][]int) (score_array [][]int) {
    score_array = make([][]int, 0)
    var score_row []int
    for i, row := range input_array {
        score_row = make([]int, 0)
        for j, pos := range row {
            //fmt.Println("Checking pos [" + strconv.Itoa(i) + "|" + strconv.Itoa(j) + "]")
            // check top
            score_top := i
            for p_i := i-1; p_i >= 0; p_i-- {
                if input_array[p_i][j] >= pos {
                    score_top = i - p_i
                    break
                }
            }
            score_bottom := (len(input_array)-1) - i
            // check bottom
            for p_i := i+1; p_i < len(input_array); p_i++ {
                if input_array[p_i][j] >= pos {
                    score_bottom = p_i - i
                    break
                }
            }
            score_left := j
            // check left
            for p_j := j-1; p_j >= 0; p_j-- {
                if input_array[i][p_j] >= pos {
                    score_left = j - p_j
                    break
                }
            }
            score_right := (len(row)-1) -j
            // check right
            for p_j := j+1; p_j < len(row); p_j++ {
                if input_array[i][p_j] >= pos {
                    score_right = p_j - j
                    break
                }
            }
            score_row = append(score_row, score_top * score_bottom * score_left * score_right)
        }
        score_array = append(score_array, score_row)
    }
    return score_array
}

func findHighestValue(input_array[][]int) (highest int, x int, y int) {
    highest = 0
    x = 0
    y = 0
    for i, row := range input_array {
        for j, pos := range row {
            if (pos > highest) {
                highest = pos
                x = i
                y = j
            }
        }
    }
    return highest, x, y
}

func solve1(input_array [][]int) {
    visibility_array := checkVisibility(input_array)
    /*
    for _, row := range visibility_array {
        fmt.Println(row)
    }
    */
    trees_visible := countValues(visibility_array, true)
    fmt.Println("[Solve1] Counted " + strconv.Itoa(trees_visible) + " visible trees in grid!")
}

func solve2(input_array [][]int) {
    score_array := scenicScore(input_array)
    /* 
    for _, row := range score_array {
        fmt.Println(row)
    }
    */
    highes_scenic_score, x, y := findHighestValue(score_array)
    fmt.Println("[Solve2] Highes scenic score is " + strconv.Itoa(highes_scenic_score) + " at position [" + strconv.Itoa(x) + "|" + strconv.Itoa(y) + "]")
    fmt.Println("[Solve2] Nothing done yet")
}

func main() {
    input_array := read_input("input.txt")
    /*
    for _, row := range input_array {
        fmt.Println(row)
    }
    */
    solve1(input_array)
    solve2(input_array)
    
}
