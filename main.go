package main

import (
	"fmt"
	"io"
	"os"
    "bytes"
    "sort"
    "strconv"
)

type FileInfoByName []os.FileInfo

func (fin FileInfoByName) Len() int      { return len(fin) }
func (fin FileInfoByName) Swap(i, j int) { fin[i], fin[j] = fin[j], fin[i] }
func (fin FileInfoByName) Less(i, j int) bool {
    return fin[i].Name() < fin[j].Name();
}

const (
	Br_horiz = "─"
	Br_vert  = "│"
	Br_split = "├"
	Br_last  = "└"
)

func dismissFiles(dirContents []os.FileInfo) (res []os.FileInfo) {
    for _, item := range dirContents {
        if item.IsDir() == true {
            res = append(res, item)
        }
    }
    return
}

func drawItem(buffer *bytes.Buffer, fInfo os.FileInfo, level int, emptyLevels int, last bool) {
    notEmpty := level - emptyLevels
    for i := 0; i < notEmpty - 1; i++ {
        buffer.WriteString(Br_vert)
        buffer.WriteString("\t")
    }
    for i := 0; i < emptyLevels; i++ {
        buffer.WriteString("\t")
    }
    if last == true {
        buffer.WriteString(Br_last  + Br_horiz + Br_horiz + Br_horiz + fInfo.Name())
    } else {
        buffer.WriteString(Br_split + Br_horiz + Br_horiz + Br_horiz + fInfo.Name())
    }
    if fInfo.IsDir() == false {
        itemSize := fInfo.Size()
        strItemSize := "empty"
        if itemSize > 0 {
            strItemSize = strconv.FormatInt(itemSize, 10) + "b"
        }
        buffer.WriteString(" (" + strItemSize + ")")
    }
    buffer.WriteString("\n")
}

func drawDirContents(buffer *bytes.Buffer, path string, printFiles bool, level int, emptyLevels int, last bool) (err error) {

    currItem, err := os.Open(path)
    if err != nil {
	    return err
	}

    dirContents, err := currItem.Readdir(-1)
    if err != nil {
	    return err
	}
    if printFiles == false {
        dirContents = dismissFiles(dirContents)
    }

    dirItemsCount := len(dirContents)

    if dirItemsCount > 0 {
        sort.Sort(FileInfoByName(dirContents))
        if last == true {
            emptyLevels++;
        }
        for i := 0; i < dirItemsCount; i++ {
            itemPath := path + string(os.PathSeparator) + dirContents[i].Name()
            itemInfo, err := os.Stat(itemPath)
            if err != nil {
		        return err
	        }
            isLast := i == dirItemsCount - 1
            drawItem(buffer, itemInfo, level + 1, emptyLevels, isLast)

            if itemInfo.IsDir() == true {
                rerr := drawDirContents(buffer, itemPath, printFiles, level + 1, emptyLevels, isLast)
                if rerr != nil {
                    return rerr;
                }
            }
        }
    }
    return nil;
}

func dirTree(out io.Writer, path string, printFiles bool) (err error) {
    var buffer bytes.Buffer
	err = drawDirContents(&buffer, path, printFiles, 0, 0, false)
    bufferContents := buffer.String()
    fmt.Fprintln(out, bufferContents[:len(bufferContents)-1])
    return err
}

func main() {
	out := os.Stdout
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
