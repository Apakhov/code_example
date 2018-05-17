package main

import (
	"io"
	"os"
	"sort"
	"strconv"
)

type byName []os.FileInfo

func (bn byName) Len() int {
	return len(bn)
}
func (bn byName) Less(i, j int) bool {
	return bn[i].Name() < bn[j].Name()
}
func (bn byName) Swap(i, j int) {
	bn[i], bn[j] = bn[j], bn[i]
}

func getSize(fi os.FileInfo) string {
	if fi.Size() == 0 {
		return "(empty)"
	}
	return "(" + strconv.FormatInt(fi.Size(), 10) + "b)"

}

func findLastDirName(fileInfos []os.FileInfo) (last string) {
	for _, fi := range fileInfos {
		if fi.IsDir() && fi.Name()[0] != '.' {
			last = fi.Name()
		}
	}
	return
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	return dirTreeRec(out, "", path, !printFiles)
}

func dirTreeRec(out io.Writer, prefix string, dirName string, noFiles bool) error {
	dir, err := os.Open(dirName)
	if err != nil {
		return err
	}
	defer dir.Close()
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return err
	}
	sort.Sort(byName(fileInfos))

	lastDirName := findLastDirName(fileInfos)
	var lastIndex int
	for i, fi := range fileInfos {
		if fi.Name()[0] != '.' {
			lastIndex = i
		}
	}
	for i, fi := range fileInfos {
		if fi.Name()[0] == '.' {
			continue
		}
		if fi.IsDir() {
			if fi.Name() == lastDirName {
				if noFiles || i == lastIndex {
					out.Write([]byte(prefix + "└───" + fi.Name() + "\n"))
					dirTreeRec(out, prefix+"\t", dirName+"/"+fi.Name(), noFiles)
				} else {
					out.Write([]byte(prefix + "├───" + fi.Name() + "\n"))
					dirTreeRec(out, prefix+"│\t", dirName+"/"+fi.Name(), noFiles)
				}

			} else {
				out.Write([]byte(prefix + "├───" + fi.Name() + "\n"))
				dirTreeRec(out, prefix+"│\t", dirName+"/"+fi.Name(), noFiles)
			}
		} else if !noFiles {
			if i == lastIndex {
				out.Write([]byte(prefix + "└───" + fi.Name() + " " + getSize(fi) + "\n"))
			} else {
				out.Write([]byte(prefix + "├───" + fi.Name() + " " + getSize(fi) + "\n"))
			}
		}
	}
	return nil
}

func main() {
	out := (os.Stdout)
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}

	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
