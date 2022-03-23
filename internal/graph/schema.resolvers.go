package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"rider-service/internal/core/domain"
	"rider-service/internal/graph/generated"
	"rider-service/internal/graph/model"
	"strconv"
)

func (r *mutationResolver) CreateRider(ctx context.Context, input model.RiderInput) (*model.Rider, error) {
	rider, err := r.RiderService.Create(input.Name, input.Status)

	if err != nil {
		return nil, err
	}

	riderModel := model.Rider{
		ID:     strconv.Itoa(int(rider.ID)),
		Name:   rider.Name,
		Status: rider.Status,
	}

	return &riderModel, err
}

func (r *mutationResolver) UpdateRider(ctx context.Context, id string, input *model.RiderInput) (*model.Rider, error) {
	rider, err := r.RiderService.Update(id, input.Name, input.Status)

	if err != nil {
		return nil, err
	}

	riderModel := model.Rider{
		ID:     strconv.Itoa(int(rider.ID)),
		Name:   rider.Name,
		Status: rider.Status,
	}

	return &riderModel, err
}

func (r *mutationResolver) UpdateLocation(ctx context.Context, id string, input *model.LocationInput) (*model.Rider, error) {
	location, err := domain.NewLocation(input.Latitude, input.Longitude)

	if err != nil {
		return nil, err
	}

	rider, err := r.RiderService.UpdateLocation(id, location)

	if err != nil {
		return nil, err
	}

	riderModel := model.Rider{
		ID:       strconv.Itoa(int(rider.ID)),
		Name:     rider.Name,
		Status:   rider.Status,
		Location: &model.Location{Latitude: rider.Location.Latitude, Longitude: rider.Location.Longitude},
	}

	return &riderModel, err
}

func (r *queryResolver) Riders(ctx context.Context) ([]*model.Rider, error) {
	riders, err := r.RiderService.GetAll()

	if err != nil {
		return nil, err
	}

	var riderModels []*model.Rider

	for _, v := range riders {
		riderModels = append(riderModels, &model.Rider{
			ID:       strconv.Itoa(int(v.ID)),
			Name:     v.Name,
			Status:   v.Status,
			Location: &model.Location{Latitude: v.Location.Latitude, Longitude: v.Location.Longitude},
		})
	}

	return riderModels, nil
}

func (r *queryResolver) Rider(ctx context.Context, id string) (*model.Rider, error) {
	rider, err := r.RiderService.Get(id)

	if err != nil {
		return nil, err
	}

	riderModel := model.Rider{
		ID:       strconv.Itoa(int(rider.ID)),
		Name:     rider.Name,
		Status:   rider.Status,
		Location: &model.Location{Latitude: rider.Location.Latitude, Longitude: rider.Location.Longitude},
	}

	return &riderModel, err
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
