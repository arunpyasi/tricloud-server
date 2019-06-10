package broker

import (
	"fmt"
	"time"
)

func (h *Hub) debugg() {
	for {
		time.Sleep(5 * time.Second)
		fmt.Printf(" ALL AGENTS %+v \n %+v \n ALL Clients %+v", h.ListOfAgents, h.AllAgentConns, h.AllUserConns)
	}

}
