package v2

import (
	"apps/center/model"
	"gorm.io/gorm"
)

type GraphService struct {
	Db *gorm.DB
}

func (g *GraphService) GetGraphData(chainId string) (model.GraphData, error) {

	var chain model.Chain
	err := g.Db.Where("id = ?", chainId).Select("origin_data").Take(&chain).Error

	if err != nil {
		return nil, err
	}

	data := []byte(chain.OriginData)

	graphData, err := model.UnmarshalGraphData(data)

	if err != nil {
		return nil, err
	}
	return graphData, nil
}

func (g *GraphService) UpdateGraphData(chainId string, data model.GraphData) error {
	bytes, err := data.Marshal()
	s := string(bytes)
	err = g.Db.Model(&model.Chain{}).Where("id = ?", chainId).Update("origin_data", s).Error
	return err
}
