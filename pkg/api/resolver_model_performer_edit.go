package api

import (
	"context"

	"github.com/stashapp/stash-box/pkg/dataloader"
	"github.com/stashapp/stash-box/pkg/models"
	"github.com/stashapp/stash-box/pkg/utils"

	"github.com/gofrs/uuid"
)

type performerEditResolver struct{ *Resolver }

func (r *performerEditResolver) Gender(ctx context.Context, obj *models.PerformerEdit) (*models.GenderEnum, error) {
	var ret models.GenderEnum
	if obj.Gender == nil || !utils.ResolveEnumString(*obj.Gender, &ret) {
		return nil, nil
	}

	return &ret, nil
}

func (r *performerEditResolver) HairColor(ctx context.Context, obj *models.PerformerEdit) (*models.HairColorEnum, error) {
	var ret models.HairColorEnum
	if obj.HairColor == nil || !utils.ResolveEnumString(*obj.HairColor, &ret) {
		return nil, nil
	}

	return &ret, nil
}

func (r *performerEditResolver) EyeColor(ctx context.Context, obj *models.PerformerEdit) (*models.EyeColorEnum, error) {
	var ret models.EyeColorEnum
	if obj.EyeColor == nil || !utils.ResolveEnumString(*obj.EyeColor, &ret) {
		return nil, nil
	}

	return &ret, nil
}

func (r *performerEditResolver) Ethnicity(ctx context.Context, obj *models.PerformerEdit) (*models.EthnicityEnum, error) {
	var ret models.EthnicityEnum
	if obj.Ethnicity == nil || !utils.ResolveEnumString(*obj.Ethnicity, &ret) {
		return nil, nil
	}

	return &ret, nil
}

func (r *performerEditResolver) BreastType(ctx context.Context, obj *models.PerformerEdit) (*models.BreastTypeEnum, error) {
	var ret models.BreastTypeEnum
	if obj.BreastType == nil || !utils.ResolveEnumString(*obj.BreastType, &ret) {
		return nil, nil
	}

	return &ret, nil
}

func (r *performerEditResolver) AddedImages(ctx context.Context, obj *models.PerformerEdit) ([]*models.Image, error) {
	if len(obj.AddedImages) == 0 {
		return nil, nil
	}

	var uuids []uuid.UUID
	for _, id := range obj.AddedImages {
		imageID, _ := uuid.FromString(id)
		uuids = append(uuids, imageID)
	}
	images, errors := dataloader.For(ctx).ImageByID.LoadAll(uuids)
	for _, err := range errors {
		if err != nil {
			return nil, err
		}
	}
	return images, nil
}

func (r *performerEditResolver) RemovedImages(ctx context.Context, obj *models.PerformerEdit) ([]*models.Image, error) {
	if len(obj.RemovedImages) == 0 {
		return nil, nil
	}

	var uuids []uuid.UUID
	for _, id := range obj.RemovedImages {
		imageID, _ := uuid.FromString(id)
		uuids = append(uuids, imageID)
	}
	images, errors := dataloader.For(ctx).ImageByID.LoadAll(uuids)
	for _, err := range errors {
		if err != nil {
			return nil, err
		}
	}
	return images, nil
}
