package main

import (
	"fmt"
	"math/big"
	"net"
)

type addr struct {
	ip4 string
	ip6 string
}

func main() {
	i := ip6ToInt("::d:1:1:0:101")
	fmt.Println(i)

	i.Add(i, big.NewInt(11))
	fmt.Println(i)
	fmt.Println(intToNet6(i))
}

func ip6ToInt(ip string) *big.Int {
	ret := big.NewInt(0)
	fmt.Println([]byte(net.ParseIP(ip).To16()))
	ret.SetBytes([]byte(net.ParseIP(ip).To16()))
	fmt.Println(ret)
	return ret
}

func intToNet6(ip *big.Int) string {
	b := make([]byte, 16-len(ip.Bytes()))
	b = append(b, ip.Bytes()...)
	return net.IP(b).String() + "/32"
}

func net6ToInt(inet string) *big.Int {
	ip, _, _ := net.ParseCIDR(inet)
	ret := big.NewInt(0)
	ret.SetBytes(ip.To16())
	return ret
}
