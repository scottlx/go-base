// Go NETCONF Client - Example
//
// Copyright (c) 2013-2018, Juniper Networks, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"

	"github.com/Juniper/go-netconf/netconf"
	"golang.org/x/crypto/ssh"
)

func main() {
	sshConfig := &ssh.ClientConfig{
		User:            "myadmin",
		Auth:            []ssh.AuthMethod{ssh.Password("Huawei@admin2021")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	s, err := netconf.DialSSH("10.50.45.209", sshConfig)

	if err != nil {
		log.Fatal(err)
	}

	defer s.Close()

	//fmt.Println(s.ServerCapabilities)
	//fmt.Println(s.SessionID)

	//getxml := `<get/>`
	editxml := `
		<edit-config>
	    <target>
	      <running/>
	    </target>
	    <config>
	      <interfaces xmlns="urn:ietf:params:xml:ns:yang:ietf-interfaces">
	        <interface>
	          <name>10GE1/0/3</name>
	          <ipv6 xmlns="urn:ietf:params:xml:ns:yang:ietf-ip">
	            <address>
	              <ip>2001:db8:c18:1::3</ip>
	              <prefix-length>128</prefix-length>
	            </address>
	          </ipv6>
	        </interface>
	      </interfaces>
	    </config>
	  </edit-config>`
	reply, err := s.Exec(netconf.RawMethod(editxml))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Reply: %+v", reply)
}
