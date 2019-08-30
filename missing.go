package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

var translateTSRegexp = regexp.MustCompile("instant\\('(.*?)'")
var translateHtmlTemplate = regexp.MustCompile("(?:{{|\")\\s*'(.*?)'\\s*\\|\\s*translate\\s*(?:}}|\")")

func findMissingTranslate(file I18nFileInfo) int {
	log.Printf("\n[Missing]: Start\n")
	result := findMissingInDirectory(srcDirectory, file.Content)
	if len(result) == 0 {
		return 0
	}
	log.Printf("%d Missing translation : \n", len(result))
	for _, r := range result {
		for _, m := range r.Translate {
			log.Printf("%s/%s  %s", r.Directory, r.File, m)
		}
	}
	return 1
}

func findMissingInDirectory(dir string, i18nContent map[string]string) Result {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("findMissingInDirectory: unable to read directory %s: %v", srcDirectory, err)
	}
	result := make([]Missing, 0)
	for _, f := range files {
		if f.IsDir() {
			r := findMissingInDirectory(dir+"/"+f.Name(), i18nContent)
			if len(r) > 0 {
				result = append(result, r...)
			}
			continue
		}
		if !strings.HasSuffix(f.Name(), ".html") && !strings.HasSuffix(f.Name(), ".ts") {
			continue
		}

		file, err := os.Open(dir + "/" + f.Name())
		if err != nil {
			log.Fatalf("findMissingInDirectory: unable to read file %s in %s: %v", f.Name(), dir, err)
		}

		missings := make([]string, 0)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if strings.HasSuffix(f.Name(), ".html") {
				ms := findInHtmlFile(scanner.Text(), i18nContent)
				if len(ms) == 0 {
					continue
				}
				missings = append(missings, ms...)
			}
			if strings.HasSuffix(f.Name(), ".ts") {
				ms := findInTsFile(scanner.Text(), i18nContent)
				if len(ms) == 0 {
					continue
				}
				missings = append(missings, ms...)
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		file.Close()

		if len(missings) == 0 {
			continue
		}
		result = append(result, Missing{
			Translate: missings,
			File:      f.Name(),
			Directory: dir,
		})
	}
	return result
}

func findInHtmlFile(line string, i18nContent map[string]string) []string {
	missings := make([]string, 0)
	results := translateHtmlTemplate.FindStringSubmatch(line)
	if len(results) > 1 {
		for i := 1; i < len(results); i = i + 2 {
			if strings.HasSuffix(results[i], "_") {
				continue
			}
			if _, has := i18nContent[results[i]]; !has {
				missings = append(missings, results[i])
			}
		}
	}
	return missings
}

func findInTsFile(line string, i18nContent map[string]string) []string {
	missings := make([]string, 0)
	results := translateTSRegexp.FindStringSubmatch(line)
	if len(results) > 1 {
		for i := 1; i < len(results); i = i + 2 {
			if strings.HasSuffix(results[i], "_") {
				continue
			}
			if _, has := i18nContent[results[i]]; !has {
				missings = append(missings, results[i])
			}
		}
	}
	return missings
}

type Result []Missing

type Missing struct {
	Translate []string
	File      string
	Directory string
}
