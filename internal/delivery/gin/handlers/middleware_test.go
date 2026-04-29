package handlers

import (
	"github.com/Shuhrat55/auth/pkg/auth"
	"github.com/Shuhrat55/auth/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	logger.Logger = zap.NewNop()