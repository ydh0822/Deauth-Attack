package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	User_interface := flag.String("i", "", "a string")
	Ap_mac := flag.String("a", "", "a string")
	Station_mac := flag.String("s", "", "a string")
	Auth_flag := flag.Bool("c", false, "a bool")
	flag.Parse()
	if *User_interface == "" {
		Usage()
	} else if *Ap_mac == "" {
		Usage()
	}
	Deauth_Attack(*User_interface, *Ap_mac, *Station_mac, *Auth_flag)
}

type StringArray []string

func (arr *StringArray) String() string {
	return fmt.Sprintf("%v", *arr)
}

func (arr *StringArray) Set(s string) error {
	*arr = strings.Split(s, ",")
	return nil
}

func Usage() {
	fmt.Println(" - AP broadcast frame")
	fmt.Println(" 	- ./Deauth-Attack -i <interface> -a <ap mac>")
	fmt.Println(" - AP unicast, Station unicast frame")
	fmt.Println(" 	- ./Deauth-Attack -i <interface> -a <ap mac> -s <station mac>")
	fmt.Println(" - authentication frame")
	fmt.Println(" 	- ./Deauth-Attack -i <interface> -a <ap mac> -s <station mac> -c")
	os.Exit(-1)
}

func ExcuteCMD(script string, arg ...string) {
	cmd := exec.Command(script, arg...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		fmt.Println((err))
		os.Exit(-1)
	} else {
		fmt.Println(string(output))
	}
}

func Turnonmon(name string) {
	ExcuteCMD("sudo", "ifconfig", name, "down")
	fmt.Println(name + " is down")
	ExcuteCMD("sudo", "iwconfig", name, "mode", "monitor")
	fmt.Println(name + " turn monitor mode")
	ExcuteCMD("sudo", "ifconfig", name, "up")
	fmt.Println(name + " is up \n")
}

func AP_broadcast(User_interface string, Ap_mac string) {
	fmt.Println("AP_broadcast")
}

func AP_unicast(User_interface string, Ap_mac string, Station_mac string) {
	fmt.Println("AP_unicast")
}

func AP_unicast_authentication(User_interface string, Ap_mac string, Station_mac string) {
	fmt.Println("AP_unicast_authentication")
}

func Deauth_Attack(User_interface string, Ap_mac string, Station_mac string, Auth_flag bool) {
	if Station_mac == "" {
		AP_broadcast(User_interface, Ap_mac)
	} else if !Auth_flag {
		AP_unicast(User_interface, Ap_mac, Station_mac)
	} else {
		AP_unicast_authentication(User_interface, Ap_mac, Station_mac)
	}
}
