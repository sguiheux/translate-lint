package main

import (
	"io/ioutil"
	"log"
	"strings"
)

func findUnusedTranslate(file I18nFileInfo) int {
	log.Printf("\n[Unused]: Start\n")
	contentData := make(map[string]string, len(file.Content))
	for k, v := range file.Content {
		contentData[k] = v
	}
	findUnusedInDirectory(srcDirectory, contentData)

	if len(contentData) == 0 {
		return 0
	}

	log.Printf("%d Unused translation : \n", len(contentData))
	for k := range contentData {
		log.Println(k)
	}
	return 1
}

func findUnusedInDirectory(dir string, i18nContent map[string]string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("findUnusedInDirectory: unable to read directory %s: %v", srcDirectory, err)
	}
	for _, f := range files {
		if f.IsDir() {
			findUnusedInDirectory(dir+"/"+f.Name(), i18nContent)
			continue
		}
		if !strings.HasSuffix(f.Name(), ".html") && !strings.HasSuffix(f.Name(), ".ts") {
			continue
		}

		contentBytes, err := ioutil.ReadFile(dir + "/" + f.Name())
		if err != nil {
			log.Fatalf("findUnusedInDirectory: unable to read file %s in %s: %v", f.Name(), dir, err)
		}

		contentString := string(contentBytes)
		for k := range i18nContent {
			if strings.Contains(contentString, k) {
				delete(i18nContent, k)
			}
		}

	}
}
