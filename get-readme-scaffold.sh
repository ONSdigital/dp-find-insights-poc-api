#!/bin/bash

thisDir=$PWD
outfile="${thisDir}/readme-scaffold.md"

# prep outfile
projName=$(basename "$thisDir")
printf "# ${projName}\n" > "$outfile"

# loop through all non-hidden folders and try 'go doc --all' in them, write output to outfile if meaningful
find . -type d -not -path '*/\.*' | while read -r dir
do
    printf "\n## dir: ${dir}\n" >> "$outfile"
    # print autodocs for any go code
    outContent=$(go doc -all "$dir" 2>/dev/null) && printf "### public Go code:\n\`\`\` ${outContent} \n\`\`\`\n" >> "${outfile}"
    # link any readmes
    cat "${dir}/README.md" &>/dev/null && printf "See [README](${dir}/README.md) for more details.\n" >> "${outfile}"

done
