$input_file = Get-Content "./input.txt"

function CheckBuffer ($Buffer) {
	$b1 = 0
	foreach ($c in $Buffer) {
		$b2 = 0
		foreach ($d in $Buffer) {
			if ($b1 -eq $b2) {
				# skip because the same index
			} else {
				if ($c -eq $d) {
					return $false
				}
			}
			$b2 += 1
		}
		$b1 += 1
	}
	return $true
}

function FindMarker($input_buffer, $marker_size) {
	Write-Host "Looking for Marker of size " $marker_size " ..."
	$marker_buffer = [System.Collections.ArrayList]::new()
	$marker_pos = 1
	foreach ($c in $input_buffer) {
		[void]$marker_buffer.Add($c)
		if ($marker_buffer.Count -eq $marker_size) {
			$is_marker = CheckBuffer $marker_buffer
			if ($is_marker) {
				Write-Host "Found Marker at position " $marker_pos
				break
			} else {
				Write-Host "No marker found in buffer " $marker_buffer
			}
			$marker_buffer.RemoveAt(0)
		}
		$marker_pos += 1	
	}
}

function Solve1() {
	Write-Host "Solving Part 1 ..."
	FindMarker $input_file.ToCharArray() 4	
}

function Solve2() {
	Write-Host "Solving Part 2 ..."
	FindMarker $input_file.ToCharArray() 14
}

Solve1
Solve2
