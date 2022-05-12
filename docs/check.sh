#!/bin/bash

perl -ne '
$cnt=1;

if (m#(/ch/)\d+(/.*)#) {
	$cnt=0;

	$a=$1;
	$b=$2;
	foreach $v (1 .. 32)
	{
		printf("%s%.2d%s\n", $a, $v, $b);
	}

} elsif (m#(/auxin/)\d+(/.*)#) {
	$cnt=0;

	$a=$1;
	$b=$2;
	foreach $v (1 .. 6)
	{
		printf("%s%.2d%s\n", $a, $v, $b);
	}

} elsif (m#(/\w+/)\d+(/.*)#) {
	$cnt=0;

	$a=$1;
	$b=$2;
	foreach $v (1 .. 16)
	{
		printf("%s%.2d%s\n", $a, $v, $b);
	}
}

if (m#(/ch/\d+/\w+/)\d+(/.*)#) {
	$cnt=0;

	$a=$1;
	$b=$2;
	foreach $v (1 .. 16)
	{
		printf("%s%.2d%s\n", $a, $v, $b);
	}

} elsif (m#(/auxin/\d+/\w+/)\d+(/.*)#) {
	$cnt=0;

	$a=$1;
	$b=$2;
	foreach $v (1 .. 16)
	{
		printf("%s%.2d%s\n", $a, $v, $b);
	}

} elsif (m#(/bus/\d+/\w+/)\d+(/.*)#) {
	$cnt=0;

	$a=$1;
	$b=$2;
	foreach $v (1 .. 16)
	{
		printf("%s%.2d%s\n", $a, $v, $b);
	}
}

if ($cnt == 1) {
	print
}
' $1 | sort -u

