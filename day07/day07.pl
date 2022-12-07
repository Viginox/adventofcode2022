#!/usr/bin/perl
use warnings;
use strict;

my $input_filename = "./input.txt";
my @input_lines;

my %parsed_structure = ( type => "dir", path => "/", name => "/", size => -1, children => [], parent => 0 );

sub read_input_file {

	my $i_filename = $_[0];
	open(my $ifh, '<', $i_filename) or die $!;

	while (<$ifh>) {
		chomp;
		push @input_lines, $_;
	}
	# close the file again
	close($ifh);
}

## Node Structure:
# Type (dir or file)
# Name
# Size
# Children
# Parent

sub parseLine {
	my $line = $_[0];
	my ($command, $command_parameter, $folder_name, $size, $file_name) = ($line =~ m/^(?:\$\s(?P<command>.*?)(?: (?P<command_parameter>.+?))?)$|^(?:dir (?P<folder_name>.+?))$|^(?:(?P<size>\d+) (?P<file_name>.*?))$/);
	#print("Command: $command\n");
	#print("Command Parameter: $command_parameter\n");
	#print("Folder Name: $folder_name\n");
	#print("Size: $size\n");
	#print("File Name: $file_name\n");
	return $command, $command_parameter, $folder_name, $size, $file_name;
}

sub computeStructure {
	my $current_path = "/";
	my $current_struct = \%parsed_structure;
	foreach my $line (@input_lines) {
		my ($command, $command_parameter, $folder_name, $size, $file_name) = parseLine($line);
		if (defined $command && $command eq "cd") {
			if ($command_parameter eq "..") {
				# changing one dir up
				my @fields = split /\//, $current_path;
				$current_path = join('/', @fields[0..$#fields-1]);
				$current_struct = %$current_struct{"parent"};
			} elsif ($command_parameter eq "/") {
				# changed to root path
				$current_path = "/";
				$current_struct = \%parsed_structure;
			} else {
				my $separator = "/";
				if ($current_path eq "/") {
					$separator = "";
				}
				$current_path .= $separator . $command_parameter;
				# change parent to correct structure
				foreach my $child (@{$current_struct->{"children"}}) {
					#print("Child: ", $child->{"name"}, "\n");
					if ($child->{"name"} eq $command_parameter) {
						#print("Updated the current_struct to child $child->{'name'}\n");
						$current_struct = $child;
						last;
					}
				}
			}
			#print("Changed path to $current_path\n");
		} elsif (defined $command && $command eq "ls") {
			# now follows a list of the current directory
			#print("listing files ...\n");
		}
		if (defined $size ) {
			# file
			my $separator = "/";
			if ($current_path eq "/") {
				$separator = "";
			}
			my $file_path = $current_path . $separator . $file_name;
			my %file_struct = ( type => "file", path => $file_path, name => $file_name, size => $size, children => [], parent => $current_struct );
			push(@{$current_struct->{"children"}}, \%file_struct);
			#print("Found file   $file_path with size $file_struct{'size'}\n");
		} elsif (defined $folder_name) {
			# folder
			my $separator = "/";
			if ($current_path eq "/") {
				$separator = "";
			}
			my $folder_path = $current_path . $separator . $folder_name;
			my %folder_struct = ( type => "dir", path => $folder_path, name => $folder_name, size => -1, children => [], parent => $current_struct );
			push(@{$current_struct->{"children"}}, \%folder_struct);
			#print("Found folder $folder_path\n");
		}
	}
	#print "$_ $parsed_structure{$_}\n" for (keys %parsed_structure);
	print("# Completed Structure Computation\n")
}

sub calculateTotalSize {
	my $node = $_[0];
	my $node_size = 0;
	#print("Calculating ", $node->{"path"}, "\n");
	# check if the current node is a file
	#print("Type: ", $node->{"type"}, "\n");
	if ($node->{"type"} eq "file") {
		# then return the file size and exit
		#print("Found file $node->{'name'} and returning size ", $node->{"size"}, "\n");
		return $node->{"size"};
	}
	# check if the dir has children
	if (defined $node->{"children"}) {
		for my $child (@{$node->{"children"}}) {
			$node_size += calculateTotalSize($child);
			#print("Processing Child ", $child->{"path"}, ". Size: $node_size\n");
		}
	}
	$node->{"size"} = $node_size;
	return $node_size;
}

sub calculateTotalDirBelow {
	my $node  = $_[0];
	my $below = $_[1];
	my $total_size = 0;

	if ($node->{"type"} eq "file") {
		# ignore files for this process
		return $total_size;
	}
	# only dir will end here
	if (defined $node->{"children"}) {
		for my $child (@{$node->{"children"}}) {
			$total_size += calculateTotalDirBelow($child, $below);
		}
	}
	if ($node->{"size"} <= $below) {
		#print("Found a dir (size: $node->{'size'}), smaller than $below\n");
		$total_size += $node->{"size"};
	}
	#print("Current total size: $total_size\n");
	return $total_size;
}

sub getDirSmallerThan {
	my $node = $_[0];
	my $required_disk_space = $_[1];
	my @possible_dirs = ();

	# ignore the root dir in this case ... that's easier
	if (defined $node->{"children"}) {
		for my $child (@{$node->{"children"}}) {
			if ($child->{"size"} >= $required_disk_space && $child->{"type"} eq "dir") {
				push(@possible_dirs, $child);
				push(@possible_dirs, getDirSmallerThan($child, $required_disk_space));
			}
		}
	}
	return @possible_dirs;
}

sub findSmallestDirToDelete {
	my $required_disk_space = $_[0];

	my @possible_dirs = getDirSmallerThan(\%parsed_structure, $required_disk_space);

	my $smallest_dir = @possible_dirs[0];
	foreach my $dir (@possible_dirs) {
		#print("$dir->{'name'}: $dir->{'size'}\n");
		if ($dir->{"size"} < $smallest_dir->{"size"}) {
			$smallest_dir = $dir;
		}
	}
	return $smallest_dir->{"size"};
}

sub solve1 {
	my $limit = 100000;
	my $total_size = calculateTotalSize(\%parsed_structure);
	print("Root Node size:       $parsed_structure{size}\n");
	print("Total Size:           $total_size\n");
	my $result = calculateTotalDirBelow(\%parsed_structure, $limit);
	print("==========================================\n");
	print("[Solve1] Result:      $result\n");
}

sub solve2 {
	my $total_disk_space = 70000000;
	my $required_unused_space = 30000000;
	my $total_size = $parsed_structure{"size"};
	my $available_disk_space = $total_disk_space - $total_size;
	my $space_to_free_up = $required_unused_space - $available_disk_space;
	print("Filesystem takes      $total_size\n");
	print("Total space is        $total_disk_space\n");
	print("Available disk space: $available_disk_space\n");
	print("Required disk space:  $required_unused_space\n");
	print("Space to free up:     $space_to_free_up\n");
	my $result = findSmallestDirToDelete($space_to_free_up);
	print("==========================================\n");
	print("[Solve2] Result:      $result\n");
}

read_input_file($input_filename);
computeStructure();

solve1();
print("\n\n");
solve2();

