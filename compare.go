package main

import "log"

func compareI18nFiles(files []I18nFileInfo) int {
	log.Printf("\n[Compare]: Start\n")
	code := 0
	if len(files) < 2 {
		log.Println("Nothing to compare")
	}

	for _, file := range files {
		// browse key
		for k := range file.Content {
			for _, fileToCheck := range files {
				if fileToCheck.Name == file.Name {
					continue
				}

				if _, has := fileToCheck.Content[k]; !has {
					code = 1
					log.Printf("Missing key %s in file %s\n", k, fileToCheck.Name)
				}
			}
		}
	}
	if code == 0 {
		log.Println("[Compare]: all files are good")
	}
	return code
}
