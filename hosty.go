package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

const (
	prefix     string = "#hosty-"
	hostsFile  string = "/etc/hosts"
	comment    string = "#"
	whitespace string = " "
	lineBreak  string = "\n"
	empty      string = ""
	enabled    string = "✔"
	disabled   string = "✖"
)

type printer func(a ...interface{}) (n int, err error)

func init() {
	flag.Usage = func() {
		fmt.Println("Usage:\n")
		fmt.Println("\thosty [command [arguments]]\n")
		fmt.Println("The commands are:\n")
		fmt.Println("\tcat, c         Echo all /etc/hosts content")
		fmt.Println("\t\thosty cat\n")
		fmt.Println("\tsave, s        Save an entry, use this to create or edit an entry")
		fmt.Println("\t\thosty save example-entry 127.0.0.1 example.com\n")
		fmt.Println("\tenable, e      Enable an entry")
		fmt.Println("\t\thosty enable example-entry\n")
		fmt.Println("\tdisable, d     Disable an entry")
		fmt.Println("\t\thosty disable example-entry\n")
		fmt.Println("\tremove, r      Remove an entry")
		fmt.Println("\t\thosty remove example-entry\n")
		os.Exit(0)
	}
}

func main() {
	fileContent := read()
	entries := parseEntries(fileContent)

	flag.Parse()

	cmd := flag.Arg(0)

	if cmd == empty {
		list(entries, fmt.Print)
		os.Exit(0)
	}

	switch cmd {
	case "cat", "c":
		fmt.Println(fileContent)
	case "save", "s":
		if len(flag.Args()) < 4 {
			flag.Usage()
		}
		entry := flag.Arg(1)
		ip := flag.Arg(2)
		domains := strings.Trim(strings.Join(flag.Args()[3:], whitespace), whitespace)

		newLine := save(fileContent, entries, entry, ip, domains)

		entries[entry] = newLine

		list(entries, fmt.Print)
	case "enable", "e":
		entry := flag.Arg(1)
		toggle(fileContent, entries, entry, comment, whitespace)
	case "disable", "d":
		entry := flag.Arg(1)
		toggle(fileContent, entries, entry, whitespace, comment)
	case "remove", "r":
		entry := flag.Arg(1)
		if line, hasEntry := entries[entry]; hasEntry {
			fileContent = strings.Replace(fileContent, prefix+entry+lineBreak, empty, 1)
			fileContent = strings.Replace(fileContent, line+lineBreak, empty, 1)

			write(fileContent)

			delete(entries, entry)

			list(entries, fmt.Print)
		} else {
			fmt.Println("hosty has no entry: " + entry)
			flag.Usage()
		}
	}

	os.Exit(0)
}

// list prints pretty entries output
func list(entries map[string]string, print printer) {
	if len(entries) > 0 {
		output := []string{"hosty entries:\n"}
		var keys []string
		for k := range entries {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			line := entries[k]
			status := enabled
			if strings.HasPrefix(line, comment) {
				status = disabled
			}

			output = append(output, fmt.Sprintf("%s %s\t%s\n", status, k, line))
		}
		print(strings.Join(output, ""))
	} else {
		print("hosty has no entries!\n")
		flag.Usage()
	}
}

// write fileContent to hostsFile
func write(fileContent string) {
	var err = ioutil.WriteFile(hostsFile, []byte(fileContent), 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//TODO toggle should be self contained about how char to replace
// toggle change entry's status from enabled to disabled and the other way around
func toggle(fileContent string, entries map[string]string, entry string, current string, replacer string) {
	line := entries[entry]
	if strings.HasPrefix(line, current) {
		newLine := strings.Replace(line, current, replacer, 1)
		entries[entry] = newLine
		fileContent = strings.Replace(fileContent, line, newLine, 1)
		write(fileContent)
	}
	list(entries, fmt.Print)
}

// read hostsFile
// return fileContent as string
func read() string {
	fileBytes, err := ioutil.ReadFile(hostsFile)
	if err != nil {
		panic(err)
		os.Exit(1)
	}

	return string(fileBytes)
}

// parse hostsFile content and put managed entries in a map
// return entries' map[string]string
func parseEntries(fileContent string) map[string]string {
	entries := make(map[string]string)

	lines := strings.Split(fileContent, lineBreak)
	for index, line := range lines {
		if strings.HasPrefix(line, prefix) {
			entry := strings.Replace(line, prefix, empty, -1)
			nextLineIndex := index + 1
			entries[entry] = lines[nextLineIndex]
		}
	}

	return entries
}

// update new entry if it already exists
// save new entry to fileContent
// invoke write
func save(fileContent string, entries map[string]string, entry string, ip string, domains string) string {
	newLine := ip + whitespace + domains
	if line, hasEntry := entries[entry]; hasEntry {
		// replacing an existing line will enable it by default
		newLine = whitespace + newLine

		fileContent = strings.Replace(fileContent, line, newLine, 1)
	} else {
		// new entry will be enabled by default
		newLine = whitespace + newLine

		fileContent += prefix + entry + lineBreak
		fileContent += newLine + lineBreak
	}

	write(fileContent)

	return newLine
}
