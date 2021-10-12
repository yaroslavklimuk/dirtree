package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

type FileInfoByName []os.FileInfo

func (fin FileInfoByName) Len() int      { return len(fin) }
func (fin FileInfoByName) Swap(i, j int) { fin[i], fin[j] = fin[j], fin[i] }
func (fin FileInfoByName) Less(i, j int) bool {
	return fin[i].Name() < fin[j].Name()
}

type DirNode struct {
	prev *DirNode
	info os.FileInfo
	last bool
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

func drawItem(buffer *bytes.Buffer, node *DirNode) {
	prevNode := node.prev
	var itemLine string
	for prevNode.prev != nil {
		if prevNode.last == true {
			itemLine = "\t" + itemLine
		} else {
			itemLine = Br_vert + "\t" + itemLine
		}
		prevNode = prevNode.prev
	}
	buffer.WriteString(itemLine)
	if node.last == true {
		buffer.WriteString(Br_last + Br_horiz + Br_horiz + Br_horiz + node.info.Name())
	} else {
		buffer.WriteString(Br_split + Br_horiz + Br_horiz + Br_horiz + node.info.Name())
	}
	if node.info.IsDir() == false {
		itemSize := node.info.Size()
		strItemSize := "empty"
		if itemSize > 0 {
			strItemSize = strconv.FormatInt(itemSize, 10) + "b"
		}
		buffer.WriteString(" (" + strItemSize + ")")
	}
	buffer.WriteString("\n")
}

func drawDirContents(buffer *bytes.Buffer, node *DirNode, path string, printFiles bool) (err error) {
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
		for i := 0; i < dirItemsCount; i++ {
			itemPath := path + string(os.PathSeparator) + dirContents[i].Name()
			itemInfo, err := os.Stat(itemPath)
			if err != nil {
				return err
			}
			newNode := &DirNode{
				prev: node,
				info: itemInfo,
				last: i == dirItemsCount-1,
			}
			drawItem(buffer, newNode)

			if itemInfo.IsDir() == true {
				err = drawDirContents(buffer, newNode, itemPath, printFiles)
				if err != nil {
					return err
				}
			}
		}
	}

	err = currItem.Close()
	return err
}

func dirTree(out io.Writer, path string, printFiles bool) (err error) {
	var buffer bytes.Buffer
	itemInfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	newNode := &DirNode{
		info: itemInfo,
		last: false,
	}
	err = drawDirContents(&buffer, newNode, path, printFiles)
	bufferContents := buffer.String()
	_, err = fmt.Fprint(out, bufferContents)
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
