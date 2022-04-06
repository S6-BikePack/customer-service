package authorization

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

type RestAuthorization struct {
	context *gin.Context
	claims  map[string]interface{}
	id      string
}

func NewRest(context *gin.Context) *RestAuthorization {
	auth := RestAuthorization{
		context: context,
	}

	claimHeader := context.GetHeader("X-User-Claims")

	if claimHeader != "" {
		json.Unmarshal([]byte(claimHeader), &auth.claims)
	}

	auth.id = context.GetHeader("X-User-Id")

	return &auth
}

func (auth *RestAuthorization) AuthorizeAdmin() bool {
	v, exist := auth.claims["admin"]
	if exist && v == true {
		fmt.Println("AUTH: is admin")
	} else {
		fmt.Println("AUTH: is not admin")
	}
	return exist && v == true
}

func (auth *RestAuthorization) AuthorizeMatchingId(id string) bool {
	if auth.id == id {
		fmt.Println("AUTH: id's match")
	} else {
		fmt.Println("AUTH: id's do not match")
	}
	return auth.id == id
}
