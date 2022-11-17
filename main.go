package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"os"
	"strings"
	"sync"
)

var (
  File string
  WaitGroup sync.WaitGroup
  
  colorReset string = "\033[0m"
  colorRed string = "\033[31m"
  colorGreen string = "\033[32m"
  colorBlue string = "\033[34m"
  colorCyan string = "\033[36m"
  colorWhite string = "\033[37m"
)

func ClearConsole() {
  fmt.Print("\033[H\033[2J")
}

func main() {
  ClearConsole()
	fmt.Printf("%v[%v?%v] %vMass Gmail Checker %v|%v By github.com/vertionn \n\n", colorBlue, colorWhite, colorBlue, colorWhite, colorBlue, colorWhite)

  fmt.Printf("%v[%v?%v] %vFile Name: %v",colorBlue, colorWhite, colorBlue, colorWhite, colorBlue)
  fmt.Scanln(&File)

  if File == "" {
    fmt.Printf("\n%vError: %vFile Name Cant Be Empty%v\n", colorRed, colorWhite, colorReset)
    time.Sleep(6 * time.Second)
    main()
  }

  ClearConsole()

  Emails := strings.Split(readFile(File), "\n")

  CreateFile("available.txt")

  start := time.Now()

  for _, EMAIL := range Emails {
    if len(EMAIL) > 0 {
      WaitGroup.Add(1)
      go func(EMAIL string) {
        defer WaitGroup.Done()

        var jsonStr = []byte(`{"username":"`+EMAIL+`","version":"3","firstName":"Github.com/","lastName":"vertionn"}`)
        requests, err := http.NewRequest("POST", "https://android.clients.google.com/setup/checkavail",bytes.NewBuffer(jsonStr))
        if err != nil {
          log.Fatalln(err)
        }
        requests.Header.Set("Content-Type","text/plain; charset=UTF-8")
        requests.Header.Set("Host","android.clients.google.com")
        requests.Header.Set("Connection","Keep-Alive")
        requests.Header.Set("user-agent","GoogleLoginService/1.3(m0 JSS15J)")
        client := &http.Client{}
        response, err := client.Do(requests)
        if err != nil {
          log.Fatalln(err)
        }

        body, err := ioutil.ReadAll(response.Body)
        if err != nil {
          log.Fatalln(err)
        }

        if strings.Contains(string(body), "SUCCESS") {
          fmt.Printf("%v[%v+%v] %v%v%v\n", colorBlue, colorGreen, colorBlue, colorGreen, EMAIL, colorReset)
          writeFile("available.txt", EMAIL+"\n")
          
        } else if strings.Contains(string(body), "USERNAME_UNAVAILABLE") {
          fmt.Printf("%v[%v-%v] %v%v%v\n", colorBlue, colorRed, colorBlue, colorRed, EMAIL, colorReset)
          
        } else {
          fmt.Printf("%v[%v?%v] %v%v%v\n", colorBlue, colorRed, colorBlue, colorRed, string(body), colorReset)
          
        }
      }(EMAIL)
    }
  }
  WaitGroup.Wait()

  end := time.Now()
  diff := end.Sub(start)

  out := time.Time{}.Add(diff)
  fmt.Printf("\n%v[%v?%v] %vFinished In %v%v %vSeconds%v\n", colorBlue, colorWhite, colorBlue, colorWhite, colorCyan, out.Second(), colorWhite, colorReset)
  
}

func readFile(path string) string {
  file, err := ioutil.ReadFile(path)
    if err != nil {
      fmt.Println(err)
    }
    return string(file)
}

func writeFile(path string, text string) {
  file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
      fmt.Println(err)
    }
    defer file.Close()
    if _, err = file.WriteString(text); err != nil {
      fmt.Println(err)
    }
}

func CreateFile(name string) {
  file, err := os.Create(name)
    if err != nil {
        log.Fatal(err)
    }
  
    defer file.Close()

}
