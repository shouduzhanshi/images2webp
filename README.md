**随着需求的增加,App 体积逐渐增大,精简App的最简单便捷的办法就是从资源文件(icon)入手.**

**谷歌推出的WEBP是个很好的选择,所以,我们需要把项目中用到的资源文件(icon)转换成webp格式**

> WebP是一种现代图像格式，可为Web上的图像提供出色的无损和有损压缩。使用WebP，网站管理员和Web开发人员可以创建更小，更丰富的图像，从而使Web更快。
与PNG相比，WebP无损图像的尺寸要[小26％](https://developers.google.cn/speed/webp/docs/webp_lossless_alpha_study#results)。在同等的[SSIM](https://en.wikipedia.org/wiki/Structural_similarity)质量指标下，WebP有损图像比同类JPEG图像 [小25-34％](https://developers.google.cn/speed/webp/docs/webp_study)。[](https://en.wikipedia.org/wiki/Structural_similarity)
无损WebP 支持透明性（也称为Alpha通道），而其[额外字节](https://developers.google.cn/speed/webp/docs/webp_lossless_alpha_study#results)仅为[22％](https://developers.google.cn/speed/webp/docs/webp_lossless_alpha_study#results)。对于可接受的有损RGB压缩，有损WebP还支持透明性，与PNG相比，文件大小通常小3倍。

#批量转换图片为webp格式流程:

下载webp库:
[Windows](https://storage.googleapis.com/downloads.webmproject.org/releases/webp/libwebp-1.0.3-windows-x86.zip)
[Mac](https://storage.googleapis.com/downloads.webmproject.org/releases/webp/libwebp-1.0.3-rc1-mac-10.14.tar.gz)
[Linux](https://storage.googleapis.com/downloads.webmproject.org/releases/webp/libwebp-1.0.3-linux-x86-64.tar.gz)

以Mac 为例:
下载得到的webp库:
![WX20191014-215407@2x.png](https://upload-images.jianshu.io/upload_images/5428263-ca64c34fd2e4eabb.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

下载得到的批量处理程序:

![WX20191014-224353@2x.png](https://upload-images.jianshu.io/upload_images/5428263-1bdf8596ed377f8b.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

将批处理程序 png2webp 拷贝到 解压后的 webp库bin目录下

![image.png](https://upload-images.jianshu.io/upload_images/5428263-7e39a7e10933a736.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)


打开终端(控制台)进入webp库bin目录:

![](https://upload-images.jianshu.io/upload_images/5428263-a1adc5139901eb56.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

`执行 ./img2webp [android项目根路径] [是否保留原文件true:false] [webp转换参数]`
``` //fmt.Println("第一个参数项目根路径")
	//fmt.Println("第二个参数是否保留原文件 true 保留 false 转换完成后删除")
	//fmt.Println("第三个参数开始包含第三个参数为webp转换自定义参数 详见:https://developers.google.cn/speed/webp/docs/using")
	//fmt.Println("作者使用参数:", "-q", "80", "-mt", "-v", "-progress", "-o")
	//fmt.Println("示例:", "/Users/lion/Android/project/android false -q 80 -mt -v -progress")
	//fmt.Println("Enter 开始转换"
```

#成果:

**转换webp前:**
![](https://upload-images.jianshu.io/upload_images/5428263-5798afd8c15b2e40.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

**转换webp后:**
![image.png](https://upload-images.jianshu.io/upload_images/5428263-d4eb56f8e8a6b253.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

**文件校验:**
**经过版本控制工具查验,不多不少 没有产生文件丢失的情况**
![结果校验](https://upload-images.jianshu.io/upload_images/5428263-18978afdb551bdaf.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

#批处理程序源码

```package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

var utilsPath string

func main() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	cmdParms := os.Args
	utilsPath = dir + "/cwebp"
	projectPath := cmdParms[1]
	isReplace, _ := strconv.ParseBool(cmdParms[2])
	parms := make([]string, 0)
	for i := 3; i < len(cmdParms); i++ {
		parms = append(parms, cmdParms[i])
	}
	parms = append(parms, "-o")
	readDir(projectPath, isReplace, parms)
}

func readDir(path string, isReplace bool, parms []string) {
	if infos, e := ioutil.ReadDir(path); e == nil {
		for _, temp := range infos {
			dir := temp.IsDir()
			name := temp.Name()
			if dir && !strings.Contains(name, "assets") {
				readDir(path+"/"+name, false, parms)
			} else if strings.HasSuffix(name, ".JPG") || strings.HasSuffix(name, ".JPEG") ||
				strings.HasSuffix(name, ".jpg") || strings.HasSuffix(name, ".jpeg") ||
				strings.HasSuffix(name, ".png") || strings.HasSuffix(name, ".PNG") {
				split := strings.Split(name, ".")
				dispFile(path+"/"+name, path, split[0], parms)
				if !isReplace {
					os.Remove(path + "/" + name)
				}
			}
		}
	}
}

func dispFile(path, dirPath, name string, parms []string) {
	parms = append(parms, dirPath+"/"+name+".webp")
	parms = append(parms, path)
	execCmd(utilsPath, parms)
}

func execCmd(shell string, raw []string) (int, error) {
	fmt.Println(shell, raw)
	cmd := exec.Command(shell, raw...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return 0, nil
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println(err)
		return 0, nil
	}
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
		return 0, nil
	}
	s := bufio.NewScanner(io.MultiReader(stdout, stderr))

	for s.Scan() {
		text := s.Text()
		fmt.Println(text)
	}
	if err := cmd.Wait(); err != nil {
		fmt.Println(err)
	}
	return 0, nil
}
```
##[github地址](https://github.com/hhhkk/images2webp.git)


#联系我:
`853151446@qq.com`
`QQ:853151446`

#如果帮助到了你,麻烦你帮我点个赞,你的支持就是我的动力,万分感谢.
