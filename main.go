package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var srcDirectory string
var i18nDirectory string
var mod string

func init() {
	flag.StringVar(&i18nDirectory, "i18nDirectory", ".", "Directory of translation files")
	flag.StringVar(&srcDirectory, "srcDirectory", "", "Source directory where files that use translation can be found")
	flag.StringVar(&mod, "mod", "compare", "Execution mode:   compare, unused, missing")
}

func main() {
	log.SetFlags(0)
	flag.Parse()
	fmt.Printf("mod: %s\n", mod)
	fmt.Printf("i18nDirectory: %s\n", i18nDirectory)
	fmt.Printf("srcDirectory: %s\n", srcDirectory)

	if i18nDirectory == "" {
		log.Fatalf("i18nDirectory is mandoratory")
	}
	if (mod == "unused" || mod == "all") && srcDirectory == "" {
		log.Fatalf("srcDirectory is mandoratory with mod %s", mod)
	}

	i18nfiles := readi18nFiles()

	if len(i18nfiles) == 0 {
		log.Printf("No i18n files found in %s\n", i18nDirectory)
	}
	switch mod {
	case "compare":
		os.Exit(compareI18nFiles(i18nfiles))
	case "unused":
		os.Exit(findUnusedTranslate(i18nfiles[0]))
	case "missing":
		os.Exit(findMissingTranslate(i18nfiles[0]))
	case "all":
		returnCode := compareI18nFiles(i18nfiles) + findUnusedTranslate(i18nfiles[0]) + findMissingTranslate(i18nfiles[0])
		os.Exit(returnCode)
	default:
		log.Fatalf("Unknown mod %s", mod)
	}

}

func readi18nFiles() []I18nFileInfo {
	files, err := ioutil.ReadDir(i18nDirectory)
	if err != nil {
		log.Fatalf("readi18nFiles: unable to list files in i18nDirectory: %v", err)
	}

	i18Files := make([]I18nFileInfo, 0, len(files))
	for _, fi := range files {
		if fi.IsDir() || !strings.HasSuffix(fi.Name(), ".json") {
			continue
		}

		filePath := i18nDirectory + "/" + fi.Name()
		fileContentByte, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatalf("readi18nFiles: unable to read file %s: %v", filePath, err)
		}
		var fileContent map[string]string
		if err := json.Unmarshal(fileContentByte, &fileContent); err != nil {
			log.Fatalf("readi18nFiles: unable to unmarshal file %s: %v", filePath, err)
		}
		i18Files = append(i18Files, I18nFileInfo{
			Name:      fi.Name(),
			Directory: i18nDirectory,
			Content:   fileContent,
		})
	}
	return i18Files
}
