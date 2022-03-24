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

func (r *mutationResolver) CreateRider(ctx context.Context, input model.RiderInput) (*domain.Rider, error) {
	rider, err := r.RiderService.Create(input.Name, input.Status)

	if err != nil {
		return nil, err
	}

	return &rider, err
}

func (r *mutationResolver) UpdateRider(ctx context.Context, id string, input *model.RiderInput) (*domain.Rider, error) {
	rider, err := r.RiderService.Update(id, input.Name, input.Status)

	if err != nil {
		return nil, err
	}
	return &rider, err
}

func (r *mutationResolver) UpdateLocation(ctx context.Context, id string, input *model.LocationInput) (*domain.Rider, error) {
	location, err := domain.NewLocation(input.Latitude, input.Longitude)

	if err != nil {
		return nil, err
	}

	rider, err := r.RiderService.UpdateLocation(id, location)

	if err != nil {
		return nil, err
	}

	return &rider, err
}

func (r *queryResolver) Riders(ctx context.Context) ([]*domain.Rider, error) {
	riders, err := r.RiderService.GetAll()

	if err != nil {
		return nil, err
	}

	var riderModels []*domain.Rider

	for i := range riders {
		riderModels = append(riderModels, &riders[i])
	}

	return riderModels, nil
}

func (r *queryResolver) Rider(ctx context.Context, id string) (*domain.Rider, error) {
	rider, err := r.RiderService.Get(id)

	if err != nil {
		return nil, err
	}

	return &rider, err
}

func (r *riderResolver) ID(ctx context.Context, obj *domain.Rider) (string, error) {
	return strconv.Itoa(int(obj.ID)), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Rider returns generated.RiderResolver implementation.
func (r *Resolver) Rider() generated.RiderResolver { return &riderResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type riderResolver struct{ *Resolver }
