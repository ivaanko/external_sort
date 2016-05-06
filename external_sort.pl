#!/usr/bin/env perl
use strict;
use warnings;

use Getopt::Long;

if (@ARGV < 2) {
  print "USAGE: $0 [-avail_mem <lines>] <input file> <output file>\n";
  exit 1;
}

my $avail_mem = 100000;
my $inter_fname_patt = "mysorted"; ## pattern to name intermediate files
GetOptions("avail_mem=i" => \$avail_mem);
my ($input_file, $output_file) = @ARGV;

## okay make my own sort to avoid using std one
## accepts list reference to avoid list copy
sub mysort ($) {
  ## placeholder
  @{$_[0]} = reverse sort @{$_[0]};
}

## 1. Read input files into chunks fitting available memory, sort the chunks
##    using mysort() and dump them to set of intermediate files
open (IN, $input_file) or die "Could not open input file '$input_file': $!";
my $tot_lines = 0;
my $eof = 0;
my $ind = 0;
while (not $eof) {
  my @lines; ## I believe here I empty the buffer
  my $line;
  while ((@lines < $avail_mem) and defined ($line = <IN>)) {
    push @lines, $line;
  }
  $tot_lines += @lines;
  $eof = 1 unless (defined $line);
  if (@lines) {
    mysort (\@lines);
    my $out_fname = "${inter_fname_patt}_$ind";
    open (OUT, ">$out_fname") or die "Could not open output file '$out_fname': $!";
    print OUT @lines;
    close OUT;
    ++$ind;
  }
}
print "read total $tot_lines lines\n";
close IN;

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

close OUT;

## cleanup
for (my $i = 0; $i < $ind; ++$i) {
  my $fname = "${inter_fname_patt}_$i";
  close ($fhandlers[$i]) or die "Could not close intermediate file '$fname': $!";
  unlink $fname or die "Could not delete intermediate file '$fname': $!";
}
