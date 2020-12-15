package structs

import "fp-dynamic-elements-manager-controller/internal/queue/structs"

type JsonExportResults struct {
	Results []structs.ListElement `json:"results"`
}
