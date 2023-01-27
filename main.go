package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	User_interface := os.Args[1]
	Turnonmon(User_interface)
	Ap_mac := os.Args[2]
	if Ap_mac == "" {
		Usage()
	}
	Station_mac := ""
	Auth_flag := ""
	if os.Args[3] != "" {
		Station_mac = os.Args[3]
	}
	if os.Args[4] != "" {
		if os.Args[4] == "-c" {
			Auth_flag = os.Args[4]
		} else {
			Usage()
		}
	}
	Deauth_Attack(User_interface, Ap_mac, Station_mac, Auth_flag)
}

func Usage() {
	fmt.Println(" - AP broadcast frame")
	fmt.Println(" 	- ./Deauth-Attack <interface> <ap mac>")
	fmt.Println(" - AP unicast, Station unicast frame")
	fmt.Println(" 	- ./Deauth-Attack <interface> <ap mac> <station mac>")
	fmt.Println(" - authentication frame")
	fmt.Println(" 	- ./Deauth-Attack <interface> <ap mac> <station mac> -c")
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

func Deauth_Attack(User_interface string, Ap_mac string, Station_mac string, Auth_flag string) {
	fmt.Println("_start_")
	println(User_interface, Ap_mac, Station_mac, Auth_flag)
}
