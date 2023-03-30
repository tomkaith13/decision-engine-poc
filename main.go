package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// definition of a node
type Node struct {
	Children []Node
	Method   string
	ApiRoute string
	Id       string
}

type Tree struct {
	StartNode Node
	NodeMap   map[string]Node
}

func setupTenantConfig() Tree {
	// config
	// Choice A
	// specialty -> services -> questionnaire -> location -> practitioner -> timeslots -> book
	// Choice B
	// 										  -> timeslots -> book
	// Choice C
	// 										  -> timeslots -> practitioner -> book

	BookNode := Node{Children: nil, ApiRoute: "/appointment", Id: "book_node", Method: "POST"}

	ChoiceATimeslotChildren := []Node{BookNode}
	ChoiceATimeslotNode := Node{Children: ChoiceATimeslotChildren, ApiRoute: "/timeslots", Id: "choice_a_timeslot_node", Method: "GET"}

	ChoiceAPractitionerChildren := []Node{ChoiceATimeslotNode}
	ChoiceAPractitionerNode := Node{Children: ChoiceAPractitionerChildren, ApiRoute: "/practitioner", Id: "choice_a_pract_node", Method: "GET"}

	ChoiceALocationChildren := []Node{ChoiceAPractitionerNode}
	ChoiceALocationNode := Node{Children: ChoiceALocationChildren, ApiRoute: "/locations", Id: "choice_a_loc_node", Method: "GET"}

	// Choice B timeslots is same as Choice A

	ChoiceCPractitionerChildren := []Node{BookNode}
	ChoiceCPracitionerNode := Node{Children: ChoiceCPractitionerChildren, ApiRoute: "/practitioner", Id: "choice_c_pract_node", Method: "GET"}

	ChoiceCTimeslotsChildren := []Node{ChoiceCPracitionerNode}
	ChoiceCTimeslotsNode := Node{Children: ChoiceCTimeslotsChildren, ApiRoute: "/timeslots", Id: "choice_c_timeslot_node", Method: "GET"}

	// from this point on, all children converge regardless of choice
	QuestionnaireChildren := []Node{ChoiceALocationNode, ChoiceATimeslotNode, ChoiceCTimeslotsNode}
	QuestionnaireNode := Node{Children: QuestionnaireChildren, ApiRoute: "/questions", Method: "GET", Id: "q_node"}

	ServicesChildren := []Node{QuestionnaireNode}
	ServicesNode := Node{Children: ServicesChildren, ApiRoute: "/services", Id: "services_node", Method: "GET"}

	SpecialtyChildren := []Node{ServicesNode}
	SpecialtyNode := Node{Children: SpecialtyChildren, ApiRoute: "/specialties", Id: "specialty_node", Method: "GET"}

	tree := Tree{StartNode: SpecialtyNode}

	nodeMap := make(map[string]Node)
	nodeMap[SpecialtyNode.Id] = SpecialtyNode
	nodeMap[ServicesNode.Id] = ServicesNode
	nodeMap[QuestionnaireNode.Id] = QuestionnaireNode
	nodeMap[ChoiceCTimeslotsNode.Id] = ChoiceCTimeslotsNode
	nodeMap[ChoiceCPracitionerNode.Id] = ChoiceCPracitionerNode
	nodeMap[ChoiceALocationNode.Id] = ChoiceALocationNode
	nodeMap[ChoiceAPractitionerNode.Id] = ChoiceAPractitionerNode
	nodeMap[ChoiceATimeslotNode.Id] = ChoiceATimeslotNode
	tree.NodeMap = nodeMap

	return tree

}

type NextRequest struct {
	CurrentNodeId string `json:"current_node"`
	CurrentInputs any    `json:"input_kvp"`
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	tree := setupTenantConfig()

	r.Post("/next", func(w http.ResponseWriter, r *http.Request) {
		var bodyPayload NextRequest
		err := json.NewDecoder(r.Body).Decode(&bodyPayload)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Unable to decode body payload"))
			return
		}

		// Insert routing logic here
		node, ok := tree.NodeMap[bodyPayload.CurrentNodeId]
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("No such current node registered!!"))
			return
		}
		nextNodeId := ""
		if node.Id != "q_node" {
			if node.Children != nil {
				nextNodeId = node.Children[0].Id
			} else {
				nextNodeId = "nil"
			}
		} else {
			// choice logic
			choiceMap := bodyPayload.CurrentInputs.(map[string]any)
			if choiceMap["choice"] == "a" {
				nextNodeId = node.Children[0].Id
			} else if choiceMap["choice"] == "b" {
				nextNodeId = node.Children[1].Id
			} else {
				nextNodeId = node.Children[2].Id
			}
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("next node id: " + nextNodeId))
	})

	r.Get("/node/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		node, ok := tree.NodeMap[id]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("route: %s\nMethod: %s", node.ApiRoute, node.Method)))

	})

	http.ListenAndServe(":8080", r)
}
