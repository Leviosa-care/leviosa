package notification

import (
	"context"
	"fmt"
	"strings"

	authuserPorts "github.com/Leviosa-care/leviosa/backend/internal/authuser/ports"
	bookingPorts "github.com/Leviosa-care/leviosa/backend/internal/booking/ports"

	"github.com/google/uuid"
)

// InProcessUserFetcher fetches user details via the authuser UserService.
type InProcessUserFetcher struct {
	userService authuserPorts.UserService
}

// NewInProcessUserFetcher creates a new user fetcher backed by the authuser service.
func NewInProcessUserFetcher(userService authuserPorts.UserService) UserFetcher {
	return &InProcessUserFetcher{userService: userService}
}

func (f *InProcessUserFetcher) GetUserByID(ctx context.Context, userID uuid.UUID) (*UserInfo, error) {
	user, err := f.userService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("fetch user %s: %w", userID, err)
	}
	return &UserInfo{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Telephone,
	}, nil
}

// InProcessRoomFetcher fetches room details via the booking RoomService.
type InProcessRoomFetcher struct {
	roomService bookingPorts.RoomService
}

// NewInProcessRoomFetcher creates a new room fetcher backed by the booking service.
func NewInProcessRoomFetcher(roomService bookingPorts.RoomService) RoomFetcher {
	return &InProcessRoomFetcher{roomService: roomService}
}

func (f *InProcessRoomFetcher) GetRoom(ctx context.Context, roomID uuid.UUID) (*RoomInfo, error) {
	room, err := f.roomService.GetRoom(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("fetch room %s: %w", roomID, err)
	}
	return &RoomInfo{
		Name:       room.Name,
		BuildingID: room.BuildingID,
	}, nil
}

// InProcessBuildingFetcher fetches building details via the booking BuildingService.
type InProcessBuildingFetcher struct {
	buildingService bookingPorts.BuildingService
}

// NewInProcessBuildingFetcher creates a new building fetcher backed by the booking service.
func NewInProcessBuildingFetcher(buildingService bookingPorts.BuildingService) BuildingFetcher {
	return &InProcessBuildingFetcher{buildingService: buildingService}
}

func (f *InProcessBuildingFetcher) GetBuilding(ctx context.Context, buildingID uuid.UUID) (*BuildingInfo, error) {
	building, err := f.buildingService.GetBuildingByID(ctx, buildingID)
	if err != nil {
		return nil, fmt.Errorf("fetch building %s: %w", buildingID, err)
	}
	return &BuildingInfo{
		Name:    building.Name,
		Address: strings.TrimSpace(fmt.Sprintf("%s, %s %s", building.Address, building.PostalCode, building.City)),
	}, nil
}
