use std::fs::File;
use std::io::{BufRead, BufReader};

fn solve1<R: BufRead>(reader: &mut R) -> u32 {
    let mut total_rucksack_type_value: u32 = 0;
    for (index, line) in reader.lines().enumerate() {
        let line = line.unwrap();
        // show the line and its numbers.
        //println!("{}. {}", index + 1, line)
        let (l1, l2) = line.split_at(line.len()/2);
        let mut common_item_type: char = ' ';
        let mut common_item_type_value: u32 = 0;
        for c in l1.chars() {
            if l2.contains(c) {
                common_item_type = c;
                common_item_type_value = get_type_value(common_item_type);
                //println!("Rucksack {}: Value of {} is {}", index+1, c, common_item_type_value);
                total_rucksack_type_value += common_item_type_value;
                break;
            }
        }
    }
    return total_rucksack_type_value;
}

fn get_type_value(type_char: char) -> u32 {
    let mut type_char_value: u32 = type_char as u32;
    if type_char_value >= 'a' as u32 && type_char_value <= 'z' as u32{
        type_char_value -= 'a' as u32 -1;
    } else if type_char_value >= 'A' as u32 && type_char_value <= 'Z' as u32 {
        type_char_value -= 'A' as u32 -1;
        type_char_value += 26;
    }
    return type_char_value;
}

fn solve2<R: BufRead>(reader: &mut R) -> u32 {
    let group_size = 3;
    let mut rucksack_group: [String; 2] = Default::default();
    let mut total_group_rucksack_value: u32 = 0;
    for (index, line) in reader.lines().enumerate() {
        let line = line.unwrap();
        if (index+1) % group_size == 0 {
            // stored a group of 3 rucksacks
            let line2 = &rucksack_group[0];
            let line3 = &rucksack_group[1];
            let mut common_item_type: char = ' ';
            let mut common_item_type_value: u32 = 0;
            for c in line.chars() {
                if line2.contains(c) && line3.contains(c) {
                    common_item_type = c;
                    common_item_type_value = get_type_value(common_item_type);
                    total_group_rucksack_value += common_item_type_value;
                    break;
                }
            }
            rucksack_group = Default::default();
        } else {
            rucksack_group[index % group_size] = line;
        }
        
    }
    return total_group_rucksack_value;
}

fn main() {
    let filename = "./input.txt";
    // open the file in read-only-mode
    let file = File::open(filename).unwrap();
    let mut reader = BufReader::new(file);

    let rucksack_total_priority_value = solve1(&mut reader);
    println!("Total rucksack priority: {}", rucksack_total_priority_value);

    let file2 = File::open(filename).unwrap();
    let mut reader2 = BufReader::new(file2);
    let group_rucksack_priority_value = solve2(&mut reader2);
    println!("Total group rucksack priority: {}", group_rucksack_priority_value);

}
