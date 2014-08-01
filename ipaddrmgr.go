package main

import (
	"fmt"
	"net"
	"log"
	"flag"
	"net/smtp"
	"bytes"
	"text/template"
	"strings"
)

func main() {
	
	var hostName string
	flag.StringVar(&hostName, "hostName", "localhost", "host name to report")
	var emailUsername string	
	flag.StringVar(&emailUsername, "emailUsername", "", "email username")
	var emailPassword string
	flag.StringVar(&emailPassword, "emailPassword", "", "email password")
	var emailServer string
	flag.StringVar(&emailServer, "emailServer", "smtp.gmail.com", "email server")
	var emailServerPort int
	flag.IntVar(&emailServerPort, "emailPort", 587, "email server port")
	var emailToAddress string
	flag.StringVar(&emailToAddress, "emailToAddress", "", "Email to address")	
	var subnetMatch string
	flag.StringVar(&subnetMatch, "subnetMatch", "192.168", "subnet match")
	flag.Parse()
			
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	var ipAddresses string
	for _, i := range ifaces {
	    addrs, err := i.Addrs()
		if err != nil {
			log.Fatal(err)
		}
	    for _, addr := range addrs {
			if strings.Contains(addr.String(), subnetMatch) {
				fmt.Println("ipaddress", addr)
				ipAddresses += "ipaddress: " + addr.String() + "\n"
			}
	    }
	}
		
	sendEmail(emailUsername, emailPassword, emailToAddress, emailServer, emailServerPort, hostName, ipAddresses)
	
}

func sendEmail(emailUsername string, emailPassword string, emailToAddress string, emailServer string, 
	emailServerPort int, hostName string, ipAddresses string) () {
		
	parameters := struct {
	        From string
	        To string
	        Subject string
	        Message string
	} {
	        emailUsername,
	        emailToAddress,
	        "ip adress update - " + hostName,
	        ipAddresses,
	}

	buffer := new(bytes.Buffer)

	template := template.Must(template.New("emailTemplate").Parse(emailScript()))
	template.Execute(buffer, &parameters)

	auth := smtp.PlainAuth("", emailUsername, emailPassword, emailServer)

	smtp.SendMail(
	      fmt.Sprintf("%s:%d", emailServer, emailServerPort),
	      auth,
	      emailUsername,
	      []string{emailToAddress},
	      buffer.Bytes())
		  	
}

func emailScript() (script string) {
    return `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}
MIME-version: 1.0
Content-Type: text/html; charset="UTF-8"

{{.Message}}`
}

