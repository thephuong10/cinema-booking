package converters

import (
	"showtimeservice/models/entities"
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
