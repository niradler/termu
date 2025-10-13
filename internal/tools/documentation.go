package tools

var ToolDocumentation = map[string]string{
	"sd": `sd - Search and replace tool (modern sed alternative)

Usage:
  sd 'old_text' 'new_text' [file...]
  sd -p 'old' 'new' file.txt  # Preview changes
  
Examples:
  sd 'var ' 'const ' file.js
  fd -e js | xargs sd 'old' 'new'`,

	"fd": `fd - Fast file finder (modern find alternative)

Usage:
  fd [pattern] [path]
  fd -e [extension] [pattern]
  fd -t [type] [pattern]
  
Examples:
  fd pattern                    # Find by name
  fd -e py                      # Find Python files
  fd --changed-within 7d        # Files modified in last 7 days
  fd -x rm {}                   # Execute command on results`,

	"rg": `ripgrep (rg) - Fast content search (modern grep alternative)

Usage:
  rg [pattern] [path]
  rg -t [type] [pattern]
  rg -C [num] [pattern]
  
Examples:
  rg "pattern"                  # Basic search
  rg -t py "import"             # Search Python files
  rg -C 3 "error"               # Show 3 lines context
  rg -i "pattern"               # Case insensitive
  rg --files-with-matches       # Show only filenames`,

	"bat": `bat - File preview with syntax highlighting (modern cat alternative)

Usage:
  bat [file...]
  bat -n [file]                 # Show line numbers
  bat -r [start]:[end] [file]   # Show line range
  
Examples:
  bat file.py                   # Preview with highlighting
  bat -A file.txt               # Show all characters
  bat -p file.md                # Plain output`,

	"xsv": `xsv - CSV manipulation toolkit

Usage:
  xsv [command] [args]
  
Commands:
  stats data.csv                # Show statistics
  select name,age data.csv      # Select columns
  search -s name "John" data.csv # Search
  sort -s age data.csv          # Sort by column
  count data.csv                # Count rows`,

	"jaq": `jaq - JSON processor (jq clone)

Usage:
  jaq [filter] [file]
  
Examples:
  jaq '.' data.json             # Parse JSON
  jaq '.items[] | .name'        # Extract field
  jaq 'select(.price > 100)'    # Filter
  jaq -r '.name'                # Raw output`,

	"yq": `yq - YAML/JSON processor

Usage:
  yq [expression] [file]
  yq -o [format] [file]
  
Examples:
  yq '.key' config.yaml         # Read value
  yq -o json config.yaml        # Convert to JSON
  yq '.key = "value"' -i file   # Update in-place`,

	"dua": `dua - Disk usage analyzer

Usage:
  dua [path]
  dua i                         # Interactive mode
  dua aggregate                 # Show top directories
  
Examples:
  dua .                         # Analyze current directory
  dua i /path                   # Interactive exploration`,

	"eza": `eza - Modern directory listing (ls alternative)

Usage:
  eza [options] [path]
  
Examples:
  eza -l                        # Long listing
  eza --tree                    # Tree view
  eza --git -l                  # Show git status
  eza -l --sort size            # Sort by size
  eza -a                        # Show hidden files`,
}

var ModernToolMap = map[string]string{
	"find": "fd",
	"grep": "rg",
	"cat":  "bat",
	"ls":   "eza",
	"sed":  "sd",
	"du":   "dua",
}

func GetDocumentation(tool string) string {
	if doc, ok := ToolDocumentation[tool]; ok {
		return doc
	}
	return ""
}

func GetModernAlternative(tool string) string {
	if modern, ok := ModernToolMap[tool]; ok {
		return modern
	}
	return tool
}
