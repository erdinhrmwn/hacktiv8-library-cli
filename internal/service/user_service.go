package service

import (
	"context"
	"fmt"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/model"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/repository"
	"github.com/erdinhrmwn/hacktiv8-library-cli/utils"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	user, err := s.userRepository.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, fmt.Errorf("user tidak ditemukan")
	}
	return user, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, nil
	}
	return user, nil
}

func (s *UserService) GetUserByRole(ctx context.Context, role string) ([]model.User, error) {
	users, err := s.userRepository.GetUsersByRole(ctx, role)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) CreateUser(ctx context.Context, user *model.User) error {
	if user.Name == "" {
		return fmt.Errorf("nama tidak boleh kosong")
	}
	if user.Email == "" {
		return fmt.Errorf("email tidak boleh kosong")
	}
	if user.Role != "staff" && user.Role != "visitor" {
		return fmt.Errorf("role tidak valid")
	}

	hashed, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed

	return s.userRepository.CreateUser(ctx, user)
}

func (s *UserService) UpdateUser(ctx context.Context, user *model.User) error {
	if user.Name == "" {
		return fmt.Errorf("nama tidak boleh kosong")
	}
	if user.Email == "" {
		return fmt.Errorf("email tidak boleh kosong")
	}
	if user.Role != "staff" && user.Role != "visitor" {
		return fmt.Errorf("role tidak valid")
	}

	existing, err := s.userRepository.GetUserByID(ctx, user.ID)
	if err != nil {
		return err
	}
	if existing.ID == 0 {
		return fmt.Errorf("user tidak ditemukan")
	}

	user.Password = existing.Password

	return s.userRepository.UpdateUser(ctx, user)
}

func (s *UserService) ChangePassword(ctx context.Context, id int, oldPassword, newPassword string) error {
	if newPassword == "" {
		return fmt.Errorf("password baru tidak boleh kosong")
	}
	if len(newPassword) < 5 {
		return fmt.Errorf("password minimal 5 karakter")
	}

	user, err := s.userRepository.GetUserByID(ctx, id)
	if err != nil {
		return err
	}
	if user.ID == 0 {
		return fmt.Errorf("user tidak ditemukan")
	}

	if !utils.VerifyPassword(oldPassword, user.Password) {
		return fmt.Errorf("password lama salah")
	}

	hashed, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	user.Password = hashed

	return s.userRepository.UpdateUser(ctx, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	existing, err := s.userRepository.GetUserByID(ctx, id)
	if err != nil {
		return err
	}
	if existing.ID == 0 {
		return fmt.Errorf("user tidak ditemukan")
	}

	return s.userRepository.DeleteUser(ctx, id)
}
