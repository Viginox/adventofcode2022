#!/usr/bin/bash

declare -a register_values=(1)
current_register_value=1

function executeCommand () {
    local command=$1
    #echo "Executing command '$command'"
    IFS=' ' read -ra command_values <<< "$command"
    local op=${command_values[0]}
    #echo "Found operation '$op'"
    if [[ $op == "noop" ]]; then
        #echo "Found noop"
        register_values+=($current_register_value)
    elif [[ $op == "addx" ]]; then
        local val=${command_values[1]}
        #echo "Found addx with value $val"
        register_values+=($current_register_value)
        #echo "Add value $val to register with value $current_register_value"
        current_register_value=$(( $current_register_value + $val ))
        #echo "Now has $current_register_value"
        register_values+=($current_register_value)
    else
        echo "PANIC: Unsupported operation"
    fi
}

function solve1 () {
    local input_filename=$1
    local start_cycle=$2
    local cycle_spacing=$3
    local end_cycle=$4
    echo "Reading file $input_filename and processing commands"
    while IFS= read -r line;
    do
        executeCommand "$line"
    done < $input_filename
    echo "Command execution complete"
    echo "Computing signal strength for selected cycles"
    local cycles=`seq $start_cycle $cycle_spacing $end_cycle`
    local total_signal_strength=0
    for cycle in ${cycles[@]}; do
        echo "Value for cycle $cycle"
        local r_val=${register_values[$(( $cycle -1 ))]}
        local signal_strength=$(( $r_val * $cycle ))
        echo "$cycle * $r_val = $signal_strength"
        total_signal_strength=$(( $total_signal_strength + $signal_strength ))
    done
    echo "[Solve1] Total signal strength is $total_signal_strength"
}

function drawCRT () {
    local current_crt_pos=0
    local crt_width=$1
    echo "Drawing values on a CRT with width $crt_width"
    for value in ${register_values[@]}; do
        if [[ $current_crt_pos == $crt_width ]]; then
            current_crt_pos=0
            echo ""
        fi
        if [[ $current_crt_pos == $value ]] || [[ $current_crt_pos == $(( $value + 1 )) ]] || [[ $current_crt_pos == $(( $value - 1 )) ]]; then
            echo -n "#"
        else
            echo -n "."
        fi
        current_crt_pos=$(( $current_crt_pos + 1 ))
    done
    echo ""
}

function solve2 () {
    local crt_width=$1
    drawCRT $crt_width
    echo "[Solve2] Please read solution from CRT screen above"
}

input_filename="input.txt"

solve1 $input_filename 20 40 220
solve2 40
