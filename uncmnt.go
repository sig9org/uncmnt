package main

import(
  "bufio"
  "flag"
  "fmt"
  "os"
  "regexp"
)

func main() {
  r1 := regexp.MustCompile(`^\s*#`)
  r2 := regexp.MustCompile(`^\s*$`)
  r3 := regexp.MustCompile(`^\s*!`)
  flag.Parse()
  args := flag.Args()
  for _, v := range args{
    fp, err := os.Open(v)
    if err != nil{
      //fmt.Println("error")
    }
    defer fp.Close()

    scanner := bufio.NewScanner(fp)
    for scanner.Scan() {
      txt := scanner.Text()
      ms1 := r1.MatchString(txt)
      ms2 := r2.MatchString(txt)
      ms3 := r3.MatchString(txt)
      if !ms1 && !ms2 && !ms3 {
        fmt.Println(scanner.Text())
      }
    }
  }
}
