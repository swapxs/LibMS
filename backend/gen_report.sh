#!/usr/bin/env sh

set -e

BLOB_DIRECTORY="./.blob"
REPORT_FILE="$BLOB_DIRECTORY/report.out"
CSV_FILE="$BLOB_DIRECTORY/report.csv"
TEXT_REPORT="$BLOB_DIRECTORY/report_summary"
FINAL_REPORT="$BLOB_DIRECTORY/final_report"

printf "\n=============================================================\n"
printf "Running Tests with Coverage\n"
printf "=============================================================\n\n"

go test -v ./... -coverprofile="$REPORT_FILE"

printf "\n=============================================================\n"
printf "Generating Test Coverage Report\n"
printf "=============================================================\n\n"

go tool cover -func="$REPORT_FILE" > "$TEXT_REPORT"

cat "$TEXT_REPORT"

awk '
    BEGIN { print "File,Coverage" > "'"$CSV_FILE"'"; }
    {
        if ($1 ~ /\.go:/) {
            split($1, arr, ":");
            file = arr[1];
            percentage = substr($3, 1, length($3)-1);
            total[file] += percentage;
            count[file] += 1;
        }
    }
    END {
        for (file in total) {
            avg_coverage = total[file] / count[file];
            printf "%-50s %.2f%%\n", file, avg_coverage >> "'"$FINAL_REPORT"'";
            print file "," avg_coverage "%" >> "'"$CSV_FILE"'";
        }
    }
' "$TEXT_REPORT"

printf "\n=============================================================\n"
printf "Consolidated Coverage Per Handler File\n"
printf "=============================================================\n\n"

cat "$FINAL_REPORT"

rm "$REPORT_FILE" "$TEXT_REPORT" "$FINAL_REPORT"

printf "\n=============================================================\n"
printf "Final Report successfully generated: %s\n" "$CSV_FILE"
printf "=============================================================\n\n"
