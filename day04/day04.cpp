#include <iostream>
#include <fstream>
#include <string>

std::fstream readFile(std::string filename) {
    std::fstream input_file;
    input_file.open(filename, std::ios::in);

    if (input_file.is_open()) {
    }
    
    return input_file;
}

int* parseLine(std::string line) {
    std::string elve_delimiter = ",";
    std::string range_delimiter = "-";
    std::string token;

    std::string elve1 = line.substr(0, line.find(elve_delimiter));
    std::string elve1_start_s = elve1.substr(0, elve1.find(range_delimiter));
    int elve1_start = std::stoi(elve1_start_s);
    std::string elve1_end_s = elve1.substr(elve1_start_s.length()+range_delimiter.length());
    int elve1_end = std::stoi(elve1_end_s);
    std::string elve2 = line.substr(elve1.length()+elve_delimiter.length());
    std::string elve2_start_s = elve2.substr(0, elve2.find(range_delimiter));
    int elve2_start = std::stoi(elve2_start_s);
    std::string elve2_end_s = elve2.substr(elve2_start_s.length()+range_delimiter.length());
    int elve2_end = std::stoi(elve2_end_s);

    std::cout << "Elve 1: " << elve1_start << " until " << elve1_end << " | Elve 2: " << elve2_start << " until " << elve2_end << std::endl;
    int* ret = new int(4);
    ret[0] = elve1_start;
    ret[1] = elve1_end;
    ret[2] = elve2_start;
    ret[3] = elve2_end;
    return ret;
}

size_t solve1(int input_arr[]) {
    if (input_arr[0] >= input_arr[2] && input_arr[1] <= input_arr[3]) {
        std::cout << ">> Elve 1 is contained in Elve 2!" << std::endl;
        return 1;
    } else if (input_arr[2] >= input_arr[0] && input_arr[3] <= input_arr[1]) {
        std::cout << ">> Elve 2 is contained in Elve 1" << std::endl;
        return 1;
    }
    return 0;
}

size_t solve2(int input_arr[]) {
    int elve1_start = input_arr[0];
    int elve1_end   = input_arr[1];
    int elve2_start = input_arr[2];
    int elve2_end   = input_arr[3];

    if ((elve1_start >= elve2_start && elve1_start <= elve2_end) || (elve1_end >= elve2_start && elve1_end <= elve2_end) || (elve1_start <= elve2_start && elve1_end >= elve2_end)) {
        std::cout << ">> >> Elve 1 overlaps Elve 2" << std::endl;
        return 1;
    } else if ((elve2_start >= elve1_start && elve2_start <= elve1_end) || (elve2_end >= elve1_start && elve2_end <= elve1_end) || (elve2_start <= elve1_start && elve2_end >= elve1_end)) {
        std::cout << ">> >> Elve2 overlaps Elve 1" << std::endl;
        return 1;
    }

    return 0;
}

int main(void) {
    std::cout << "Let's go!" << std::endl;

    std::fstream input_file = readFile("input.txt");
    std::string line;

    size_t noFullyContianedSectors = 0;
    size_t noOverlapingSecotrs = 0;
    while (getline(input_file, line)) {
        int* line_result = parseLine(line);
        noFullyContianedSectors += solve1(line_result);
        noOverlapingSecotrs += solve2(line_result);
        delete line_result;
    }
    std::cout << "There are " << noFullyContianedSectors << " fully contained sectors!" << std::endl;
    std::cout << "There are " << noOverlapingSecotrs << " overlapping sectors!" << std::endl;


    return 0;
}
