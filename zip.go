package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func main() {
	const dir = "./1/"

	CompressZip(dir)              //压缩
	DeCompressZip("1.zip", "1_2") //解压缩
}

func innerCompress(zipWriter *zip.Writer, dir string) {
	fmt.Println("now innerCompress ", dir)

	//获取源文件列表
	f, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range f {
		var filepath string

		if strings.HasSuffix(dir, "/") {
			filepath = fmt.Sprintf("%s%s", dir, file.Name())
		} else {
			filepath = fmt.Sprintf("%s%s%s", dir, "/", file.Name())
		}

		if file.IsDir() {
			//fmt.Println(fmt.Sprintf("%s is a dir", file.Name()))
			if strings.HasSuffix(dir, "/") {
				innerCompress(zipWriter, filepath)
			} else {
				innerCompress(zipWriter, filepath)
			}
			continue
		}

		fw, _ := zipWriter.Create(filepath)
		filecontent, err := ioutil.ReadFile(filepath)
		if err != nil {
			fmt.Println(err)
		}
		n, err := fw.Write(filecontent)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(filepath, n)
	}
}

func CompressZip(root string) {
	fzip, _ := os.Create("1.zip")
	w := zip.NewWriter(fzip)
	defer w.Close()

	innerCompress(w, root)
}

func DeCompressZip(srcpath, dstpath string) {
	File := srcpath
	dir := dstpath

	os.Mkdir(dir, 0777) //创建一个目录

	cf, err := zip.OpenReader(File) //读取zip文件
	if err != nil {
		fmt.Println("OpenReader", err)
	}
	defer cf.Close()

	for _, file := range cf.File {
		rc, err := file.Open()
		if err != nil {
			fmt.Println("Open", err)
		}

		var filename string
		if strings.HasPrefix(file.Name, "./") {
			filename = file.Name[1:]
		} else {
			filename = file.Name
		}

		os.MkdirAll(path.Dir(dir+filename), 0755)
		f, err := os.Create(dir + filename)
		if err != nil {
			fmt.Println("Create", err)
		}
		defer f.Close()

		n, err := io.Copy(f, rc)
		if err != nil {
			fmt.Println("Copy", err)
		}

		fmt.Println(n)
	}

}
