package usecases

import (
	"github.com/lucastomic/snoopycoon_back/database"
	"github.com/lucastomic/snoopycoon_back/domain"
)

func CreateTopic(topic *domain.Topic) error {
	err := database.DB.Create(&topic).Error
	if err != nil {
		return err
	}
	return nil
}

func GetTopics() ([]domain.Topic, error) {
	topics := []domain.Topic{}
	err := database.DB.Find(&topics).Error
	if err != nil {
		return nil, err
	}
	return topics, nil
}

func DeleteTopic(id uint) error {
	err := database.DB.Delete(&domain.Topic{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateTopic(id uint, topic domain.Topic) error {
	err := database.DB.Model(&domain.Topic{}).Where("id = ?", id).Updates(&topic).Error
	if err != nil {
		return err
	}
	return nil
}
