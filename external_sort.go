package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strconv"
)

var debug bool

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
func read_into_buffer(input_buffer_ref *[]string, fhandler *os.File, chunk_size int) int {
	scanner := bufio.NewScanner(fhandler)
	eof := false
	if debug {
		log.Printf("len(*input_buffer_ref)=%d, chunk_size=%d\n", len(*input_buffer_ref), chunk_size)
	}
	for (len(*input_buffer_ref) < chunk_size) && !eof {
		eof = !scanner.Scan()
		if eof {
			if err := scanner.Err(); err != nil {
				log.Println(err)
			}
			break
		}
		if debug {
			log.Printf("scanner.Text()=%v\n", scanner.Text())
		}
		*input_buffer_ref = append(*input_buffer_ref, scanner.Text())
	}
	if debug {
		log.Printf("len(*input_buffer_ref)=%d\n", len(*input_buffer_ref))
	}
	return len(*input_buffer_ref)
}

func main() {
	/* I make an assumption that we measure memory in lines.
	I could set this limit in bytes, look through the input file with one more pass
	to find the longest line length but I think we can omit this here. */
	var avail_mem int
	flag.IntVar(&avail_mem, "avail_mem", 100000, "Number of lines to fit the memory")
	flag.BoolVar(&debug, "debug", false, "Pring debug info")
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		log.Fatal("Usage: " + os.Args[0] + " [-avail_mem <lines>] <input file> <output file>")
	}
	input_file := args[0]
	output_file := args[1]
	if debug {
		log.Println("DEBUG: debug is on")
	}
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
		var lines []string
		if debug {
			log.Printf("DEBUG: len(lines)=%d, avail_mem=%d\n", len(lines), avail_mem)
		}
		for len(lines) < (avail_mem - 1) {
			if debug {
				log.Printf("read %d lines\n", len(lines))
			}
			eof = !scanner.Scan()
			if eof {
				if err := scanner.Err(); err != nil {
					log.Println(err)
				}
				break
			}
			lines = append(lines, scanner.Text())
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
				outfile.WriteString(lines[i] + "\n")
			}
			outfile.Close()
			ind += 1
		}
	}
	if debug {
		log.Println("read total " + strconv.Itoa(tot_lines) + " lines")
		log.Printf("ind=%d\n", ind)
	}
	infile.Close()

	// 2. Merge the intermediate files into output file
	outfile, err := os.Create(output_file)
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()
	var fhandlers []*os.File
	for i := 0; i < ind; i += 1 {
		fname := inter_fname_patt + "_" + strconv.Itoa(i)
		fh, err := os.Open(fname)
		if err != nil {
			log.Fatal(err)
		}
		defer fh.Close()
		fhandlers = append(fhandlers, fh)
	}
	chunk_size := int((avail_mem - 1) / (ind + 1)) // -1 for current line, +1 for output buffer
	if debug {
		log.Printf("chunk_size=%d\n", chunk_size)
	}

	var output_buffer []string
	var input_buffers [][]string
	input_buffers = make([][]string, ind)
	// read into input buffers
	for i := 0; i < ind; i += 1 {
		read_into_buffer(&input_buffers[i], fhandlers[i], chunk_size)
	}

	if debug {
		log.Printf("input_buffers=%v\n", input_buffers)
	}
	empty_handlers := 0
	// top-level loop here
	for {
		// k-merge here (ind-merge)
		// find "minimum" of ind heads, put it into output buffer and pop it
		min_ind := 0
		// skip empty buffer
		for min_ind = 0; min_ind < ind; min_ind += 1 {
			if debug {
				log.Printf("min_ind=%d\n", min_ind)
				log.Printf("len(input_buffers[min_ind])=%d\n", len(input_buffers[min_ind]))
			}
			if len(input_buffers[min_ind]) > 0 {
				break
			}
		}
		if debug {
			log.Printf("min_ind=%d\n", min_ind)
		}
		for i := min_ind + 1; i < ind; i += 1 {
			// skip empty buffer
			if len(input_buffers[i]) == 0 {
				continue
			}
			// ok I won't implement my own string comparison
			if input_buffers[i][0] < input_buffers[min_ind][0] {
				min_ind = i
			}
		}
		if debug {
			log.Printf("min_ind=%d\n", min_ind)
		}
		// shift
		x := input_buffers[min_ind][0]
		input_buffers[min_ind] = input_buffers[min_ind][1:]
		output_buffer = append(output_buffer, x)
		if len(output_buffer) == chunk_size {
			// flush output buffer
			for i := 0; i < len(output_buffer); i += 1 {
				outfile.WriteString(output_buffer[i] + "\n")
			}
			output_buffer = nil
		}
		if len(input_buffers[min_ind]) == 0 {
			if 0 == read_into_buffer(&input_buffers[min_ind], fhandlers[min_ind], chunk_size) {
				empty_handlers += 1
			}
		}
		if empty_handlers == ind {
			break
		}
	}

	// final flush of out buffer
	for i := 0; i < len(output_buffer); i += 1 {
		outfile.WriteString(output_buffer[i] + "\n")
	}
	outfile.Close()

	// cleanup
	for i := 0; i < ind; i += 1 {
		fname := inter_fname_patt + "_" + strconv.Itoa(ind)
		fhandlers[i].Close()
		os.Remove(fname)
	}
}
