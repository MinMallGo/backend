package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
)

func ReadFile(filename string) ([]byte, error) {
	files, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return reverse(files), nil
}

func WriteFile(filename string, data []byte, perm os.FileMode) error {
	data = reverse(data)
	err := os.WriteFile(filename, data, perm)
	if err != nil {
		return err
	}
	return nil
}

func reverse(s []byte) []byte {
	str := []rune(string(s))
	left, right := 0, len(str)-1
	for left < right {
		str[left], str[right] = str[right], str[left]
		left++
		right--
	}
	return []byte(string(str))
}

// 分段读取，然后转换，再文件头追加内容
// 1. 分段读取
// 2. reverse
// 3. 往文件头写入

func ReadFile2(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine() // 按行读取文件
		if err == io.EOF {                // 用于判断文件是否读取到结尾
			break
		}
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", line)
		// todo 往文件头部追加
	}
	return nil
}

// 分多个文件存储  1T 100
// 一个文件保存所有的文件名
// 分片存放。size <= 2G

// 读取写入之后判断文件的大小，超过之后新建文件再进行写入

// 如果是新建则需要返回文件名。
const gb = 2 ^ 30
const singleFileSize = gb * 10
const fileChunk = 2 ^ 21

func WriteFile2(filename string, data []byte) (string, error) {

	// 数据反转
	data = reverse(data)

	stat, err := os.Stat(filename)
	if err != nil {
		return filename, err
	}
	// todo 这里需要判断一下文件大小然后再决定是创建文件还是继续写入
	if stat.Size() < singleFileSize {
		//
		file, err := os.OpenFile(filename, os.O_RDWR, 0644)
		if err != nil {
			return "", err
		}
		_, err = file.WriteAt(data, 0)
		if err != nil {
			return "", err
		}
	} else {
		// todo 往文件头部追加
		filename = filename + "_1"
		file, err := os.OpenFile(filename, os.O_RDWR, 0644)
		if err != nil {
			return "", err
		}
		_, err = file.WriteAt(data, 0)
		if err != nil {
			return "", err
		}
	}

	if err != nil {
		return filename, err
	}

	return filename, nil
}

func BigFileSplite(filename string) error {

	fileToBeChunked := filename
	file, err := os.Open(fileToBeChunked)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	var fileSize int64 = fileInfo.Size()

	//const fileChunk = 1 * (1 << 20) // 1 MB, change this to your requirement
	// calculate total number of parts the file will be chunked into
	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))
	fmt.Printf("Splitting to %d pieces.\n", totalPartsNum)

	totalFileNum := uint64(fileSize / singleFileSize) // 总文件个数
	fileNames := make([]string, 0, totalFileNum)

	for i := uint64(1); i < totalPartsNum; i++ {
		partSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
		partBuffer := make([]byte, partSize)

		_, err = file.Read(partBuffer)
		if err != nil {
			return err
		}

		// write to disk
		fileName := "somebigfile_" + strconv.FormatUint(i, 10)
		_, err := os.Create(fileName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// write/save buffer to disk
		// todo reverse

		targetFileName, err := WriteFile2(fileName, partBuffer) // 需要保存文件名。
		if err != nil {
			return err
		}
		//ioutil.WriteFile(fileName, partBuffer, os.ModeAppend)
		fmt.Println("Split to : ", targetFileName)
		// 最后才是写入这个文件
		fileNames = append(fileNames, targetFileName)
	}
	return nil
}

// 直接追加
func Write3(filename string, data []byte) {
	// 从字符串/文件尾部读取
	// 直接append到文件
	// 1234 43
}

// 读文件则是从末尾开始读取，记录文件的偏移量
func Read3(fileName string) error {

	return nil
}

// todo 这里有个文件来保存写入的那些文件？ 或者文件夹就行了吧
