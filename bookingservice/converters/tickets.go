package converters

import (
	"bookingservice/models/entities"
	"bookingservice/models/responses"
)

func ConvertTicketEntityToResponse(entities []entities.Ticket) []responses.TicketResponse {
	var res []responses.TicketResponse

	for _, entity := range entities {
		res = append(res, responses.TicketResponse{
			Id:         entity.ID,
			Row:        entity.Row,
			Column:     entity.Column,
			Price:      entity.Price,
			ShowTimeId: entity.ShowTimeId,
		})
	}

	return res

}
