package user

import (
	"fmt"
	"math/rand/v2" // v1 대신 v2를 사용합니다.
	"sync"
)

type User struct {
	ID string
}

type UserPool struct {
	mu    sync.RWMutex
	users []*User
}

func NewUserPool() *UserPool {
	return &UserPool{
		users: make([]*User, 0),
	}
}

// LoadController가 필요한 유저 수만큼 확보
func (up *UserPool) EnsureUsers(required int) {
	up.mu.Lock()
	defer up.mu.Unlock()

	current := len(up.users)
	if current >= required {
		return
	}

	needed := required - current

	for i := 0; i < needed; i++ {
		newUser := &User{
			ID: fmt.Sprintf("user_%d", current+i+1),
		}
		up.users = append(up.users, newUser)
	}
}

// 세션 매니저가 랜덤 유저 뽑을 때 사용
func (up *UserPool) GetRandomUser() *User {
	up.mu.RLock()
	defer up.mu.RUnlock()

	n := len(up.users)
	if n == 0 {
		return nil
	}

	// [수정] rand.Intn -> rand.IntN (v2 전역 함수 사용하여 Thread-safe 확보)
	idx := rand.IntN(n)
	return up.users[idx]
}

func (up *UserPool) TotalCount() int {
	up.mu.RLock()
	defer up.mu.RUnlock()
	return len(up.users)
}
