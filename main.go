package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	User_interface := flag.String("i", "", "interface name")
	Ap_mac := flag.String("a", "", "AP MAC Address == BSSID")
	Station_mac := flag.String("s", "", "Station MAC Address")
	Auth_flag := flag.Bool("c", false, "Authentication Attack")
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
	fmt.Println("[*] " + name + " is down")
	ExcuteCMD("sudo", "iwconfig", name, "mode", "monitor")
	fmt.Println("[*] " + name + " turn monitor mode")
	ExcuteCMD("sudo", "ifconfig", name, "up")
	fmt.Println("[*] " + name + " is up \n")
}

func _Init_(User_interface string) int {
	Turnonmon(User_interface)
	tmp_CH := ""
	fmt.Println("input Your AP Channel : ")
	fmt.Scanln(&tmp_CH)
	ExcuteCMD("sudo", "iwconfig", User_interface, "channel", tmp_CH)
	CH, err := strconv.Atoi(tmp_CH)
	if err != nil {
		fmt.Println(err)
	}
	return CH
}

func AP_broadcast(User_interface string, Ap_mac string) {
	CH := _Init_(User_interface)
	fmt.Println("AP_broadcast | Channel : ", CH)
	buffer := new(bytes.Buffer)
	Deauth_Header := [21]byte{0x00, 0x00, 0x0b, 0x00, 0x00, 0x80, 0x02, 0x00, 0x00, 0x00, 0x00, 0xc0, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	binary.Write(buffer, binary.LittleEndian, Deauth_Header)

	addr, err := net.ParseMAC(Ap_mac)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(addr)

	// Deauth_AP_MAC := [12]byte{0xb2, 0x9f, 0x30, 0x4e, 0xed, 0xdb, 0xb2, 0x9f, 0x30, 0x4e, 0xed, 0xdb}
	// Deauth_Footer := [4]byte{0x50, 0x4f, 0x07, 0x00}

}

func AP_unicast(User_interface string, Ap_mac string, Station_mac string) {
	CH := _Init_(User_interface)
	fmt.Println("AP_unicast | Channel : ", CH)
}

func AP_unicast_authentication(User_interface string, Ap_mac string, Station_mac string) {
	CH := _Init_(User_interface)
	fmt.Println("AP_unicast_authentication | Channel : ", CH)
}

func Deauth_Attack(User_interface string, Ap_mac string, Station_mac string, Auth_flag bool) {
	if Station_mac == "" {
		if Auth_flag {
			Usage()
		}
		AP_broadcast(User_interface, Ap_mac)
	} else if !Auth_flag {
		AP_unicast(User_interface, Ap_mac, Station_mac)
	} else {
		AP_unicast_authentication(User_interface, Ap_mac, Station_mac)
	}
}
