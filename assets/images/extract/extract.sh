#!/bin/bash

rm out?.png
perl -e '
$ly = $y = 22;
foreach $yi (0 .. 6)
{
	$lx = $x = 22;
	foreach $xi (0 .. 10)
	{
		# $i = ($yi * 10) + $xi;
		$i++;

		printf"convert OriginalIcons1.png -crop 91x58+%d+%d Icon%d.png\n", $x, $y, $i;

		$x = $lx + 2 + 91; $lx = $x;
	}

	$y = $ly + 58; $ly = $y;
}
' | tee /tmp/z;
bash /tmp/z

# 3d3d3d


# 22,  22
# 113, 80

# 115, 80
# 206, 138

# 208, 138
# 299, 196

# 394, 370
# 485, 428

# 487, 370
# 578, 428
