package gowebserver

import "github.com/gin-gonic/gin"

func GetFaviconIco(c *gin.Context) {
	c.File("./urmid.svg")
}

func GreetUser(c *gin.Context) {
	c.IndentedJSON(200, "Hello World!")
}

func GreetUserByName(c *gin.Context) {
	name := c.Param("name")
	c.String(200, "Hello %s", name)
}
