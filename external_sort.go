package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strconv"
)

/* okay make my own sort to avoid using std one
accepts list reference to avoid list copy
and say let it be classic quicksort */
func partition(array_ref []string, lo int, hi int) int {
	pivot := array_ref[hi]
	i := lo
	for j := lo; j < hi; j++ {
		if array_ref[j] <= pivot {
			// swap A[i] with A[j]
			array_ref[i], array_ref[j] = array_ref[j], array_ref[i]
			i += 1
		}
	}
	// swap A[i] with A[hi]
	array_ref[i], array_ref[hi] = array_ref[hi], array_ref[i]
	return i
}
func mysort(array_ref []string, lo int, hi int) {
	if len(array_ref) == 0 {
		return
	}
	if lo < hi {
		p := partition(array_ref, lo, hi)
		mysort(array_ref, lo, p-1)
		mysort(array_ref, p+1, hi)
	}
}

func main() {
	/* I make an assumption that we measure memory in lines.
	I could set this limit in bytes, look through the input file with one more pass
	to find the longest line length but I think we can omit this here. */
	var avail_mem int
	var input_file, output_file string
	flag.IntVar(&avail_mem, "avail_mem", 100000, "Number of lines to fit the memory")
	flag.StringVar(&input_file, "input_file", "", "Input file path")
	flag.StringVar(&output_file, "output_file", "", "Output file path")
	flag.Parse()
	inter_fname_patt := "mysorted" // pattern to name intermediate files
	infile, err := os.Open(input_file)
	if err != nil {
		log.Fatal(err)
	}
	//defer infile.Close()
	/*
	* 1. Read input files into chunks fitting available memory, sort the chunks
	* using mysort() and dump them to set of intermediate files
	 */
	tot_lines := 0
	eof := false
	ind := 0 // intermediate file name index, also number of sorted pieces for further merge
	scanner := bufio.NewScanner(infile)
	for !eof {
		lines := make([]string, avail_mem)
		for len(lines) < avail_mem {
			eof = !scanner.Scan()
			if eof {
				if err := scanner.Err(); err != nil {
					log.Println(err)
				}
				break
			}
			lines[len(lines)] = scanner.Text()
		}
		tot_lines += len(lines)
		if len(lines) > 0 {
			mysort(lines, 0, len(lines)-1)
			out_fname := inter_fname_patt + "_" + strconv.Itoa(ind)
			outfile, err := os.Create(out_fname)
			if err != nil {
				log.Fatal(err)
			}
			for i := 0; i < len(lines); i += 1 {
				outfile.WriteString(lines[i])
			}
			outfile.Close()
			ind += 1
		}
	}
	log.Println("read total " + strconv.Itoa(tot_lines) + " lines")
	infile.Close()
	/*
	   ## 2. Merge the intermediate files into output file
	   open (OUT, ">$output_file") or die "Could not create output file '$output_file': $!";
	   my @fhandlers;
	   for (my $i = 0; $i < $ind; ++$i) {
	     my $fname = "${inter_fname_patt}_$i";
	     open ($fhandlers[$i], $fname) or die "Could not open intermediate file '$fname': $!";
	   }

	   my $chunk_size = int($avail_mem / ($ind + 1)); ## +1 for output buffer
	   print "chunk_size='$chunk_size'\n";
	   my @output_buffer;
	   my @input_buffers;

	   sub read_into_buffer {
	     my ($i, $input_buffers_ref, $fhandler, $chunk_size) = @_;
	     my @lines;
	     my $line;
	     while ((@lines < $chunk_size) and defined ($line = <$fhandler>)) {
	       push @lines, $line;
	     }
	     $input_buffers_ref->[$i] = \@lines;
	     return int(@lines);
	   }

	   ## read into input buffers
	   for (my $i = 0; $i < $ind; ++$i) {
	     read_into_buffer($i, \@input_buffers, $fhandlers[$i], $chunk_size);
	   }
	   my $empty_handlers = 0;
	   ## top-level loop here
	   do {
	     ## k-merge here (ind-merge)
	     ## find "minimum" of ind heads, put it into output buffer and pop it
	     my $min_ind;
	     ## skip empty buffer
	     for ($min_ind = 0; $min_ind < $ind; ++$min_ind) {
	       if (@{$input_buffers[$min_ind]}) {
	         last;
	       }
	     }
	     for (my $i = $min_ind + 1; $i < $ind; ++$i) {
	       ## skip empty buffer
	       next unless (@{$input_buffers[$i]});
	       ## ok I won't implement my own string comparison
	       if ($input_buffers[$i][0] lt $input_buffers[$min_ind][0]) {
	         $min_ind = $i;
	       }
	     }
	     push @output_buffer, shift @{$input_buffers[$min_ind]};
	     if (@output_buffer == $chunk_size) {
	       ## flush output buffer
	       print OUT @output_buffer;
	       @output_buffer = ();
	     }
	     unless (@{$input_buffers[$min_ind]}) {
	       unless (read_into_buffer($min_ind, \@input_buffers, $fhandlers[$min_ind], $chunk_size)) {
	         ++$empty_handlers;
	       }
	     }
	   } while ($empty_handlers < $ind);

	   ## final flush of out buffer
	   print OUT @output_buffer;
	   close OUT;

	   ## cleanup
	   for (my $i = 0; $i < $ind; ++$i) {
	     my $fname = "${inter_fname_patt}_$i";
	     close ($fhandlers[$i]) or die "Could not close intermediate file '$fname': $!";
	     unlink $fname or die "Could not delete intermediate file '$fname': $!";
	   }
	*/
}
