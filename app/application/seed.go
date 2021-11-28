package application

import (
	"context"
	"log"
	"os"

	"github.com/linzhengen/ddd-gin-admin/app/domain/menu"
	"github.com/linzhengen/ddd-gin-admin/app/domain/menu/menuaction"
	"github.com/linzhengen/ddd-gin-admin/app/domain/menu/menuactionresource"
	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"
	"github.com/linzhengen/ddd-gin-admin/app/domain/trans"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/yaml"
)

type Seed interface {
	Execute(ctx context.Context, menuSeedPath string) error
}

type SeedMenus []struct {
	Name     string `yaml:"name"`
	Icon     string `yaml:"icon"`
	Router   string `yaml:"router,omitempty"`
	Sequence int    `yaml:"sequence"`
	Actions  []struct {
		Code      string `yaml:"code"`
		Name      string `yaml:"name"`
		Resources []struct {
			Method string `yaml:"method"`
			Path   string `yaml:"path"`
		} `yaml:"resources"`
	} `yaml:"actions,omitempty"`
	Children SeedMenus
}

func NewSeed(
	menuSvc menu.Service,
	transRepo trans.Repository,
) Seed {
	return &seedApp{
		menuSvc:   menuSvc,
		transRepo: transRepo,
	}
}

type seedApp struct {
	menuSvc   menu.Service
	transRepo trans.Repository
}

func (s seedApp) Execute(ctx context.Context, menuSeedPath string) error {
	if err := s.menuSeed(ctx, menuSeedPath); err != nil {
		return err
	}
	return nil
}

func (s seedApp) menuSeed(ctx context.Context, menuSeedPath string) error {
	_, pr, err := s.menuSvc.Query(ctx, menu.QueryParam{
		PaginationParam: pagination.Param{OnlyCount: true},
	})
	if err != nil {
		return err
	}
	if pr.Total > 0 {
		return nil
	}
	data, err := s.readMenuData(menuSeedPath)
	if err != nil {
		return err
	}

	return s.createMenus(ctx, "", data)
}

func (s seedApp) readMenuData(name string) (SeedMenus, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data SeedMenus
	d := yaml.NewDecoder(file)
	d.SetStrict(true)
	err = d.Decode(&data)
	log.Printf("%+v", data)
	return data, err
}

func (s seedApp) createMenus(ctx context.Context, parentID string, list SeedMenus) error {
	return s.transRepo.Exec(ctx, func(ctx context.Context) error {
		for _, item := range list {
			var as menuaction.MenuActions
			for _, action := range item.Actions {
				var ars menuactionresource.MenuActionResources
				for _, r := range action.Resources {
					ars = append(ars, &menuactionresource.MenuActionResource{
						Method: r.Method,
						Path:   r.Path,
					})
				}
				as = append(as, &menuaction.MenuAction{
					Code:      action.Code,
					Name:      action.Name,
					Resources: ars,
				})
			}
			sitem := &menu.Menu{
				Name:       item.Name,
				Sequence:   item.Sequence,
				Icon:       item.Icon,
				Router:     item.Router,
				ParentID:   parentID,
				Status:     1,
				ShowStatus: 1,
				Actions:    as,
			}

			menuID, err := s.menuSvc.Create(ctx, sitem)
			if err != nil {
				return err
			}

			if item.Children != nil && len(item.Children) > 0 {
				err := s.createMenus(ctx, menuID, item.Children)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}
