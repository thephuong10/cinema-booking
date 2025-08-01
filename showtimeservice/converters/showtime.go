package converters

import (
	"github.com/google/uuid"
	"showtimeservice/models/entities"
	"showtimeservice/models/requests"
	"showtimeservice/models/responses"
)

func ToResponse(showTimes []entities.ShowTimes) []responses.ShowTime {
	var res []responses.ShowTime

	for _, st := range showTimes {
		res = append(res, responses.ShowTime{
			ID:        st.ID,
			CinemaID:  st.CinemaID,
			RoomID:    st.RoomID,
			MovieID:   st.MovieID,
			Status:    st.Status,
			Price:     st.Price,
			StartTime: st.StartTime,
			EndTime:   st.EndTime,
		})
	}

	return res

}

func ConvertShowTimeEntityToResponse(st *entities.ShowTimes) *responses.ShowTime {
	return &responses.ShowTime{
		ID:        st.ID,
		CinemaID:  st.CinemaID,
		RoomID:    st.RoomID,
		MovieID:   st.MovieID,
		Status:    st.Status,
		Price:     st.Price,
		StartTime: st.StartTime,
		EndTime:   st.EndTime,
	}

}

func ConvertCreateShowTimeRequestToEntity(req *requests.CreateShowTime) entities.ShowTimes {
	return entities.ShowTimes{
		ID:        uuid.New().String(),
		CinemaID:  req.CinemaID,
		RoomID:    req.RoomID,
		MovieID:   req.MovieID,
		Status:    req.Status,
		Price:     req.Price,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}
}
