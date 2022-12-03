#!/usr/bin/bash

file="day1.txt"

elfCalories=()

idx=0
elfCal=0

while read line; do
    if [[ -n $line ]]; then
        #echo -e "$line"
        elfCal=$((elfCal + line))
    else
        #echo -e "empty line"
        elfCalories[$idx]=$elfCal
        elfCal=0
        idx=$((idx + 1))
    fi
done < $file

elfCalories[$ids]=$elfCal

echo $elfCalories
