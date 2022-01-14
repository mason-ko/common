package recursive_move

import (
	"fmt"
	"io/ioutil"
	"os"
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

func MoveSubtitle(targetDir string) {
	dir, err := ioutil.ReadDir(targetDir)
	if err != nil {
		panic(err)
	}

	for _, v := range dir {
		fmt.Println(v)
	}
}
