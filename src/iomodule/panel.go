package iomodule

import(
   "driver"
)

func PanelLights(order_list chan driver.Data, command_list chan driver.Data){
   i:= 0
   for (i < 4){
      order_list_copy := <- order_list
      order_list <- order_list_copy
      
      command_list_copy := <- command_list
      command_list <- command_list_copy
      
      driver.SetButtonLamp("command", i, command_list_copy.Array[i])
      
      if (i > 0){   
         driver.SetButtonLamp("down", i*2, order_list_copy.Array[i*2])
      }
      
      if (i < 3){
         driver.SetButtonLamp("up", i*2 +1, order_list_copy.Array[i*2+1])
      }
      i += 1
   }
}
