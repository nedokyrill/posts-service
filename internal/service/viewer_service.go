package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nedokyrill/posts-service/internal/models"
	"github.com/nedokyrill/posts-service/pkg/consts"
	"github.com/nedokyrill/posts-service/pkg/logger"
	"github.com/nedokyrill/posts-service/pkg/utils"
)

type Viewer struct {
	ch chan *models.Comment
	id int
}

type ViewerServiceImpl struct {
	viewers map[uuid.UUID][]Viewer
	cnt     int
	mu      sync.Mutex
}

func NewViewerService() *ViewerServiceImpl {
	return &ViewerServiceImpl{
		viewers: make(map[uuid.UUID][]Viewer),
		cnt:     0,
		mu:      sync.Mutex{},
	}
}

// добавляем подписчика
func (s *ViewerServiceImpl) CreateViewer(_ context.Context, postId uuid.UUID) (int, chan *models.Comment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	comm := make(chan *models.Comment)
	s.viewers[postId] = append(s.viewers[postId], Viewer{ch: comm, id: s.cnt})
	s.cnt++

	logger.Logger.Infof("create viewer for post id %s", postId.String())
	return s.cnt, comm, nil
}

// удаляем подписчика из пула при закрытии подписки
func (s *ViewerServiceImpl) DeleteViewer(_ context.Context, postId uuid.UUID, id int) error {
	s.mu.Lock()

	viewers, ok := s.viewers[postId]
	if !ok {
		s.mu.Unlock()
		logger.Logger.Error(fmt.Sprintf("no post with postId: %s", postId.String()))
		return utils.GqlError{Msg: fmt.Sprintf("no post with postId: %s", postId.String()),
			Type: consts.BadRequestType}
	}
	for i, viewer := range viewers {
		if viewer.id == id {
			s.viewers[postId] = append(viewers[:i], viewers[i+1:]...)
			close(viewer.ch)
			break
		}
	}

	s.mu.Unlock()

	logger.Logger.Infof("delete viewer for post id %s", postId.String())
	return nil
}

// отправляем уведомление в виде комментария всем подписчикам
func (s *ViewerServiceImpl) NotifyViewers(ctx context.Context, postId uuid.UUID, comm models.Comment) error {
	s.mu.Lock()

	viewers, ok := s.viewers[postId]
	if !ok {
		s.mu.Unlock()
		logger.Logger.Error(fmt.Sprintf("no post with postId: %s", postId.String()))
		return utils.GqlError{Msg: fmt.Sprintf("no post with postId: %s", postId.String()),
			Type: consts.BadRequestType}
	}

	snap := make([]Viewer, len(viewers))
	copy(snap, viewers)
	s.mu.Unlock() // сделали копию каналов и разлочили мьютекс, чтобы далее не блокировать данные

	for _, v := range snap {
		select { // селект для неблокирующей отправки (канал для нотификаций может быть закрыт или заполнен)
		case v.ch <- &comm: // отправляем нотификацию
		case <-ctx.Done():
			logger.Logger.Info("context canceled")
		case <-time.After(500 * time.Millisecond): // если за 500мс канал не освободился, логируем айдишник и пропускаем
			logger.Logger.Error(fmt.Sprintf("viewer with id: %v, channel full", v.id))
		}
	}

	logger.Logger.Info(fmt.Sprintf("notify viewer for postId %s successfully", postId.String()))
	return nil
}
