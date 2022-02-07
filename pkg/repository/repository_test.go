package repository

import (
	"context"
	"testing"

	"github.com/crikke/cms/pkg/config"
	"github.com/crikke/cms/pkg/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetFromDb(t *testing.T) {

	cfg := config.LoadServerConfiguration()
	r, err := NewRepository(context.TODO(), cfg)

	assert.NoError(t, err)

	r.GetContent(context.Background(), domain.ContentReference{ID: uuid.MustParse("a1f6da93-80c9-4315-a012-1ea4249d7413")})
}
