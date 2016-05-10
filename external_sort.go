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

//func read_into_buffer(input_buffer_ref *[]string, fhandler *os.File, chunk_size int) int {
func read_into_buffer(input_buffer_ref *[]string, scanner *bufio.Scanner, chunk_size int) int {
	//scanner := bufio.NewScanner(fhandler)
	eof := false
	if debug {
		log.Printf("len(*input_buffer_ref)=%d, chunk_size=%d\n", len(*input_buffer_ref), chunk_size)
	}
	for (len(*input_buffer_ref) < chunk_size) && !eof {
		eof = !scanner.Scan()
		if debug {
			log.Printf("eof=%v\n", eof)
		}
		if eof {
			if debug {
				log.Println("DEBUG: eof")
				log.Printf("[last]=%d\n", len(*input_buffer_ref))
				if len(*input_buffer_ref) > 0 {
					log.Printf("last=%d\n", (*input_buffer_ref)[len(*input_buffer_ref)-1])
				}
			}
			if err := scanner.Err(); err != nil {
				log.Println(err)
			}
			break
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
	scanner := bufio.NewScanner(infile)
	tot_lines := 0
	ind := 0 // intermediate file name index, also number of sorted pieces for further merge
	for {
		var lines []string
		lines_read := read_into_buffer(&lines, scanner, avail_mem-1)
		if debug {
			log.Println("read " + strconv.Itoa(lines_read) + " lines")
			log.Printf("ind=%d\n", ind)
		}
		if lines_read == 0 {
			break
		}
		tot_lines += lines_read
		mysort(lines, 0, lines_read-1)
		out_fname := inter_fname_patt + "_" + strconv.Itoa(ind)
		outfile, err := os.Create(out_fname)
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < lines_read; i += 1 {
			outfile.WriteString(lines[i] + "\n")
		}
		outfile.Close()
		ind += 1
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
	//defer outfile.Close()
	var fhandlers []*os.File
	var scanners []*bufio.Scanner
	for i := 0; i < ind; i += 1 {
		fname := inter_fname_patt + "_" + strconv.Itoa(i)
		fh, err := os.Open(fname)
		if err != nil {
			log.Fatal(err)
		}
		scanner := bufio.NewScanner(fh)
		scanners = append(scanners, scanner)
		//defer fh.Close()
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
		read_into_buffer(&input_buffers[i], scanners[i], chunk_size)
	}

	empty_handlers := 0
	// top-level loop here
	for {
		// k-merge here (ind-merge)
		// find "minimum" of ind heads, put it into output buffer and pop it
		min_ind := 0
		// skip empty buffer
		for min_ind = 0; min_ind < ind; min_ind += 1 {
			if len(input_buffers[min_ind]) > 0 {
				break
			}
			if debug {
				log.Printf("skip empty buf min_ind=%d\n", min_ind)
			}
		}
		for i := min_ind + 1; i < ind; i += 1 {
			// skip empty buffer
			if len(input_buffers[i]) == 0 {
				if debug {
					log.Printf("skip empty buf i=%d\n", i)
				}
				continue
			}
			// ok I won't implement my own string comparison
			if input_buffers[i][0] < input_buffers[min_ind][0] {
				min_ind = i
			}
		}
		// shift
		x := input_buffers[min_ind][0]
		input_buffers[min_ind] = input_buffers[min_ind][1:]
		// push
		output_buffer = append(output_buffer, x)
		if len(output_buffer) == chunk_size {
			// flush output buffer
			for i := 0; i < len(output_buffer); i += 1 {
				outfile.WriteString(output_buffer[i] + "\n")
			}
			output_buffer = nil
		}
		if len(input_buffers[min_ind]) == 0 {
			if debug {
				log.Printf("reading into buffer %d\n", min_ind)
			}
			if read_into_buffer(&input_buffers[min_ind], scanners[min_ind], chunk_size) == 0 {
				if debug {
					log.Printf("EMPTY\n")
				}
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
		fname := inter_fname_patt + "_" + strconv.Itoa(i)
		if debug {
			log.Printf("DEBUG: fname=%s\n", fname)
		}
		fhandlers[i].Close()
		err = os.Remove(fname)
		if err != nil {
			log.Println(err)
		}
	}
}
