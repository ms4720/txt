//Command txt is a templating language for shell programming.
//
//Input
//
//The input to the template comes from stdin.
//It is parsed in one of four ways.
//
//The default is to split stdin into records and fields, using the -R and -F
//flags respectively, similar to awk(1), and dot is set to a list of lists
//of string.
//
//If the -csv flag, or the -header flag, is specified, stdin is treated as
//a CSV file, as recognized by the encoding/csv package.
//If the -header flag is not specified, the first record is used as the header.
//Dot is set to a list of maps, with the header for each column as the key.
//
//If the -json flag is specified, stdin is treated as JSON.
//Dot is set as the decoded JSON.
//
//If the -no-stdin flag is specified, stdin is not read.
//Dot is not set.
//
//Templates
//
//The templating language is documented at http://golang.org/pkg/text/template
//with the single difference that if the first line at the top of the file
//begins with #! that line is skipped.
//If the -html flag is used, escaping functions are automatically added to all
//outputs based on context.
//
//Any command line arguments after the flags are treated as filenames
//of templates.
//The templates are named after the basename of the respective filename.
//The first file listed is the main template, unless the -template flag
//specifies otherwise.
//If the -e flag is used to define an inline template, it is always the main
//template, and the -template flag is illegal.
//
//Functions
//
//Built in functions are documented at
//http://golang.org/pkg/text/template#hdr-Functions
//
//The following additional functions are defined:
//
//	readCSV headerspec filename
//		headerspec is a comma-separated list of headers or "" to use the headers
//		in filename.
//		Dot is set to the contents of the CSV file as with -csv.
//		If the file cannot be opened or its contents are malformed, execution
//		stops.
//
//	readJSON filename
//		Read the JSON encoded file into dot or halt execution if decoding fails
//		or the file cannot be opened.
//		Dot is set to the contents of the JSON file as with -json.
//
//	read RS FS filename
//		Read filename with the default record and file splitting as specified
//		by the RS and FS regular expressions.
//
//	quoteCSV string
//		Apply the appropriate CSV quoting rules to string.
//
//	toJSON what
//		Encode what as JSON. Execution halts if
//		http://golang.org/pkg/encoding/json/#Marshal errors.
//
//	readFile filename
//		Read filename completely as a single string.
//		Execution halts if the file cannot be read.
//
//	equalFold string-one string-two
//		Reports whether the UTF-8 encoded string-one and string-two are equal
//		under Unicode case-folding.
//
//	fields string
//		Split string around whitespace.
//
//	join separator strings
//		Join the list in strings by the string separator.
//
//	lower string
//		Lowercase string.
//
//	upper string
//		Uppercase string.
//
//	title string
//		Titlecase string.
//
//	trim cutset string
//		Return string with all leading and trailing runes in cutset removed.
//
//	trimLeft cutset string
//		Return string with all leading runes in cutset removed.
//
//	trimRight cutset string
//		Return string with all trailing runes in cutset removed.
//
//	trimPrefix prefix string
//		Return string with prefix removed.
//
//	trimSuffix suffix string
//		Return string with suffix removed.
//
//	trimSpace string
//		Return string with all leading and trailing whitespace removed.
//
//	match pattern string
//		Return whether string matches the regex in pattern.
//		Execution halts if pattern is not a valid regular expression.
//
//	find pattern string
//		Returns all substrings of string that match pattern.
//		Execution halts if pattern is not a valid regular expression.
//
//	replace pattern spec string
//		Replace all substrings in string matching pattern by spec.
//		Execution halts if pattern is not a valid regular expression.
//
//	split pattern string
//		Split string into a list of substrings separated by pattern.
//		Execution halts if pattern is not a valid regular expression.
//
//	env key
//		Returns the environment variable key or "".
//
//	exec name args*
//		Execute command name with args. Stdin is nil.
//		Stderr shares the stderr of txt(1).
//		Stdout is returned as a string.
//
//	pipe name args* input
//		Execute command name with args with input as stdin.
//		Otherwise, like exec.
package main
