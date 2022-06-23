package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"

	"github.com/robfig/cron/v3"
	"golang.org/x/text/encoding/simplifiedchinese"
)
type Charset string

const (
    UTF8    = Charset("UTF-8")
    GB18030 = Charset("GB18030")
    HELP = `
# 分   时   日   月   周 命令行
# -    -    -    -    -
# |    |    |    |    |
# |    |    |    |    +----- 星期中星期几 (0-6) (星期天为0)
# |    |    |    +---------- 月份 (1-12)
# |    |    +--------------- 一个月中的第几天 (1-31)
# |    +-------------------- 小时 (0-23)
# +------------------------- 分钟 (0-59)
#
# * 时表示每单位时间都要执行
# a-b 时表示从第 a 单位时间到第 b 单位时间这段时间内要执行
# */n 时表示每 n 单位时间个时间间隔执行一次
# a, b, c,... 时表示第 a, b, c,... 单位时间要执行
#
# 示例:每5分钟执行bat批处理
# */5 * * * * d:/subversion/update.bat
`
)

var version = "1.0.0"

type RunJob struct {
    cmd string
}

func (r *RunJob) Run() {
    Command(r.cmd)
}

func main() {
    args := os.Args
    if len(args) > 1 && (args[1] == "-v" || args[1] == "--version"){
        fmt.Println("Crontab Version:",version)
        return
    }
    if len(args) > 1 && (args[1] == "-h" || args[1] == "--help"){
        fmt.Println("Crontab Help")
        fmt.Println("Setting config.cfg like:")
        fmt.Println(HELP)
        return
    }
    tasks, err := Config("./config.cfg")
    if err != nil {
        fmt.Println("config error:", err)
        return
    }

    c := cron.New()
    for _, t := range tasks {
        fmt.Println(t)
        params := strings.Split(t," ")
        Task(c,params)
    }
    c.Start()
    select {}
}

func Task(c *cron.Cron, params []string) error {
    time := params[0:6]
    timeStr := strings.Join(time, ", ")
    // fmt.Println(timeStr)

    cmd := strings.Join(params[6:]," ")
    // fmt.Println(cmd)

    j := &RunJob{
        cmd:cmd,
    }
    c.AddJob(timeStr,j)
    return nil
}

func MakeConfigTemplate(filePth string) {
    cfg := strings.Replace(HELP,"\n","",1) //去除句首空行
    err := ioutil.WriteFile(filePth, []byte(cfg), 0666) //写入文件(字节数组)
    if err != nil {
        fmt.Println(err.Error())
    }
}

func Config(filePth string) ([]string, error) {
    f, err := os.Open(filePth)
    if err != nil {
        if _, err := os.Stat(filePth); os.IsNotExist(err) {
            MakeConfigTemplate(filePth)
        }
        fmt.Println("Please setting config <config.cfg>")
        return nil,err
    }
    defer f.Close()
    conf, _ := ioutil.ReadAll(f)
    conf_string := string(conf)
    tasks := strings.Split(conf_string,"\n")
    var result []string
    for _, task := range tasks {
        if task == "" {
            continue
        }
        if strings.HasPrefix(task, "#") {
            continue
        }
        result = append(result,task)
    }
    return result,nil
}

func Command(cmd string) error {
    sysType := runtime.GOOS
    Msg().Info(cmd)
    var params []string
    if sysType == "linux" {
       // LINUX系统
       params = append(params,"bash", "-c", cmd)
    }

    if sysType == "windows" {
        // windows系统
        params = append(params,"cmd", "/C", cmd)
    }
    c := exec.Command(params[0],params[1:]...)

    c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

    stdout, err := c.StdoutPipe()
    if err != nil {
        fmt.Println(err)
        return err
    }
    c.Start()
    in := bufio.NewScanner(stdout)
    for in.Scan() {
        cmdRe:=ConvertByte2String(in.Bytes(),"GB18030")
        fmt.Println(cmdRe)
    }
    c.Wait()
    fmt.Println()
    return nil
}

func WindowsShell(params []string,hideWindow bool) error {

    var temArg []string

    // 隐藏powershell窗口
    temArg = append(temArg, "-WindowStyle")
    temArg = append(temArg, "Hidden")

    // 启动目标程序
    temArg = append(temArg, "-Command")
    temArg = append(temArg, "Start-Process")
    temArg = append(temArg, params...)

    // 静默启动参数
    if hideWindow {
        temArg = append(temArg, "-WindowStyle")
        temArg = append(temArg, "Hidden")
    }
    cmd := exec.Command("PowerShell.exe", temArg...)

    // 启动时隐藏powershell窗口,没有这句会闪一下powershell窗口
    cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

    err := cmd.Run()
    if nil != err {
        fmt.Println(err)
        return err
    }
    fmt.Println("start exe:", params)
    return nil
}

func ConvertByte2String(byte []byte, charset Charset) string {
    var str string
    switch charset {
    case GB18030:
        var decodeBytes,_=simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
        str= string(decodeBytes)
    case UTF8:
        fallthrough
    default:
        str = string(byte)
    }
    return str
}
