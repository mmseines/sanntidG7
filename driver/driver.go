package driver  // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.c and driver.go
/*
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
#include "channels.h"
#include "elev.h"
*/
import "C"

func Init()(int){
    return C.elev_init()
}

func SetSpeed(speed int){
    C.elev_set_speed(speed)
}

func GetFloor()(int){
    return C.elev_get_floor_sensor_signal()
}

func GetButtonSignal(button string, floor int)(int){
    if button == "down" {
        return C.elev_get_button_signal(C.BUTTON_CALL_DOWN, int floor)
    }
    if button == "up" {
        return C.elev_get_button_signal(C.BUTTON_CALL_UP, int floor)
    }
    
        


} 
