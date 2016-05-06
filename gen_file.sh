#!/bin/bash
ruby -e 'a=STDIN.readlines;500000.times do;b=[];16.times do; b << a[rand(a.size)].chomp end; puts b.join(" "); end' < /usr/share/dict/words > file.txt
ruby -e 'a=STDIN.readlines;2.times do;b=[];10000.times do; b << a[rand(a.size)].chomp end; puts b.join(" "); end' < /usr/share/dict/words > long.txt
ruby -e 'a=STDIN.readlines;10.times do;b=[];10.times do; b << a[rand(a.size)].chomp end; puts b.join(" "); end' < /usr/share/dict/words > small.txt
