Very large file sorting

Hint: External sorting is a term for a class of sorting algorithms that can handle massive amounts of data.

Problem: You have a *very* large text file, so it does not fit in memory, with text lines. Sort the file into an output file where all the lines are sorted in alphabetic order, taking into account all words per line. The lines themselves do not need to be sorted and are not to be modified. Lines are considered to be average in length so edge cases such as a file with just two very large lines should still work but it is OK if performance suffers in that case.

Boundary: Use any programming language you feel comfortable with. Please use standard libraries only, no batch or stream processing frameworks. Be as efficient as possible while avoiding using standard library sorting routines. Provide a rationale for your approach. Design schemas are welcome.

Please note that the file.txt that we use to measure the performance of the result is generated via: ruby -e 'a=STDIN.readlines;5000000.times do;b=[];16.times do; b << a[rand(a.size)].chomp end; puts b.join(" "); end' < /usr/share/dict/words > file.txt

