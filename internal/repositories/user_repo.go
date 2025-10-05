package repositories

import (
	"context"
	"fmt"
	"locntp-user-counter/internal/models"
	"locntp-user-counter/pkg/utils"
	"sync"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserRepository struct {
	db    *gorm.DB
	cache *redis.Client
	mu    sync.Mutex // Protects concurrent counter updates
}

func NewUserRepository(db *gorm.DB, cache *redis.Client) *UserRepository {
	return &UserRepository{
		db:    db,
		cache: cache,
	}
}

// CreateUser creates a new user with counter initialized to 0
func (r *UserRepository) CreateUser(
	ctx context.Context,
	req *models.CreateUserRequest,
) (*models.User, error) {
	user := &models.User{
		Username: req.Username,
		Counter:  0,
	}

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Cache the initial counter value
	if err := utils.SetUserCount(ctx, r.cache, user.ID, user.Counter); err != nil {
		// Log but don't fail - utils is not critical for creation
		fmt.Printf("Warning: failed to utils user counter: %v\n", err)
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	var user models.User

	if err := r.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetUserCount retrieves the counter value (from utils first, fallback to DB)
func (r *UserRepository) GetUserCount(ctx context.Context, userID uint) (int64, bool, error) {
	// Try utils first
	count, cached, err := utils.GetUserCount(ctx, r.cache, userID)
	if err != nil {
		fmt.Printf("Warning: utils error: %v\n", err)
	}

	if cached {
		return count, true, nil
	}

	// Cache miss - fetch from database
	var user models.User
	if err := r.db.WithContext(ctx).Select("counter").First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, false, fmt.Errorf("user not found")
		}
		return 0, false, fmt.Errorf("failed to get user counter: %w", err)
	}

	// Update utils with DB value
	if err := utils.SetUserCount(ctx, r.cache, userID, user.Counter); err != nil {
		fmt.Printf("Warning: failed to update utils: %v\n", err)
	}

	return user.Counter, false, nil
}

// IncrementUserCount increments the user's counter (with concurrency protection)
func (r *UserRepository) IncrementUserCount(
	ctx context.Context,
	userID uint,
	amount int64,
) (*models.User, error) {
	// Lock to prevent race conditions on concurrent increments
	r.mu.Lock()
	defer r.mu.Unlock()

	// Start a transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Lock the row for update
	var user models.User
	if err := tx.Clauses().First(&user, userID).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Increment counter
	user.Counter += amount

	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update counter: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Update utils
	if err := utils.SetUserCount(ctx, r.cache, user.ID, user.Counter); err != nil {
		fmt.Printf("Warning: failed to update utils: %v\n", err)
	}

	return &user, nil
}

// GetAllUsers retrieves all users
func (r *UserRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User

	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil
}

// DeleteUser deletes a user and their utils
func (r *UserRepository) DeleteUser(ctx context.Context, userID uint) error {
	if err := r.db.WithContext(ctx).Delete(&models.User{}, userID).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Delete from utils
	if err := utils.DeleteUserCount(ctx, r.cache, userID); err != nil {
		fmt.Printf("Warning: failed to delete from utils: %v\n", err)
	}

	return nil
}
