package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	command    = "/usr/sbin/aprontest"
	ledCommand = "/usr/sbin/set_rgb"
)

func setRGB(args []string) {

	_, err := exec.Command(ledCommand, args...).Output()
	if err != nil {
		fmt.Println("Error setting LED")
	}

}

func executeAction(c *gin.Context) {

	var args []string
	for i := 0; i < 10; i++ {
		operation := c.Query(strconv.Itoa(i))
		if len(operation) > 0 {
			args = append(args, operation)
		} else {
			break
		}
	}
	fmt.Println(args)
	//Turn LED Blue to Indicate event.
	setRGB([]string{"0", "255", "255"})
	time.Sleep(1 * time.Second)

	statusCode := 200
	out, err := exec.Command(command, args...).Output()
	if err != nil {
		statusCode = 500
		//Flash LED red to indicate Error
		setRGB([]string{"255", "0", "0"})
		time.Sleep(1 * time.Second)
	}

	//Turn LED back to green
	setRGB([]string{"0", "255", "0"})

	c.JSON(statusCode, gin.H{
		"status": "received",
		"output": string(out[:]),
	})

}

func main() {
	router := gin.Default()

	v1 := router.Group("/v1")
	{
		v1.GET("/action", executeAction)
	}

	router.Run(":4141")
}
