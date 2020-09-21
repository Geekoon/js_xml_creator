package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func main() {
	const myurl = "http://localhost:8000"
	const xmlbody = `
<request>
    <parameters>
        <email>test@test.com</email>
        <password>test</password>
    </parameters>
</request>`

	resp, err := http.Post(myurl, "text/xml", strings.NewReader(xmlbody))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.Header)
	// Do something with resp.Body
}
