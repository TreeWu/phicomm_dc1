package data

import (
	"context"
	v1 "github.com/treewu/phicomm_dc1/api/gateway/v1"
	"gorm.io/gorm"
)

type PlanModel struct {
	gorm.Model
	UserId   int64             `gorm:"column:user_id;type:varchar(64);index" json:"user_id"`
	Name     string            `gorm:"column:name;type:varchar(64)" json:"name"`
	Cron     string            `gorm:"column:cron" json:"cron"`
	Enabled  bool              `gorm:"column:enabled" json:"enabled"`
	PlanType int32             `gorm:"column:plan_type" json:"plan_type"`
	Devices  []PlanModelDevice `gorm:"column:devices;serializer:json" json:"devices"`
	History  []CommandLog      `gorm:"-"`
}

type PlanModelDevice struct {
	DeviceId   string `json:"deviceId"`
	Switch1    *int32 `json:"switch1,omitempty"`
	Switch2    *int32 `json:"switch2,omitempty"`
	Switch3    *int32 `json:"switch3,omitempty"`
	SwitchMain *int32 `json:"switchMain,omitempty"`
}

func (p *PlanModel) TableName() string {
	return "plan"
}

type PlanRepo struct {
	db *gorm.DB
}

func NewPlanRepo(data *Data) *PlanRepo {
	return &PlanRepo{
		db: data.db,
	}
}

func (p *PlanRepo) Create(ctx context.Context, plan *PlanModel) (*PlanModel, error) {
	err := p.db.WithContext(ctx).Create(plan).Error
	if err != nil {
		return nil, err
	}
	return plan, nil
}

func (p *PlanRepo) Update(ctx context.Context, plan *PlanModel) error {
	return p.db.WithContext(ctx).Model(&PlanModel{}).Where("id = ?", plan.ID).Updates(plan).Error
}
func (p *PlanRepo) Get(ctx context.Context, id int64) (*PlanModel, error) {
	var plan PlanModel
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&plan).Error
	return &plan, err
}

func (p *PlanRepo) Delete(ctx context.Context, id int64) error {
	return p.db.WithContext(ctx).Model(&PlanModel{}).Delete("id = ? ", id).Error
}
func (p *PlanRepo) List(ctx context.Context, id int64) ([]PlanModel, error) {
	var plans []PlanModel
	if err := p.db.WithContext(ctx).Model(&PlanModel{}).Where("user_id = ?", id).Find(&plans).Error; err != nil {
		return nil, err
	}
	for i, plan := range plans {
		var history []CommandLog
		_ = p.db.WithContext(ctx).Debug().Model(&CommandLog{}).Where(CommandLog{
			PlanId: plan.ID,
		}).Order("id desc").Limit(10).Find(&history).Error
		plan.History = history
		plans[i] = plan
	}
	return plans, nil
}

func (p *PlanRepo) CountPlan(ctx context.Context, id int64) (int64, error) {
	var count int64
	if err := p.db.WithContext(ctx).Model(&PlanModel{}).Where("user_id = ?", id).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
func (p *PlanRepo) FindEnable(ctx context.Context) ([]PlanModel, error) {
	var plans []PlanModel
	err := p.db.WithContext(ctx).Model(&PlanModel{}).Where("enabled = ?", true).Where("plan_type = ?", v1.PlanType_PLAN_TYPE_AUTO).Find(&plans).Error
	return plans, err
}

func (p *PlanRepo) SwitchPlan(ctx context.Context, id int64, enable bool) error {
	return p.db.WithContext(ctx).Model(&PlanModel{}).Where("id = ?", id).Update("enabled", enable).Error
}
