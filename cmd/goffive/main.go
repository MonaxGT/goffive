package main

import (
	"flag"
	"fmt"

	"github.com/MonaxGT/goffive"
)

func main() {
	userPtr := flag.String("u", "", "User")
	passPtr := flag.String("p", "", "Password")
	urlPtr := flag.String("url", "", "URL F5 Big-IP")
	flag.Parse()
	var err error
	var conf *goffive.Client

	conf, err = goffive.New(*userPtr, *passPtr, *urlPtr)
	if err != nil {
		panic(err)
	}

	/*err = conf.LTM()
	if err != nil {
		panic(err)
	}*/

	pools, err := conf.LTM.Pools()
	if err != nil {
		panic(err)
	}
	for _, v := range pools {
		fmt.Println(v.FullPath)
		fmt.Println(v.Name)
	}

	policies, err := conf.ASM.Policies()
	if err != nil {
		panic(err)
	}
	for _, v := range policies {
		fmt.Println(v.Name)
		fmt.Println(v.VirtualServers)
	}


	signatories, err :=conf.ASM.Signatories("08E2vm9ejesvAX8paBY1bg")
	if err != nil {
		panic(err)
	}

	for _, v := range signatories {
		fmt.Println(v.SignatureReference.Name)
		fmt.Println(v.Alarm)
	}
}
