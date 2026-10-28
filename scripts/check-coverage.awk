BEGIN {
	total_found = 0
	threshold = target
	gsub(/%/, "", threshold)
}

$1 == "total:" {
	coverage = $3
	gsub(/%/, "", coverage)
	total_found = 1
}

END {
	if (!total_found) {
		print "go tool cover did not report total coverage" > "/dev/stderr"
		exit 1
	}

	printf "Total coverage: %s%%\n", coverage
	if (threshold != "" && coverage + 0 < threshold + 0) {
		printf "Coverage %s%% is below the required %s%%\n", coverage, threshold > "/dev/stderr"
		exit 1
	}
}
