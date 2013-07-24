package main

import (
    "flag"
    "fmt"
    "os"
    "syscall"
    "index/suffixarray"
    "regexp"
)

func main() {
    flag.Usage = func () {
	fmt.Fprintf(os.Stderr, "usage: read [inputfile]\n")
	flag.PrintDefaults()
	os.Exit(2)
    }
    flag.Parse()
    argv := flag.Args()
    if len(argv) < 1 {
	fmt.Println("No input file")
	os.Exit(1)
    }

    map_file, err := os.Open(argv[0])
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer map_file.Close()

    fileLen, err := map_file.Seek(0, 2)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    _, err = map_file.Seek(0, 0)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    mmap, err := syscall.Mmap(int(map_file.Fd()), 0, int(fileLen), syscall.PROT_READ, syscall.MAP_SHARED)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    idx := suffixarray.New(mmap)
    re := regexp.MustCompile("(?P<first>A)(?P<last>a)")
    allAa := idx.FindAllIndex(re, -1)
    fmt.Println("All substr of form ?Aa? in file are :")
    for _, A := range allAa {
	fmt.Println(string(mmap[A[0]-1:A[1]+1]), A)
    }

    err = syscall.Munmap(mmap)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
