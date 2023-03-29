package main

import (
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

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	setupTenantConfig()

	r.Post("/next", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("hi"))
	})

	http.ListenAndServe(":8080", r)
}
