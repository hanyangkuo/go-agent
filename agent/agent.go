package main

import (
	"com.young/agent/util"
	"crypto/aes"
	"crypto/cipher"
	"flag"
	"fmt"
	_ "github.com/kardianos/service"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const version = "1.0.0"

var (
	rootDir = ""
)

func init() {
	flagVerbose := flag.Bool("verbose", false, "")
	flagVer := flag.Bool("version", false, "print")
	flag.Parse()
	if *flagVer {
		fmt.Print(version)
		os.Exit(0)
	}

	exepath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	rootDir = filepath.Dir(exepath)
	os.Mkdir(filepath.Join(rootDir, "log"), os.ModePerm)
	file, err := rotatelogs.New(
		filepath.Join(rootDir, "log", "agent.log.%Y%m%d"),
		rotatelogs.WithMaxAge(time.Duration(72)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.SetFormatter(&log.JSONFormatter{})
	if *flagVerbose {
		log.SetOutput(io.MultiWriter(file, os.Stdout))
	} else {
		log.SetOutput(file)
	}
}


func main() {
	go func(){
		for {
			calculate()
		}
	}()
	for {

		//userTime, kernelTime, residentSize, totalSize := getProcessRow()
		////fmt.Println(userTime, kernelTime, residentSize, totalSize)
		//fmt.Println(residentSize, totalSize)
		//fmt.Printf("residentSize: %d\ttotalSize: %d\n", residentSize, totalSize)
		//p := &process.Process{Pid: int32(os.Getpid())}
		//m, err := p.MemoryInfo()
		//if err != nil {
		//	log.Println(err)
		//}
		//fmt.Printf("RSS: %d\tVMS: %d\n", m.RSS, m.VMS)
		time.Sleep(time.Second)
	}
}
var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
func calculate() {
	//需要去加密的字串
	plaintext := []byte("My name is Astaxie")
	//如果傳入加密串的話，plaint 就是傳入的字串
	if len(os.Args) > 1 {
		plaintext = []byte(os.Args[1])
	}

	//aes 的加密字串
	key_text := "astaxie12798akljzmknm.ahkjkljl;k"
	if len(os.Args) > 2 {
		key_text = os.Args[2]
	}

	//fmt.Println(len(key_text))

	// 建立加密演算法 aes

	c, err := aes.NewCipher([]byte(key_text))
	if err != nil {
		fmt.Printf("Error: NewCipher(%d bytes) = %s", len(key_text), err)
		os.Exit(-1)
	}

	//加密字串
	cfb := cipher.NewCFBEncrypter(c, commonIV)
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)
	//fmt.Printf("%s=>%x\n", plaintext, ciphertext)

	// 解密字串
	cfbdec := cipher.NewCFBDecrypter(c, commonIV)
	plaintextCopy := make([]byte, len(plaintext))
	cfbdec.XORKeyStream(plaintextCopy, ciphertext)
	//fmt.Printf("%x=>%s\n", ciphertext, plaintextCopy)
}

func testFunc(){
	defer util.TimeTrack(time.Now(), "testFunc")
	time.Sleep(1*time.Second)
	runtime.GC()
}
