package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	rc = ".envrc"
)

var (
	evMap = map[string]rcRecord{
		"AWS_ACCESS_KEY_ID":           {"", false},
		"AWS_SECRET_ACCESS_KEY":       {"", false},
		"AWS_DEFAULT_REGION":          {"", false},
		"AWS_CA_BUNDLE":               {"", false},
		"AWS_CONFIG_FILE":             {"", false},
		"AWS_DEFAULT_OUTPUT":          {"", false},
		"AWS_PAGER":                   {"", false},
		"AWS_PROFILE":                 {"", false},
		"AWS_ROLE_SESSION_NAME":       {"", false},
		"AWS_SESSION_TOKEN":           {"", false},
		"AWS_SHARED_CREDENTIALS_FILE": {"", false},
	}
)

type rcRecord struct {
	value string
	exist bool
}

func main() {
	read := readExistEnvrc(rc)
	if read == nil {
		if err := createEnvrc(); err != nil {
			fmt.Errorf("create failed")
			os.Exit(1)
		}
		fmt.Printf("%s created\n", rc)
		os.Exit(0)
	}

	if err := updateEnvrc(read); err != nil {
		fmt.Errorf("update failed")
		os.Exit(1)
	}
	fmt.Printf("%s updated\n", rc)
}

func printPrompt(key, val string) {
	fmt.Printf("%s[%s]: ", key, val)
}

func formatExport(key, val string) string {
	return fmt.Sprintf("export %s=%s\n", key, val)
}

func readExistEnvrc(filename string) map[string]rcRecord {
	f, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer f.Close()
	return mapFromEnvrc(f)
}

func createEnvrc() error {
	f, err := os.Create(rc)
	if err != nil {
		return err
	}
	defer f.Close()

	sc := bufio.NewScanner(os.Stdin)
	for i, v := range evMap {
		printPrompt(i, v.value)
		sc.Scan()
		data := sc.Text()
		if data != "" {
			f.WriteString(formatExport(i, data))
		}
	}

	return nil
}

func updateEnvrc(maps map[string]rcRecord) error {
	f, err := os.Create(rc)
	if err != nil {
		return err
	}
	defer f.Close()

	in := bufio.NewScanner(os.Stdin)
	for i, v := range evMap {
		if val, ok := maps[i]; ok {
			printPrompt(i, val.value)
			in.Scan()
			input := in.Text()
			if input == "" {
				input = val.value
			}
			maps[i] = rcRecord{input, true}
			f.WriteString(formatExport(i, input))
			continue
		}
		printPrompt(i, v.value)
		in.Scan()
		input := in.Text()
		if input != "" {
			f.WriteString(formatExport(i, input))
		}
	}
	for i, v := range maps {
		if !v.exist {
			f.WriteString(formatExport(i, v.value))
		}
	}
	return nil
}

func mapFromEnvrc(fp *os.File) map[string]rcRecord {
	sc := bufio.NewScanner(fp)
	res := make(map[string]rcRecord, 3)
	for sc.Scan() {
		line := sc.Text()
		v := strings.Split(strings.TrimLeft(line, "export "), "=")
		res[v[0]] = rcRecord{v[1], false}
	}
	return res
}
