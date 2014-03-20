package heis

import(
   "driver"
   "time"  
   "math"
  // "fmt"
)

func HeisInit()(int, int, int){
   direction := 0;
   for driver.GetFloor() == -1 {
      driver.SetSpeed(-300)
   }
   driver.SetSpeed(0)
   current_floor := driver.GetFloor()
   destination := -1
   return direction, current_floor, destination
}

func Heis(order_list chan driver.Data, command_list chan driver.Data, cost chan driver.Data, remove_order chan driver.Data, remove_command chan driver.Data, elevator_number chan int){ 
   //Initialize variables. 
   direction, current_floor, destination := HeisInit()
   //var cost_copy driver.Data
   var command_list_copy driver.Data
   var command_list_temp driver.Data
   var order_list_copy driver.Data
   var order_list_temp driver.Data
   var remove_orders driver.Data
   var remove_commands driver.Data
   var cost_copy driver.Data
   var elevator_nr int
   order_list_copy.Array = [8]int{0,0,0,0,0,0,0,0}
   command_list_copy.Array = [8]int{0,0,0,0,0,0,0,0}
   
   elevator_nr = -1
   //Starts a gorutine for continously reading relevant channels. 
   go func(){
      for{
         select{
         case data := <- order_list:
            order_list_temp = data
            order_list_copy = order_list_temp
         case data := <- command_list:
            command_list_temp = data
            command_list_copy = command_list_temp
         case data := <- elevator_number:
         	elevator_nr = data
         }
         time.Sleep(1*time.Millisecond)
      }
   }()
   
   for{
      //direction is initialized to zero, this function returns the first found destination if the direction
      // is zero, and optimalizes the destination if the direction is positive or negative. 
      destination = GetDestination(direction, current_floor, order_list_copy.Array, command_list_copy.Array, elevator_nr)
 	  
      // decides direction required to reach destination from current floor.
      direction = GetDirection(destination, current_floor)
      driver.SetSpeed(direction*300)
      
      if(driver.GetFloor() != -1){
      	cost_copy.Array = CostFunction(current_floor, direction, destination, order_list_copy.Array, elevator_nr)
      	cost <- cost_copy
      }	//Loop where the elevator is in running mode.  
      for(destination != -1){
      	// Updates current floor and optimizes the destination
         destination = GetDestination(direction, current_floor, order_list_copy.Array, command_list_copy.Array, elevator_nr)
         floor := driver.GetFloor() 
         if(floor != -1){
            current_floor = floor
         }
         if(driver.GetFloor() != -1){
      		cost_copy.Array = CostFunction(current_floor, direction, destination, order_list_copy.Array, elevator_nr)
      		cost <- cost_copy
      	 }	
         //If sentence with the requirements for a stop. 
         if( (current_floor == driver.GetFloor()) && ((direction==-1 && order_list_copy.Array[2*current_floor]==elevator_nr) || (direction==1 && order_list_copy.Array[2*current_floor+1]==elevator_nr) || command_list_copy.Array[current_floor] == 1 || (destination == current_floor))){
        
            //stopping the elevator
            driver.SetSpeed(-1* direction*300)
            time.Sleep(15*time.Millisecond)
            driver.SetSpeed(0)
          	
          	//Sending what orders have been accomplished. 
            remove_orders.Array, remove_commands.Array = RemoveOrders(current_floor, direction, destination)
            remove_order <- remove_orders
            remove_command <- remove_commands
            
            //opening/closing doors 
            driver.SetDoorLamp(1)
            time.Sleep(3*time.Second)
            driver.SetDoorLamp(0)
            
            //resets the destination if it is the current floor.
            if current_floor == destination{
               destination = -1
            }
            //was apparantly neccesary with a break here. 
            break
         }
         time.Sleep(1*time.Millisecond)
      }
      time.Sleep(1*time.Millisecond) 
   }
}








func GetDirection(destination int, current_floor int)(int){
   direction := 0
   if (destination == -1){
      return direction
   }else if(destination > current_floor){
      direction = 1  
   }else if(destination < current_floor){
      direction = -1
   }
   return direction
}


func GetDestination(direction int, current_floor int, order_list [8]int, command_list [8]int, elevator_nr int)(int){
   var i int
   candidate := -1
   if(direction == 1){
      i = 3
      for(i >= current_floor){
         if (order_list[i*2+1] == elevator_nr || command_list[i] == 1){
            if i > candidate{
            	candidate =  i
            }
         }else if(order_list[i*2] == elevator_nr){
         	if i > candidate{
         		candidate = i
         	}
         }
         i -= 1
      }
      return candidate
   }else if (direction == -1){
      i = 0
      for(i <= current_floor){
         if (order_list[i*2] == elevator_nr || command_list[i] == 1){
            return i
         }
         i += 1
      }
      return -1
   }else{
      i = 0
      for(i < 4){
         if (order_list[i*2] == elevator_nr || order_list[i*2+1] == elevator_nr || command_list[i] == 1){
            return i
         }
         i += 1
      }
      return -1         
   }

}


func CostFunction(current_floor int,direction int, destination int, order_list [8]int, elevator_nr int)([8]int){
   i := 0
   var cost [8]int
   for i<8{
      if (direction == 0){
         cost[i] = int(math.Abs(float64(i/2 - current_floor)))
      }else if(direction == 1){
         if(i%2 == 1 && i/2 > current_floor){
            cost[i] = i/2 - current_floor - 1
         }else if (i%2 == 1 && i/2 <= current_floor || i%2 == 0){
            cost[i] = int (math.Abs(float64(i/2 - destination)) + math.Abs(float64(destination - current_floor - 1)))
         	if (destination != 3){
         		cost[i] += 1
         	}
         }else{
            cost[i] = 6
         }
      }else{
         if(i%2 == 0 && i/2 < current_floor){
            cost[i] =  current_floor - i/2 - 1
         }else if (i%2 == 0 && i/2 >= current_floor || i%2 ==  1){
            cost[i] = int(math.Abs(float64(i/2 - destination)) + math.Abs(float64(current_floor - destination - 1)))
         	if (destination != 0){
         		cost[i] += 1
         	}
         }else{
            cost[i] = 6
         }
      }
      
      i += 1
   }
   i = 0
   
   for i<8{
      if order_list[i] == elevator_nr{
         cost[i] = 0
      }
      i += 1
   }
   
   return cost
}


func RemoveOrders(current_floor int,direction int, destination int)([8]int,[8]int){
   remove_order := [8]int{0,0,0,0,0,0,0,0}
   remove_command :=[8]int{0,0,0,0,0,0,0,0}
   i := 0
   for (i < 4){
      if (current_floor == i){
         remove_command[i] = 1
         if (destination == i){
         	remove_command[i] = 1
         	remove_order[i*2] = 1
         	remove_order[i*2+1] = 1
         	
         }else if (direction == 1){
            remove_order[i*2+1] = 1
         } else if (direction == -1){
            remove_order[i*2] = 1
         } 
      }
      i +=1
   }
   return remove_order, remove_command
}

