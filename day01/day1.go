package main

import (
    "fmt"
    "bufio"
    "os"
    "strconv"
)

func main() {
    readFile, err := os.Open("day1.txt")

    if err != nil {
        panic(err)
    }

    fileScanner := bufio.NewScanner(readFile)

    fileScanner.Split(bufio.ScanLines)

    var elves_data []int64

    var current_elve int64 = 0
    for fileScanner.Scan() {
        line_text := fileScanner.Text()
        if line_text == "" {
            fmt.Println("empty Line")
            elves_data = append(elves_data, current_elve)
            current_elve = 0
        } else {
            line_number, err := strconv.ParseInt(line_text, 10, 64)
            if err != nil {
                panic(err)
            }
            current_elve += line_number
            fmt.Println(line_number)
        }
    }
    // find biggest value
    var biggest_elve_nr int = 0
    var firstbiggest_elve_cal int64 = 0
    var secbiggest_elve_cal int64 = 0
    var thirdbiggest_elve_cal int64 = 0
    for i := 0; i < len(elves_data); i++ {
        if elves_data[i] > firstbiggest_elve_cal {
            biggest_elve_nr = i+1
            thirdbiggest_elve_cal = secbiggest_elve_cal
            secbiggest_elve_cal = firstbiggest_elve_cal
            firstbiggest_elve_cal = elves_data[i]
        } else if elves_data[i] > secbiggest_elve_cal {
            thirdbiggest_elve_cal = secbiggest_elve_cal
            secbiggest_elve_cal = elves_data[i]
        } else if elves_data[i] > thirdbiggest_elve_cal {
            thirdbiggest_elve_cal = elves_data[i]
        }
    }

    fmt.Println(elves_data)

    fmt.Println("Elve Nr " + strconv.Itoa(biggest_elve_nr) + " carries " + strconv.Itoa(int(firstbiggest_elve_cal)) + " calories")

    fmt.Println("Total cal of top3 elves: " + strconv.Itoa(int(firstbiggest_elve_cal + secbiggest_elve_cal + thirdbiggest_elve_cal)))
}
