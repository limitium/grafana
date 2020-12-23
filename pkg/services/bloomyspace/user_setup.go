package bloomyspace

import (
	"fmt"
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/events"
	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/registry"
	"github.com/grafana/grafana/pkg/setting"
)

type UserSetupService struct {
	log log.Logger
	Cfg *setting.Cfg `inject:""`
	Bus bus.Bus      `inject:""`
}

func init() {
	registry.RegisterService(&UserSetupService{})
}

func (bs *UserSetupService) Init() error {
	bs.log = log.New("UserSetupService")
	bs.Bus.AddEventListener(bs.SignUpCompleted)
	return nil
}

func (bs *UserSetupService) SignUpCompleted(event *events.SignUpCompleted) error {
	userQuery := &models.GetUserByEmailQuery{
		Email: event.Email,
	}
	if err := bus.Dispatch(userQuery); err != nil {
		return fmt.Errorf("could not find user: %w", err)
	}
	user := userQuery.Result

	if user.Id == 1 {
		return nil
	}

	err2 := bs.AddDummyAdmin(user)
	if err2 != nil {
		return err2
	}

	err3 := bs.RemoveAdminRole(user)
	if err3 != nil {
		return err3
	}

	err := bs.CopyDatasource(user)
	if err != nil {
		return err
	}

	return nil
}

func (bs *UserSetupService) RemoveAdminRole(user *models.User) error {
	orgUserCmd := &models.UpdateOrgUserCommand{
		Role:   models.ROLE_EDITOR,
		OrgId:  user.OrgId,
		UserId: user.Id,
	}
	if err := bus.Dispatch(orgUserCmd); err != nil {
		return fmt.Errorf("unable to update user role: %w", err)
	}
	return nil
}

func (bs *UserSetupService) AddDummyAdmin(user *models.User) error {
	bloomyspace := bs.Cfg.Raw.Section("bloomyspace")
	addUserCmd := &models.AddOrgUserCommand{
		Role:   models.ROLE_ADMIN,
		OrgId:  user.OrgId,
		UserId: bloomyspace.Key("dummy_org_admin_id").MustInt64(1),
	}
	if err := bus.Dispatch(addUserCmd); err != nil {
		return fmt.Errorf("unable to add admin to Org: %w", err)
	}
	return nil
}

func (bs *UserSetupService) CopyDatasource(user *models.User) error {
	bloomyspace := bs.Cfg.Raw.Section("bloomyspace")
	dsQuery := &models.GetDataSourceByIdQuery{
		OrgId: bloomyspace.Key("main_datasource_org_id").MustInt64(1),
		Id:    bloomyspace.Key("main_datasource_id").MustInt64(1),
	}
	if err := bus.Dispatch(dsQuery); err != nil {
		return fmt.Errorf("could not find datasource: %w", err)
	}
	ds := dsQuery.Result
	dsCmd := &models.AddDataSourceCommand{
		OrgId:     user.OrgId,
		Name:      ds.Name,
		Type:      ds.Type,
		Access:    models.DS_ACCESS_PROXY,
		Url:       ds.Url,
		IsDefault: true,
		ReadOnly:  true,
	}
	if err := bus.Dispatch(dsCmd); err != nil {
		return fmt.Errorf("unable to add datasource: %w", err)
	}
	return nil
}
