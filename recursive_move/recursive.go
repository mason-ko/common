package recursive_move

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

func MoveRecursive(rootDir, targetDir string, isRoot bool) {
	dir, err := ioutil.ReadDir(targetDir)
	if err != nil {
		panic(err)
	}

	var checkSize int64 = 0
	checkFile := ""

	for _, v := range dir {
		if v.IsDir() {
			fmt.Println(v.Name())
			MoveRecursive(rootDir, targetDir+"/"+v.Name(), false)
		} else if !isRoot {
			if checkSize < v.Size() {
				checkSize = v.Size()
				checkFile = v.Name()
			}
		}
	}

	if checkSize != 0 {
		//move
		os.Rename(targetDir+"/"+checkFile, rootDir+"/"+checkFile)
	}
}

func MoveSubtitle(targetDir string, movieKey, smiKey []string) {
	dir, err := ioutil.ReadDir(targetDir)
	if err != nil {
		panic(err)
	}

	var movies []string
	var smis []string

	for _, v := range dir {
		if v.IsDir() {
			continue
		}

		name := v.Name()
		ok := true
		for _, k := range movieKey {
			if !strings.Contains(name, k) {
				ok = false
			}
		}
		if ok {
			movies = append(movies, name)
		}

		ok = true
		for _, k := range smiKey {
			if !strings.Contains(name, k) {
				ok = false
			}
		}
		if ok {
			smis = append(smis, name)
		}
	}

	sort.Strings(movies)
	sort.Strings(smis)

	if len(movies) != len(smis) {
		return
	}

	for i, v := range movies {
		smi := smis[i]
		smiSp := strings.Split(smi, ".")
		smiExt := smiSp[len(smiSp)-1]

		movieSp := strings.Split(v, ".")
		movieExt := movieSp[len(movieSp)-1]

		fmt.Println(smiExt, movieExt)

		newName := strings.Replace(v, "."+movieExt, "."+smiExt, 1)
		fmt.Println(newName)

		os.Rename(targetDir+"/"+smi, targetDir+"/"+newName)
	}

	fmt.Println(movies, len(movies))
	fmt.Println(smis, len(smis))
}
