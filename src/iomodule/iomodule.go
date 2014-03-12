package iomodule

import(
   "driver"
   "fmt"
   "time"
)

func IoManager(order_queue chan driver.Data, command_list chan driver.Data, order_list chan driver.Data, cost chan driver.Data, remove_order chan driver.Data, remove_command chan driver.Data){
   var order_queue_copy driver.Data
   order_queue_copy.Array = [8]int{0,0,0,0,0,0,0,0}
   var order_list_copy driver.Data
   order_list_copy.Array = [8]int{0,0,0,0,0,0,0,0}
   var command_list_copy driver.Data
   command_list_copy.Array = [8]int{0,0,0,0,0,0,0,0}
   
   remove_orders := [8]int{0,0,0,0,0,0,0,0}
   remove_commands := [8]int{0,0,0,0,0,0,0,0}
   
   go func(){
      for{
         select{
         case data := <- remove_order:
            remove_orders = data.Array
         case data := <- remove_command:
            remove_commands = data.Array          
         }
         time.Sleep(1*time.Millisecond)
      }
   }()
   
   for {
   
      
      i := 0
      for i<4{
         if driver.GetButtonSignal("command", i) == 1{
            command_list_copy.Array[i] = 1
         } 
         if driver.GetButtonSignal("down", i) == 1{
            order_list_copy.Array[2*i] = 1
         }
         if driver.GetButtonSignal("up", i) == 1{
            fmt.Println("knapp nr: ", i, " er blitt satt til 1")
            order_list_copy.Array[2*i + 1] = 1
         }
         i += 1
         
      }
      
      // Panel thing
      
      i= 0 
      for (i < 4){
         driver.SetButtonLamp("command", i, command_list_copy.Array[i])
         if (i == driver.GetFloor()){
            driver.SetFloorIndicator(i)
         }
         if (i > 0){   
            driver.SetButtonLamp("down", i , order_list_copy.Array[i*2])
         }
         if (i < 3){
            driver.SetButtonLamp("up", i, order_list_copy.Array[i*2+1])
         }
         i += 1
      }
      //fmt.Printf("Deadlock yo?\n")
      i = 0
      for i<4{
         if remove_commands[i] == 1{
            command_list_copy.Array[i] = 0
         } 
         if remove_orders[2*i] == 1{
            order_list_copy.Array[2*i] = 0
         }
         if remove_orders[2*i+1] == 1{
            order_list_copy.Array[2*i+1] =  0
         }
         i += 1
      }
      i = 0
      
      for i<3{
         select{
         case data := <- order_list:
            order_list_copy = data
         case data := <- command_list:
            command_list_copy = data
         default:
            order_list <- order_list_copy
            command_list <- command_list_copy
         
         }
         i += 1 
      }  
      fmt.Printf("Nope.yo\n")
   
   time.Sleep(1*time.Millisecond)
   }

}

