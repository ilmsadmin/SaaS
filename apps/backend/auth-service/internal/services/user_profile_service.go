package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"zplus-saas/apps/backend/auth-service/internal/models"
)

type UserProfileService struct {
	db *sqlx.DB
}

func NewUserProfileService(db *sqlx.DB) *UserProfileService {
	return &UserProfileService{db: db}
}

// GetProfile gets user profile
func (s *UserProfileService) GetProfile(userID uuid.UUID) (*models.UserProfile, error) {
	var profile models.UserProfile
	err := s.db.Get(&profile, "SELECT * FROM user_profiles WHERE user_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("profile not found")
	}
	return &profile, nil
}

// GetUserWithProfile gets user with profile information
func (s *UserProfileService) GetUserWithProfile(userID uuid.UUID) (*models.User, *models.UserProfile, error) {
	var user models.User
	err := s.db.Get(&user, "SELECT * FROM users WHERE id = $1", userID)
	if err != nil {
		return nil, nil, fmt.Errorf("user not found")
	}

	var profile models.UserProfile
	err = s.db.Get(&profile, "SELECT * FROM user_profiles WHERE user_id = $1", userID)
	if err != nil {
		// Create default profile if not exists
		profile = models.UserProfile{
			ID:        uuid.New(),
			UserID:    userID,
			Language:  "en",
			Timezone:  "UTC",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		query := `
			INSERT INTO user_profiles (id, user_id, language, timezone, created_at, updated_at)
			VALUES (:id, :user_id, :language, :timezone, :created_at, :updated_at)`

		_, err = s.db.NamedExec(query, &profile)
		if err != nil {
			return &user, nil, fmt.Errorf("failed to create default profile: %w", err)
		}
	}

	return &user, &profile, nil
}

// UpdateProfile updates user profile
func (s *UserProfileService) UpdateProfile(userID uuid.UUID, req *models.UpdateProfileRequest) error {
	// Start transaction
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Update user basic info if provided
	if req.FirstName != nil || req.LastName != nil {
		setParts := []string{}
		args := []interface{}{}
		argIndex := 1

		if req.FirstName != nil {
			setParts = append(setParts, fmt.Sprintf("first_name = $%d", argIndex))
			args = append(args, *req.FirstName)
			argIndex++
		}

		if req.LastName != nil {
			setParts = append(setParts, fmt.Sprintf("last_name = $%d", argIndex))
			args = append(args, *req.LastName)
			argIndex++
		}

		setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
		args = append(args, time.Now())
		argIndex++

		args = append(args, userID)

		query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d",
			joinStrings(setParts, ", "), argIndex)

		_, err = tx.Exec(query, args...)
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
	}

	// Update profile
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Avatar != nil {
		setParts = append(setParts, fmt.Sprintf("avatar = $%d", argIndex))
		args = append(args, *req.Avatar)
		argIndex++
	}

	if req.Phone != nil {
		setParts = append(setParts, fmt.Sprintf("phone = $%d", argIndex))
		args = append(args, *req.Phone)
		argIndex++
	}

	if req.Address != nil {
		setParts = append(setParts, fmt.Sprintf("address = $%d", argIndex))
		args = append(args, *req.Address)
		argIndex++
	}

	if req.City != nil {
		setParts = append(setParts, fmt.Sprintf("city = $%d", argIndex))
		args = append(args, *req.City)
		argIndex++
	}

	if req.Country != nil {
		setParts = append(setParts, fmt.Sprintf("country = $%d", argIndex))
		args = append(args, *req.Country)
		argIndex++
	}

	if req.PostalCode != nil {
		setParts = append(setParts, fmt.Sprintf("postal_code = $%d", argIndex))
		args = append(args, *req.PostalCode)
		argIndex++
	}

	if req.DateOfBirth != nil {
		setParts = append(setParts, fmt.Sprintf("date_of_birth = $%d", argIndex))
		args = append(args, *req.DateOfBirth)
		argIndex++
	}

	if req.Bio != nil {
		setParts = append(setParts, fmt.Sprintf("bio = $%d", argIndex))
		args = append(args, *req.Bio)
		argIndex++
	}

	if req.Language != nil {
		setParts = append(setParts, fmt.Sprintf("language = $%d", argIndex))
		args = append(args, *req.Language)
		argIndex++
	}

	if req.Timezone != nil {
		setParts = append(setParts, fmt.Sprintf("timezone = $%d", argIndex))
		args = append(args, *req.Timezone)
		argIndex++
	}

	if len(setParts) > 0 {
		setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
		args = append(args, time.Now())
		argIndex++

		args = append(args, userID)

		query := fmt.Sprintf("UPDATE user_profiles SET %s WHERE user_id = $%d",
			joinStrings(setParts, ", "), argIndex)

		_, err = tx.Exec(query, args...)
		if err != nil {
			return fmt.Errorf("failed to update profile: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// DeleteProfile deletes user profile
func (s *UserProfileService) DeleteProfile(userID uuid.UUID) error {
	_, err := s.db.Exec("DELETE FROM user_profiles WHERE user_id = $1", userID)
	if err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}
	return nil
}

// Helper function to join strings
func joinStrings(parts []string, separator string) string {
	if len(parts) == 0 {
		return ""
	}

	result := parts[0]
	for i := 1; i < len(parts); i++ {
		result += separator + parts[i]
	}
	return result
}
