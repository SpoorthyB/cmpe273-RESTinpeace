package main

import (
	"fmt"
	"log"
	"time"
//	"encoding/json"
//	"reflect"
	"net/http"
	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
)


func onStateChanged(d gatt.Device, s gatt.State) {
	fmt.Println("State:", s)
	switch s {
	case gatt.StatePoweredOn:
		fmt.Println("scanning...")
		d.Scan([]gatt.UUID{}, false)
		return
	default:
		d.StopScanning()
	}
}

func onPeriphDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	fmt.Printf("\nPeripheral ID:%s, NAME:(%s)\n", p.ID(), p.Name())
	fmt.Println("  Local Name        =", a.LocalName)
	fmt.Println("  TX Power Level    =", a.TxPowerLevel)
	fmt.Println("  Manufacturer Data =", a.ManufacturerData)
	fmt.Println("  Service Data      =", a.ServiceData)
	if(len(a.ManufacturerData)!=0){
	      go Request(a.ManufacturerData)
	}
}
func Request(ManufacturerData []uint8){
	t := time.Now()
	client := &http.Client{}
        fmt.Println(string(ManufacturerData[2:]))
        id:=string(ManufacturerData[2:])
                if (ManufacturerData[0]==0){
                fmt.Println("Checkout")
                req, _ := http.NewRequest("PUT", "http://52.37.72.212:3005/update/"+id+"/left", nil)
                resp, _ := client.Do(req)
                fmt.Println(resp)
                }
                if(ManufacturerData[0]==1){
		fmt.Println("TIME: "+t.Format("2006-01-02 15:04:05"))
                fmt.Println("CheckIn")
                req, _ := http.NewRequest("PUT", "http://52.37.72.212:3005/update/"+id+"/attended", nil)
                resp, _ := client.Do(req)
                fmt.Println(resp)
                }
}




func main() {
	d, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		log.Fatalf("Failed to open device, err: %s\n", err)
		return
	}

	// Register handlers.
	d.Handle(gatt.PeripheralDiscovered(onPeriphDiscovered))
	fmt.Println(d)
	d.Init(onStateChanged)
	select {}
}
