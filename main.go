package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
)

var GLOBAL_PATH_CONFIGURATION_FALLBACK = ""

// Structurize the configuration
type Configuration struct {
    Path []string `json:"path"`
    ExcludeExt []string `json:"excludeExt"`
}

var PATH_TO_TARGET_DELETION []string = []string{}
var TOTAL_BYTE_FOUND int64 = 0


// source from : https://gist.github.com/anikitenko/b41206a49727b83a530142c76b1cb82d?permalink_comment_id=4467913#gistcomment-4467913
// thanks to https://gist.github.com/maxmcd
func byteSize(b int) string {
    bf := float64(b)
    for _, unit := range []string{"", "KB", "MB", "GB", "TB"} {
        if math.Abs(bf) < 1024.0 {
            return fmt.Sprintf("%3.1f%sB", bf, unit)
        }
        bf /= 1024.0
    }
    return fmt.Sprintf("%.1fYiB", bf)
}

func main(){
    sp := spinner.New(spinner.CharSets[11], 100* time.Millisecond)
    promptDeletion := promptui.Prompt{
        Label: "Delete the file ?",
        IsConfirm: true,
    }
    config := Configuration{
        Path: []string{ "/home/bluespada/.cache", "/tmp", "/home/bluespada/.local/share/trash", "/home/bluespada" },
    }
    
    sp.Start()
    for _, v := range config.Path {
        sp.Color("green", "bold")
        sp.Suffix = fmt.Sprintf(" Scanning %s", v)
        _ = filepath.Walk(v, func(path string, info os.FileInfo, err error) error {
            sp.Color("green bold")
            sp.Suffix = fmt.Sprintf(" Calculating %s", v);
            file, err := os.Open(path)
            if err != nil {
                return err
            }
            defer file.Close()
            fstat, err := file.Stat()
            if err != nil {
                return err;
            }
            TOTAL_BYTE_FOUND += fstat.Size()
            return nil
        })
    }
    sp.Stop()
    fmt.Println("total size", byteSize(int(TOTAL_BYTE_FOUND)))
    result, err := promptDeletion.Run()
    if err != nil {
        return
    }
    if strings.ToLower(result) == "y" {
        fmt.Println("Delete")
    }
}
