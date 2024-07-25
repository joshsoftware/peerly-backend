package badges

import (
	"context"
	"testing"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/dto"
	"github.com/joshsoftware/peerly-backend/internal/repository"
	"github.com/joshsoftware/peerly-backend/internal/repository/mocks"
	"github.com/stretchr/testify/mock"
)

func TestListBadges(t *testing.T) {
	badgeRepo := mocks.NewBadgesStorer(t)
	service := NewService(badgeRepo)

	tests := []struct {
		name            string
		context         context.Context
		setup           func(badgeMock *mocks.BadgesStorer)
		isErrorExpected bool
	}{
		{
			name:    "Success",
			context: context.Background(),
			setup: func(badgeMock *mocks.BadgesStorer) {
				badgeMock.On("ListBadges", mock.Anything, nil).Return([]repository.Badge{
					{
						ID:           1,
						Name:         "Gold",
						RewardPoints: 1000,
					},
					{
						ID:           2,
						Name:         "Platinum",
						RewardPoints: 2000,
					},
				}, nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name:    "Error in list badges",
			context: context.Background(),
			setup: func(badgeMock *mocks.BadgesStorer) {
				badgeMock.On("ListBadges", mock.Anything, nil).Return([]repository.Badge{}, apperrors.InternalServerError).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(badgeRepo)

			// test service
			_, err := service.ListBadges(test.context)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test Failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
		})
	}
}

func TestGetBadge(t *testing.T) {
	badgeRepo := mocks.NewBadgesStorer(t)
	service := NewService(badgeRepo)

	tests := []struct {
		name            string
		context         context.Context
		badgeID         int8
		setup           func(badgeMock *mocks.BadgesStorer)
		isErrorExpected bool
	}{
		{
			name:    "Success",
			context: context.Background(),
			badgeID: 1,
			setup: func(badgeMock *mocks.BadgesStorer) {
				badgeMock.On("GetBadge", mock.Anything, nil, int8(1)).Return(repository.Badge{
					ID:           1,
					Name:         "Gold",
					RewardPoints: 1000,
				}, nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name:    "Error in fetch badge",
			badgeID: 1,
			context: context.Background(),
			setup: func(badgeMock *mocks.BadgesStorer) {
				badgeMock.On("GetBadge", mock.Anything, nil, int8(1)).Return(repository.Badge{}, apperrors.InternalServerError).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(badgeRepo)

			// test service
			_, err := service.GetBadge(test.context, test.badgeID)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test Failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
		})
	}
}

func TestDeleteBadge(t *testing.T) {
	badgeRepo := mocks.NewBadgesStorer(t)
	service := NewService(badgeRepo)

	tests := []struct {
		name            string
		context         context.Context
		badgeID         int8
		setup           func(badgeMock *mocks.BadgesStorer)
		isErrorExpected bool
	}{
		{
			name:    "Success",
			context: context.Background(),
			badgeID: 1,
			setup: func(badgeMock *mocks.BadgesStorer) {
				badgeMock.On("DeleteBadge", mock.Anything, nil, int8(1)).Return(nil).Once()
			},
			isErrorExpected: false,
		},
		{
			name:    "Error in list corevalues",
			badgeID: 1,
			context: context.Background(),
			setup: func(badgeMock *mocks.BadgesStorer) {
				badgeMock.On("DeleteBadge", mock.Anything, nil, int8(1)).Return(apperrors.BadgeNotFound).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(badgeRepo)

			// test service
			err := service.DeleteBadge(test.context, test.badgeID)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test Failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
		})
	}
}

func TestUpdateBadge(t *testing.T) {
	badgeRepo := mocks.NewBadgesStorer(t)
	service := NewService(badgeRepo)

	tests := []struct {
		name             string
		context          context.Context
		badgeUpdatedInfo dto.Badge
		setup            func(badgeMock *mocks.BadgesStorer)
		isErrorExpected  bool
	}{
		{
			name:    "Success",
			context: context.Background(),
			badgeUpdatedInfo: dto.Badge{
				ID:           int8(1),
				Name:         "Gold",
				RewardPoints: int16(1000),
			},
			setup: func(badgeMock *mocks.BadgesStorer) {

				badgeMock.On("GetBadge", mock.Anything, nil, int8(1)).Return(repository.Badge{
					ID:           1,
					Name:         "Gold",
					RewardPoints: 1000,
				}, nil).Once()

				badgeMock.On("GetBadgeByName", mock.Anything, nil, "Gold").Return(int8(0)).Once()

				badgeMock.On("GetBadgeByRewardPoints", mock.Anything, nil, int16(1000)).Return(int8(0)).Once()

				badgeMock.On("UpdateBadge", mock.Anything, nil, dto.Badge{
					ID:           1,
					Name:         "Gold",
					RewardPoints: 1000,
				}).Return(repository.Badge{
					ID:           1,
					Name:         "Gold",
					RewardPoints: 1000,
				}, nil).Once()

			},
			isErrorExpected: false,
		},
		{
			name:    "Error badge not found",
			context: context.Background(),
			badgeUpdatedInfo: dto.Badge{
				ID:           int8(1),
				Name:         "Gold",
				RewardPoints: int16(1000),
			},
			setup: func(badgeMock *mocks.BadgesStorer) {
				badgeMock.On("GetBadge", mock.Anything, nil, int8(1)).Return(repository.Badge{}, apperrors.BadgeNotFound).Once()
			},
			isErrorExpected: true,
		},
		{
			name:    "Error badge name already exists",
			context: context.Background(),
			badgeUpdatedInfo: dto.Badge{
				ID:           int8(1),
				Name:         "Gold",
				RewardPoints: int16(1000),
			},
			setup: func(badgeMock *mocks.BadgesStorer) {
				badgeMock.On("GetBadge", mock.Anything, nil, int8(1)).Return(repository.Badge{
					ID:           1,
					Name:         "Gold",
					RewardPoints: 1000,
				}, nil).Once()

				badgeMock.On("GetBadgeByName", mock.Anything, nil, "Gold").Return(int8(2)).Once()
			},
			isErrorExpected: true,
		},
		{
			name:    "Error badge reward points already exists",
			context: context.Background(),
			badgeUpdatedInfo: dto.Badge{
				ID:           int8(1),
				Name:         "Gold",
				RewardPoints: int16(1000),
			},
			setup: func(badgeMock *mocks.BadgesStorer) {
				badgeMock.On("GetBadge", mock.Anything, nil, int8(1)).Return(repository.Badge{
					ID:           1,
					Name:         "Gold",
					RewardPoints: 1000,
				}, nil).Once()

				badgeMock.On("GetBadgeByName", mock.Anything, nil, "Gold").Return(int8(0)).Once()

				badgeMock.On("GetBadgeByRewardPoints", mock.Anything, nil, int16(1000)).Return(int8(2)).Once()
			},
			isErrorExpected: true,
		},
		{
			name:    "Error in updating badge",
			context: context.Background(),
			badgeUpdatedInfo: dto.Badge{
				ID:           int8(1),
				Name:         "Gold",
				RewardPoints: int16(1000),
			},
			setup: func(badgeMock *mocks.BadgesStorer) {
				badgeMock.On("GetBadge", mock.Anything, nil, int8(1)).Return(repository.Badge{
					ID:           1,
					Name:         "Gold",
					RewardPoints: 1000,
				}, nil).Once()

				badgeMock.On("GetBadgeByName", mock.Anything, nil, "Gold").Return(int8(0)).Once()

				badgeMock.On("GetBadgeByRewardPoints", mock.Anything, nil, int16(1000)).Return(int8(0)).Once()

				badgeMock.On("UpdateBadge", mock.Anything, nil, dto.Badge{
					ID:           1,
					Name:         "Gold",
					RewardPoints: 1000,
				}).Return(repository.Badge{}, apperrors.BadgeNotFound).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(badgeRepo)

			// test service
			_, err := service.UpdateBadge(test.context, test.badgeUpdatedInfo)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test Failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
		})
	}
}

func TestCreateBadge(t *testing.T) {
	badgeRepo := mocks.NewBadgesStorer(t)
	service := NewService(badgeRepo)

	tests := []struct {
		name            string
		context         context.Context
		badge           dto.Badge
		setup           func(badgeMock *mocks.BadgesStorer)
		isErrorExpected bool
	}{
		{
			name:    "Success",
			context: context.Background(),
			badge: dto.Badge{
				ID:           int8(1),
				Name:         "Gold",
				RewardPoints: int16(1000),
			},
			setup: func(badgeMock *mocks.BadgesStorer) {

				badgeMock.On("GetBadgeByName", mock.Anything, nil, "Gold").Return(int8(0)).Once()

				badgeMock.On("GetBadgeByRewardPoints", mock.Anything, nil, int16(1000)).Return(int8(0)).Once()

				badgeMock.On("CreateBadge", mock.Anything, nil, dto.Badge{
					ID:           1,
					Name:         "Gold",
					RewardPoints: 1000,
				}).Return(repository.Badge{
					ID:           1,
					Name:         "Gold",
					RewardPoints: 1000,
				}, nil).Once()

			},
			isErrorExpected: false,
		},
		{
			name:    "Error badge name already exists",
			context: context.Background(),
			badge: dto.Badge{
				ID:           int8(1),
				Name:         "Gold",
				RewardPoints: int16(1000),
			},
			setup: func(badgeMock *mocks.BadgesStorer) {
				badgeMock.On("GetBadgeByName", mock.Anything, nil, "Gold").Return(int8(2)).Once()
			},
			isErrorExpected: true,
		},
		{
			name:    "Error badge reward points already exists",
			context: context.Background(),
			badge: dto.Badge{
				ID:           int8(1),
				Name:         "Gold",
				RewardPoints: int16(1000),
			},
			setup: func(badgeMock *mocks.BadgesStorer) {
				badgeMock.On("GetBadgeByName", mock.Anything, nil, "Gold").Return(int8(0)).Once()

				badgeMock.On("GetBadgeByRewardPoints", mock.Anything, nil, int16(1000)).Return(int8(2)).Once()
			},
			isErrorExpected: true,
		},
		{
			name:    "Error in creating badge",
			context: context.Background(),
			badge: dto.Badge{
				ID:           int8(1),
				Name:         "Gold",
				RewardPoints: int16(1000),
			},
			setup: func(badgeMock *mocks.BadgesStorer) {

				badgeMock.On("GetBadgeByName", mock.Anything, nil, "Gold").Return(int8(0)).Once()

				badgeMock.On("GetBadgeByRewardPoints", mock.Anything, nil, int16(1000)).Return(int8(0)).Once()

				badgeMock.On("CreateBadge", mock.Anything, nil, dto.Badge{
					ID:           1,
					Name:         "Gold",
					RewardPoints: 1000,
				}).Return(repository.Badge{}, apperrors.BadgeNotFound).Once()
			},
			isErrorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(badgeRepo)

			// test service
			_, err := service.CreateBadge(test.context, test.badge)

			if (err != nil) != test.isErrorExpected {
				t.Errorf("Test Failed, expected error to be %v, but got err %v", test.isErrorExpected, err != nil)
			}
		})
	}
}
