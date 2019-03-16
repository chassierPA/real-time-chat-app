package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/graphql-go/graphql"
)

//Message structure of messaging
type Message struct {
	ID  string `json:id`
	Msg string `json:message`
}

//Messages stock of messages
var Messages []Message = []Message{
	Message{
		ID:  "1",
		Msg: "premier messsage",
	},
	Message{
		ID:  "2",
		Msg: "second message",
	},
}

var ID int = 2

func main() {
	msgType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Message",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"msg": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"messages": &graphql.Field{
				Type: graphql.NewList(msgType),
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return Messages, nil
				},
			},
		},
	})

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"addMessage": &graphql.Field{
				Type: graphql.NewList(msgType),
				Args: graphql.FieldConfigArgument{
					"msg": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					msg := params.Args["msg"].(string)
					ID += 1
					idty := strconv.Itoa(ID)
					Messages = append(Messages, Message{idty, msg})
					return Messages, nil
				},
			},
		},
	})

	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})

	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: r.URL.Query().Get("query"),
		})
		json.NewEncoder(w).Encode(result)
	})
	http.ListenAndServe(":12345", nil)
}

func FilterMsg(messages []Message, id string) Message {
	for _, value := range messages {
		if value.ID == id {
			return value
		}
	}
	return Message{}
}
