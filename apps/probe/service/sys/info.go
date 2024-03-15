package sys

import "apps/common/message/data"

type InfoService struct {
}

// TODO implement GetInfo
// GetInfo gather system information
func (s *InfoService) GetInfo() (*data.PerformanceMetrics, error) {
	return nil, nil
}
