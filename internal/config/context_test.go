package config

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromContext(t *testing.T) {
	ctx := context.WithValue(t.Context(), configCtx, New())
	conf, ok := FromContext(ctx)
	assert.True(t, ok)
	assert.NotNil(t, conf)
}

func TestWithContext(t *testing.T) {
	ctx := WithContext(t.Context(), New())
	require.NotNil(t, ctx)
	conf, ok := ctx.Value(configCtx).(*Config)
	assert.True(t, ok)
	assert.NotNil(t, conf)
}
